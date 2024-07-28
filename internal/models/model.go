package models

import validation "github.com/go-ozzo/ozzo-validation"

type Book struct {
	BookId        int    `json:"book_id"`
	Name          string `json:"name"`
	PublishedYear int    `json:"PublishedYear"`
	ShelfId       int    `json:"ShelfId"`
}

type Shelf struct {
	ShelfId int `json:"shelf_id"`
}

func (b *Book) Validate() error {
	return validation.ValidateStruct(b,
		validation.Field(&b.Name, validation.Required, validation.Length(1,255)),
		validation.Field(&b.PublishedYear, validation.Required, validation.Min(0)),
		validation.Field(&b.ShelfId, validation.Required),
	)
}