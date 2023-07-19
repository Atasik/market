package model

import "github.com/go-playground/validator/v10"

func RegisterCustomValidations(v *validator.Validate) error {
	err := v.RegisterValidation("user_role", ValidateRole)
	if err != nil {
		return err
	}
	return v.RegisterValidation("review_category", ValidateReviewCategory)
}
