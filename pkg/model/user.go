package model

type User struct {
	ID       int    `db:"id" json:"-"`
	Role     Role   `db:"role" json:"role"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
}

type Role string

const (
	ADMIN  Role = "admin"
	SELLER Role = "seller"
	USER   Role = "user"
)
