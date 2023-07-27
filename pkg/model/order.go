package model

import (
	"errors"
	"time"
)

type Order struct {
	ID          int       `db:"id" json:"id"`
	UserID      int       `db:"user_id" json:"user_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	DeliveredAt time.Time `db:"delivered_at" json:"delivered_at"`
	Products    []Product `json:"products,omitempty"`
}

type OrderQueryInput struct {
	Limit     int
	Offset    int
	SortBy    string
	SortOrder string
}

func (i OrderQueryInput) Validate() error {
	if i.SortBy != SortByDate || (i.SortOrder != ASCENDING && i.SortOrder != DESCENDING) {
		return errors.New("invalid sort query")
	}

	return nil
}
