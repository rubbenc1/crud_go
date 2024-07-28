package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"myproject/internal/models"
	"myproject/internal/repository/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBookHandler_CreateBook(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    bookRepo := mocks.NewMockBookRepository(ctrl)
    shelfRepo := mocks.NewMockShelfRepository(ctrl)
    handler := NewBookHandler(bookRepo, shelfRepo)

    testCases := []struct {
        name            string
        shelfId         int
        book            models.Book
        mockShelfResp   *models.Shelf
        mockAddBookResp *models.Book
        mockAddBookErr  error
        expectedCode    int
    }{
        {
            name: "Successful creation",
            shelfId: 1,
            book: models.Book{
                Name: "Test Book",
                PublishedYear: 2024,
                ShelfId: 1,
            },
            mockShelfResp: &models.Shelf{ShelfId: 1},
            mockAddBookResp: &models.Book{
                BookId: 1,
                Name: "Test Book",
                PublishedYear: 2024,
                ShelfId: 1,
            },
            mockAddBookErr: nil,
            expectedCode: http.StatusCreated,
        },
        {
            name: "Shelf not found",
            shelfId: 2,
            book: models.Book{
                Name: "Another Book",
                PublishedYear: 2024,
                ShelfId: 2,
            },
            mockShelfResp: nil,
            mockAddBookResp: nil,
            mockAddBookErr: nil,
            expectedCode: http.StatusNotFound,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Set up mock expectations
            shelfRepo.EXPECT().GetShelf(tc.shelfId).Return(tc.mockShelfResp, nil).Times(1)

            if tc.mockShelfResp != nil {
                if tc.mockAddBookResp != nil {
                    bookRepo.EXPECT().AddBook(gomock.Eq(&tc.book)).Return(tc.mockAddBookResp, tc.mockAddBookErr).Times(1)
                } else {
                    bookRepo.EXPECT().AddBook(gomock.Eq(&tc.book)).Return(nil, tc.mockAddBookErr).Times(1)
                }
            }

            // Create request body
            body, err := json.Marshal(tc.book)
            assert.NoError(t, err)

            // Create HTTP request
            req, err := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
            assert.NoError(t, err)
            req.Header.Set("Content-Type", "application/json")

            // Create response recorder
            rr := httptest.NewRecorder()

            // Call the handler
            handler.CreateBook(rr, req)


            // Check response code
            assert.Equal(t, tc.expectedCode, rr.Code)
        })
    }
}

func TestBookHandler_GetBook(t *testing.T){
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    bookRepo := mocks.NewMockBookRepository(ctrl)
    shelfRepo := mocks.NewMockShelfRepository(ctrl)
    handler := NewBookHandler(bookRepo, shelfRepo)

	testCases := []struct {
		name string
		id int
		mockBookResp *models.Book
		mockBookErr error
		expectedCode int
		expectedBody map[string]interface{}
	}{
		{
			name: "Successful retrieval",
			id: 1,
			mockBookResp: &models.Book {
				BookId: 1,
				Name: "Test Book",
				PublishedYear: 2024,
				ShelfId: 1,
			},
			mockBookErr: nil,
			expectedCode: http.StatusOK,
			expectedBody: map[string]interface{}{
				"book_id":       float64(1),
				"name":          "Test Book",
				"PublishedYear": float64(2024),
				"ShelfId":       float64(1),
			},
		},
		{
			name: "Book not found",
			id: 2,
			mockBookResp: nil,
			mockBookErr: nil,
			expectedCode: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "Book not found",
			},
		},
	}

	for _,tc := range testCases {
		t.Run(tc.name, func(t *testing.T){
			bookRepo.EXPECT().GetBook(tc.id).Return(tc.mockBookResp, tc.mockBookErr).Times(1)
			req, err := http.NewRequest("GET", fmt.Sprintf("/books/%d", tc.id), nil)
			assert.NoError(t, err)
			rr := httptest.NewRecorder()
			// Set up the router and context
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", fmt.Sprintf("%d", tc.id))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Call the handler
			handler.GetBook(rr, req)

			// Check response code
			assert.Equal(t, tc.expectedCode, rr.Code)

			// Check response body
			var actualBody map[string]interface{}
			err = json.Unmarshal(rr.Body.Bytes(), &actualBody)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBody, actualBody)
			
		})
	}
}

