package repository

import (
	"fmt"
	"market/internal/model"
	"market/pkg/database/postgres"

	"github.com/jmoiron/sqlx"
)

type UserPostgresqlRepository struct {
	db *sqlx.DB
}

func NewUserPostgresqlRepo(db *sqlx.DB) *UserPostgresqlRepository {
	return &UserPostgresqlRepository{db: db}
}

func (repo *UserPostgresqlRepository) GetUser(login string) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE username = $1", usersTable)

	if err := repo.db.Get(&user, query, login); err != nil {
		return model.User{}, postgres.ParsePostgresError(err)
	}
	return user, nil
}

func (repo *UserPostgresqlRepository) GetUserByID(userID int) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", usersTable)

	if err := repo.db.Get(&user, query, userID); err != nil {
		return model.User{}, postgres.ParsePostgresError(err)
	}
	return user, nil
}

func (repo *UserPostgresqlRepository) CreateUser(user model.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, role, password) VALUES ($1, $2, $3) RETURNING id", usersTable)

	row := repo.db.QueryRow(query, user.Username, user.Role, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, postgres.ParsePostgresError(err)
	}
	return id, nil
}
