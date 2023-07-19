package model

import (
	"errors"
	"time"
)

type Review struct {
	ID        int       `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	ProductID int       `db:"product_id" json:"product_id"`
	UserID    int       `db:"user_id" json:"user_id"`
	Username  string    `db:"username" json:"username"`
	Text      string    `db:"text" json:"text"`
	Category  Category  `db:"category" json:"category"`
}

type UpdateReviewInput struct {
	Text      *string    `json:"text"`
	Category  Category   `json:"category"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type Category string

const (
	POSITIVE Category = "positive"
	NEUTRAL  Category = "neutral"
	NEGATIVE Category = "negative"
)

func (i UpdateReviewInput) Validate() error {
	if i.Text == nil && i.UpdatedAt == nil && (i.Category != POSITIVE && i.Category != NEUTRAL && i.Category != NEGATIVE) {
		return errors.New("update structure has no values")
	}

	return nil
}