func TestHandler_UpdateBook(t *testing.T){
	ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    bookRepo := mocks.NewMockBookRepository(ctrl)
    shelfRepo := mocks.NewMockShelfRepository(ctrl)
    handler := NewBookHandler(bookRepo, shelfRepo)

	testCases := []struct {
		name string
		id int
		book models.Book
		mockShelfResp *models.Shelf
		mockShelfErr error
        mockUpdatedBookErr error
        expectedCode int
		expectedBody map[string]interface{}
	}{
		{
			name: "Successful update",
			id: 1,
			book: models.Book{
				Name: "Updated Book",
				PublishedYear: 2025,
				ShelfId: 1,
			},
			mockShelfResp: &models.Shelf{ShelfId: 1},
			mockShelfErr: nil,
			mockUpdatedBookErr: nil,
			expectedCode: http.StatusNoContent,
			expectedBody: nil,
		},
		{
			name: "Invalid book ID",
			id:   0,
			book: models.Book{
				Name:          "Invalid Book",
				PublishedYear: 2025,
				ShelfId:       2,
			},
			mockShelfResp:     nil,
			mockShelfErr:      nil,
			mockUpdatedBookErr: nil,
			expectedCode:      http.StatusBadRequest,
			expectedBody:      map[string]interface{}{"error": "Invalid book ID"},
		},
		{
			name: "Shelf does not exist",
			id:   1,
			book: models.Book{
				Name:          "Nonexistent Shelf",
				PublishedYear: 2025,
				ShelfId:       2,
			},
			mockShelfResp:     nil,
			mockShelfErr:      fmt.Errorf("shelf not found"),
			mockUpdatedBookErr: nil,
			expectedCode:      http.StatusBadRequest,
			expectedBody:      map[string]interface{}{"error": "Shelf does not exist"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up mock expectations for shelfRepo
			if tc.id != 0 && (tc.mockShelfResp != nil || tc.mockShelfErr != nil) {
				shelfRepo.EXPECT().GetShelf(tc.book.ShelfId).Return(tc.mockShelfResp, tc.mockShelfErr).Times(1)
			}

			if tc.mockShelfResp != nil {
				// Set up mock expectations for bookRepo only if shelf exists and ID is valid
				bookRepo.EXPECT().UpdateBook(tc.id, &tc.book).Return(tc.mockUpdatedBookErr).Times(1)
			}

			// Create request body
			body, err := json.Marshal(tc.book)
			assert.NoError(t, err)

			// Create HTTP request
			req, err := http.NewRequest("PUT", fmt.Sprintf("/books/%d", tc.id), bytes.NewBuffer(body))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Set up the router and context
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", fmt.Sprintf("%d", tc.id))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handler.UpdateBook(rr, req)

			// Check response code
			assert.Equal(t, tc.expectedCode, rr.Code)

			// Check response body if expected
			if tc.expectedBody != nil {
				var responseBody map[string]interface{}
				err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedBody, responseBody)
			}
		})
	}
}


func TestHandler_DeleteBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

    bookRepo := mocks.NewMockBookRepository(ctrl)
    shelfRepo := mocks.NewMockShelfRepository(ctrl)
    handler := NewBookHandler(bookRepo, shelfRepo)

	testCases := []struct {
		name             string
		id               int
		mockDeleteErr    error
		expectedCode     int
		expectedBody     map[string]interface{}
	}{
		{
			name:          "Successful delete",
			id:            1,
			mockDeleteErr: nil,
			expectedCode:  http.StatusNoContent,
			expectedBody:  nil,
		},
		{
			name:          "Invalid book ID",
			id:            0, // invalid ID for testing
			expectedCode:  http.StatusBadRequest,
			expectedBody:  map[string]interface{}{"error": "Invalid book ID"},
		},
		{
			name:          "Book not found",
			id:            1,
			mockDeleteErr: fmt.Errorf("book not found"),
			expectedCode:  http.StatusNotFound,
			expectedBody:  map[string]interface{}{"error": "Book not found"},
		},
		{
			name:          "Internal server error",
			id:            1,
			mockDeleteErr: fmt.Errorf("internal error"),
			expectedCode:  http.StatusInternalServerError,
			expectedBody:  map[string]interface{}{"error": "internal error"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up mock expectations for bookRepo if necessary
			if tc.id != 0 && (tc.mockDeleteErr != nil || tc.mockDeleteErr == nil) {
				bookRepo.EXPECT().DeleteBook(tc.id).Return(tc.mockDeleteErr).Times(1)
			}

			// Create HTTP request
			req, err := http.NewRequest("DELETE", fmt.Sprintf("/books/%d", tc.id), nil)
			assert.NoError(t, err)

			// Set up the router and context
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", fmt.Sprintf("%d", tc.id))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handler.DeleteBook(rr, req)

			// Check response code
			assert.Equal(t, tc.expectedCode, rr.Code)

			// Check response body if expected
			if tc.expectedBody != nil {
				var responseBody map[string]interface{}
				err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedBody, responseBody)
			}
		})
	}
}