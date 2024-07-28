package repository

import (
	"database/sql"
	"fmt"
	"myproject/internal/models"
)

type ShelfRepository interface {
    AddShelf() (*models.Shelf, error)
    GetShelf(id int) (*models.Shelf, error)
    DeleteShelf(id int) error
}

type shelfRepository struct {
    DB *sql.DB
}

func NewShelfRepository(db *sql.DB) ShelfRepository {
    return &shelfRepository{DB: db}
}


func (r *shelfRepository) AddShelf() (*models.Shelf, error) {
    query := `INSERT INTO shelves DEFAULT VALUES RETURNING ShelfId`
    shelf := &models.Shelf{}
    err := r.DB.QueryRow(query).Scan(&shelf.ShelfId)
    if err != nil {
        return nil, err
    }
    return shelf, nil
}

func (r *shelfRepository) GetShelf(id int) (*models.Shelf, error) {
	query := `SELECT ShelfId FROM shelves WHERE ShelfId = $1`
	shelf := &models.Shelf{}
	err := r.DB.QueryRow(query, id).Scan(&shelf.ShelfId)
	if err != nil {
		return nil, err
	}
	return shelf, nil
}

func (r *shelfRepository) DeleteShelf(id int) error {
	query := `DELETE FROM shelves WHERE ShelfId = $1`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete shelf: %w", err)
	}
	return nil
}
