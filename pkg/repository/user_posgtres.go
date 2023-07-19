package repository

import (
	"fmt"
	"market/pkg/model"

	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	GetUser(login string) (model.User, error)
	GetUserById(UserID int) (model.User, error)
	CreateUser(model.User) (int, error)
}

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
		return model.User{}, ParsePostgresError(err)
	}

	return user, nil
}

func (repo *UserPostgresqlRepository) GetUserById(UserID int) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", usersTable)

	if err := repo.DB.Get(&user, query, UserID); err != nil {
		return model.User{}, ParsePostgresError(err)
	}

	return user, nil
}

func (repo *UserPostgresqlRepository) CreateUser(user model.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, role, password) VALUES ($1, $2, $3) RETURNING id", usersTable)

	row := repo.DB.QueryRow(query, user.Username, user.Role, user.Password)
	err := row.Scan(&id)
	if err != nil {
		return 0, ParsePostgresError(err)
	}

	return id, nil
}
