package model

import (
	"errors"
	"time"
)

type Product struct {
	ID             int       `db:"id"`
	Title          string    `schema:"title" db:"title"`
	Price          float32   `schema:"price" db:"price"`
	Tag            string    `schema:"tag" db:"tag"`
	Type           string    `schema:"type" db:"type"`
	Description    string    `schema:"description" db:"description"`
	Count          int       `db:"count"`
	PurchasedCount int       `db:"purchased_count"`
	OrderID        int       `db:"order_id"`
	CreationDate   time.Time `schema:"creation_date" db:"creation_date"`
	Views          int       `db:"views"`
	ImageURL       string    `db:"image_url"`
}

type UpdateProductInput struct {
	Title       *string  `schema:"title"`
	Price       *float32 `schema:"price"`
	Tag         *string  `schema:"tag"`
	Type        *string  `schema:"type"`
	Description *string  `schema:"description"`
	Count       *int     `schema:"count"`
	Views       *int     `schema:"views"`
	ImageURL    *string  `schema:"image_url"`
}

func (i UpdateProductInput) Validate() error {
	if i.Title == nil && i.Price == nil && i.Tag == nil && i.Type == nil && i.Description == nil && i.Count == nil && i.Views == nil && i.ImageURL == nil {
		return errors.New("update structure has no values")
	}

	return nil
}
