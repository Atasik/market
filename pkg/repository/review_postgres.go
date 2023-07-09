package repository

import (
	"market/pkg/model"

	"github.com/jmoiron/sqlx"
)

type ReviewRepo interface {
	Add(userID, productID, rating int, username, text string) (model.Review, error)
	Delete(userID, productID int) (bool, error)
	GetAll(productID int, orderBy string) ([]model.Review, error)
}

type ReviewPostgresqlRepository struct {
	DB *sqlx.DB
}

func NewReviewPostgresqlRepo(db *sqlx.DB) *ReviewPostgresqlRepository {
	return &ReviewPostgresqlRepository{DB: db}
}

func (repo *ReviewPostgresqlRepository) Add(userID, productID, rating int, username, text string) (model.Review, error) {
	return model.Review{}, nil
}

func (repo *ReviewPostgresqlRepository) Delete(userID, productID int) (bool, error) {
	return false, nil
}

func (repo *ReviewPostgresqlRepository) GetAll(productID int, orderBy string) ([]model.Review, error) {
	return nil, nil
}
