package model

type User struct {
	ID       int    `db:"id"`
	UserMode string `db:"user_mode"`
	Username string `db:"username"`
	Password string `db:"password"`
}