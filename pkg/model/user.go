package model

type User struct {
	ID       int    `db:"id"`
	UserMode string `db:"user_mode"`
	Username string `db:"username" schema:"username"`
	Password string `db:"password" schema:"password"`
}
