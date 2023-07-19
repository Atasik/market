package model

import "github.com/go-playground/validator/v10"

type User struct {
	ID       int    `db:"id" json:"-"`
	Role     string `db:"role" json:"role" validate:"user_role,required"`
	Username string `db:"username" json:"username" validate:"required"`
	Password string `db:"password" json:"password" validate:"required"`
}

const (
	ADMIN  string = "admin"
	SELLER string = "seller"
	USER   string = "user"
)

func ValidateRole(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	return !(role != SELLER && role != USER)
}
