package model

import (
	"errors"
	"time"
)

type Product struct {
	ID              int       `db:"id" json:"id"`
	UserID          int       `db:"user_id" json:"user_id"`
	Title           string    `db:"title" json:"title" schema:"title" validate:"required"`
	Price           float32   `db:"price" json:"price" schema:"price" validate:"required"`
	Tag             string    `db:"tag" json:"tag" schema:"tag"`
	Category        string    `db:"category" json:"category" schema:"category" validate:"required"`
	Description     string    `db:"description" json:"description" schema:"description"`
	Amount          int       `db:"amount" json:"amount" schema:"amount" validate:"required"`
	PurchasedAmount int       `db:"purchased_amount" json:"purchased_amount,omitempty"`
	OrderID         int       `db:"order_id" json:"order_id,omitempty"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
	Views           int       `db:"views" json:"views"`
	ImageURL        string    `db:"image_url" json:"image_url"`
	ImageID         string    `db:"image_id"`
	Reviews         []Review  `json:"reviews"`
	RelatedProducts []Product `json:"related_products"`
}

type UpdateProductInput struct {
	Title       *string    `json:"title"`
	Price       *float32   `json:"price"`
	Tag         *string    `json:"tag"`
	Type        *string    `json:"type"`
	Description *string    `json:"description"`
	UpdatedAt   *time.Time `json:"updated_at"`
	Amount      *int       `json:"amount"`
	Views       *int       `json:"views"`
	ImageURL    *string    `json:"image_url"`
	ImageID     *string
}

type ProductQueryInput struct {
	Limit     int
	Offset    int
	ProductID int
	SortBy    string
	SortOrder string
}

func (i ProductQueryInput) Validate() error {
	if i.SortBy != SortByViews && i.SortBy != SortByPrice && i.SortBy != SortByDate || (i.SortOrder != ASCENDING && i.SortOrder != DESCENDING) {
		return errors.New("invalid sort query")
	}

	return nil
}

func (i UpdateProductInput) Validate() error {
	if i.Title == nil && i.Price == nil && i.Tag == nil && i.Type == nil && i.Description == nil && i.Amount == nil && i.Views == nil && i.UpdatedAt == nil && (i.ImageURL == nil && i.ImageURL != i.ImageID) {
		return errors.New("update structure has no values")
	}

	return nil
}
