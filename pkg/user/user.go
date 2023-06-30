package user

type User struct {
	ID       int    `db:"id"`
	UserMode string `db:"user_mode"`
	Username string `db:"username"`
	Password string `db:"password"`
}

type UserRepo interface {
	Authorize(login, pass string) (User, error)
	Register(login, pass string) (int, error)
}
