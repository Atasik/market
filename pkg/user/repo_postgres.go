package user

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	ErrNoUser     = errors.New("no user found")
	ErrUserExists = errors.New("user already exists")
	ErrBadPass    = errors.New("invalid password")
)

type UserPostgresqlRepository struct {
	DB *sqlx.DB
}

func NewPostgresqlRepo(db *sqlx.DB) *UserPostgresqlRepository {
	return &UserPostgresqlRepository{DB: db}
}

func (repo *UserPostgresqlRepository) Authorize(login, pass string) (User, error) {
	var user User
	query := "SELECT * FROM users WHERE username = $1"

	if err := repo.DB.Get(&user, query, login); err != nil {
		fmt.Print(user)
		if err == sql.ErrNoRows {
			return User{}, ErrNoUser
		}
		return User{}, err
	}

	if user.Password != pass {
		return User{}, ErrBadPass
	}

	return user, nil
}

func (repo *UserPostgresqlRepository) Register(login, pass string) (int, error) {
	tx, err := repo.DB.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	query := "INSERT INTO users (username, user_mode, password) VALUES ($1, $2, $3) RETURNING id"

	row := tx.QueryRow(query, login, "user", pass)
	err = row.Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, ErrUserExists
	}

	return id, tx.Commit()
}
