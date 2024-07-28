package handlers

import (
	"encoding/json"
	"myproject/internal/models"
	"myproject/internal/repository"
	"net/http"
	"strconv"
	"github.com/go-chi/chi/v5"
)

type BookHandler struct {
	bookRepo repository.BookRepository
	shelfRepo repository.ShelfRepository
}

func NewBookHandler(bookRepo repository.BookRepository, shelfRepo repository.ShelfRepository) *BookHandler {
	return &BookHandler{bookRepo: bookRepo, shelfRepo: shelfRepo}
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
    var book models.Book

    // Decode the request body
    if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Validate the book
    if err := book.Validate(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Check if the shelf exists
    shelf, err := h.shelfRepo.GetShelf(book.ShelfId)
    if err != nil || shelf == nil {
        http.Error(w, "Shelf not found", http.StatusNotFound)
        return
    }

    // Add the book to the repository
    createdBook, err := h.bookRepo.AddBook(&book)
    if err != nil {
        http.Error(w, "Failed to create book", http.StatusInternalServerError)
        return
    }

    // Set response header and encode the created book
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(createdBook); err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    }
}


func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	bookId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	book, err := h.bookRepo.GetBook(bookId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if book == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Book not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request){
	bookId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || bookId <= 0{
		http.Error(w, `{"error":"Invalid book ID"}`, http.StatusBadRequest)
		return
	}

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
        return
	}

	if err := book.Validate(); err != nil {
        http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
        return
    }

	_, err = h.shelfRepo.GetShelf(book.ShelfId)
    if err != nil {
        http.Error(w, `{"error": "Shelf does not exist"}`, http.StatusBadRequest)
        return
    }

    err = h.bookRepo.UpdateBook(bookId, &book)
    if err != nil {
        http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)

}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request){
	bookId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || bookId <= 0 {
		http.Error(w,`{"error": "Invalid book ID"}`, http.StatusBadRequest)
		return
	}
	err = h.bookRepo.DeleteBook(bookId)
	if err != nil {
		if err.Error() == "book not found" {
			http.Error(w, `{"error": "Book not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}