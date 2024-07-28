package repository

import (
	"database/sql"
	"fmt"
	"myproject/internal/models"
)

type BookRepository interface {
    AddBook(book *models.Book) (*models.Book, error)
    GetBook(id int) (*models.Book, error)
    UpdateBook(id int, updatedBook *models.Book) error
    DeleteBook(id int) error
}

type bookRepository struct {
    DB *sql.DB
}

func NewBookRepository(db *sql.DB) BookRepository {
    return &bookRepository{DB: db}
}

func (r *bookRepository) AddBook(book *models.Book) (*models.Book, error) {
    query := `INSERT INTO books (Name, PublishedYear, ShelfId) VALUES ($1, $2, $3) RETURNING BookId`
    err := r.DB.QueryRow(query, book.Name, book.PublishedYear, book.ShelfId).Scan(&book.BookId)
    if err != nil {
        return nil, fmt.Errorf("failed to add book: %w", err)
    }
    return book, nil
}

func (r *bookRepository) GetBook(id int)(*models.Book,error){
	query :=`SELECT BookId, Name, PublishedYear, ShelfId FROM books WHERE BookId = $1`
	book := &models.Book{}
	err := r.DB.QueryRow(query, id).Scan(&book.BookId, &book.Name, &book.PublishedYear, &book.ShelfId)
	if err != nil {
		return nil, err
	}
	return book, nil
} 

func (r *bookRepository) UpdateBook(id int, updatedBook *models.Book) error{
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p:=recover(); p!= nil{
			tx.Rollback()
			panic(p)
		}else if err !=nil{
			tx.Rollback()
		}else {
			err = tx.Commit()
		}
	}()

	query :=`UPDATE books SET Name = $1, PublishedYear = $2, ShelfId = $3 WHERE BookId = $4`
	_, err = r.DB.Exec(query, updatedBook.Name, updatedBook.PublishedYear, updatedBook.ShelfId, id)
	if err != nil{
		return err
	}
	return nil
}

func (r *bookRepository) DeleteBook(id int) error {
	query := `DELETE FROM books WHERE BookId = $1`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}
	return nil
}

