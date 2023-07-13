package model

type User struct {
	ID       int    `db:"id"`
	Role     Role   `db:"role"`
	Username string `db:"username" schema:"username"`
	Password string `db:"password" schema:"password"`
}

type Role string

const (
	ADMIN Role = "admin"
	USER  Role = "user"
)
