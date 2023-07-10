package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"market/pkg/model"

	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	GetUser(login string) (model.User, error)
	CreateUser(login, hash string) (int, error)
}

var (
	ErrNoUser     = errors.New("no user found")
	ErrUserExists = errors.New("user already exists")
	ErrBadPass    = errors.New("invalid password")
)

type UserPostgresqlRepository struct {
	DB *sqlx.DB
}

func NewUserPostgresqlRepo(db *sqlx.DB) *UserPostgresqlRepository {
	return &UserPostgresqlRepository{DB: db}
}

func (repo *UserPostgresqlRepository) GetUser(login string) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE username = $1", usersTable)

	if err := repo.DB.Get(&user, query, login); err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, ErrNoUser
		}
		return model.User{}, err
	}

	return user, nil
}

func (repo *UserPostgresqlRepository) CreateUser(login, hash string) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, user_mode, password) VALUES ($1, $2, $3) RETURNING id", usersTable)

	row := repo.DB.QueryRow(query, login, "user", hash)
	err := row.Scan(&id)
	if err != nil {
		return 0, ErrUserExists
	}

	return id, nil
}
