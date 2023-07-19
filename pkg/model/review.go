package model

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

type Review struct {
	ID        int       `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	ProductID int       `db:"product_id" json:"product_id"`
	UserID    int       `db:"user_id" json:"user_id"`
	Username  string    `db:"username" json:"username"`
	Text      string    `db:"text" json:"text" validate:"required"`
	Category  string    `db:"category" json:"category" validate:"review_category,required"`
}

type UpdateReviewInput struct {
	Text      *string    `json:"text"`
	Category  *string    `json:"category" validate:"review_category"`
	UpdatedAt *time.Time `json:"updated_at"`
}

const (
	POSITIVE string = "positive"
	NEUTRAL  string = "neutral"
	NEGATIVE string = "negative"
)

func (i UpdateReviewInput) Validate() error {
	if i.Text == nil && i.UpdatedAt == nil && i.Category != nil {
		return errors.New("update structure has no values")
	}

	return nil
}

func ValidateReviewCategory(fl validator.FieldLevel) bool {
	category := fl.Field().String()
	return !(category != POSITIVE && category != NEGATIVE && category != NEUTRAL)
}
