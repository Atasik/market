package model

import "github.com/go-playground/validator/v10"

const (
	SortByViews = "views"
	SortByPrice = "price"
	SortByDate  = "created_at"
	ASCENDING   = "ASC"
	DESCENDING  = "DESC"
)

func RegisterCustomValidations(v *validator.Validate) error {
	if err := v.RegisterValidation("user_role", ValidateRole); err != nil {
		return err
	}
	return v.RegisterValidation("review_category", ValidateReviewCategory)
}

type QueryInput struct {
	Limit     int
	Offset    int
	SortBy    string
	SortOrder string
}
