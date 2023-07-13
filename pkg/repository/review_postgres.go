package repository

import (
	"fmt"
	"market/pkg/model"

	"github.com/jmoiron/sqlx"
)

type ReviewRepo interface {
	Create(review model.Review) (int, error)
	Delete(reviewID int) (bool, error)
	Update(userID, productID int, text string) (bool, error)
	GetAll(productID int, orderBy string) ([]model.Review, error)
}

type ReviewPostgresqlRepository struct {
	DB *sqlx.DB
}

func NewReviewPostgresqlRepo(db *sqlx.DB) *ReviewPostgresqlRepository {
	return &ReviewPostgresqlRepository{DB: db}
}

func (repo *ReviewPostgresqlRepository) Create(review model.Review) (int, error) {
	var reviewId int
	query := fmt.Sprintf("INSERT INTO %s (creation_date, product_id, user_id, username, review_text, rating) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id", reviewsTable)

	row := repo.DB.QueryRow(query, review.CreationDate, review.ProductID, review.UserID, review.Username, review.Text, review.Rating)
	err := row.Scan(&reviewId)
	if err != nil {
		return 0, err
	}

	return reviewId, nil
}

func (repo *ReviewPostgresqlRepository) Delete(reviewID int) (bool, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", reviewsTable)

	_, err := repo.DB.Exec(query, reviewID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo *ReviewPostgresqlRepository) Update(userID, productID int, text string) (bool, error) {
	query := fmt.Sprintf("UPDATE %s SET review_text = $1 WHERE (user_id = $2 AND product_id = $3)", reviewsTable)

	_, err := repo.DB.Exec(query, text, userID, productID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo *ReviewPostgresqlRepository) GetAll(productID int, orderBy string) ([]model.Review, error) {
	var rewiews []model.Review
	query := fmt.Sprintf("SELECT * FROM %s WHERE product_id = $1", reviewsTable)

	if err := repo.DB.Select(&rewiews, query, productID); err != nil {
		return nil, err
	}

	return rewiews, nil
}
