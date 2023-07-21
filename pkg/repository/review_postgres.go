package repository

import (
	"fmt"
	"market/pkg/model"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ReviewRepo interface {
	Create(review model.Review) (int, error)
	Delete(reviewID int) error
	Update(userID, productID int, input model.UpdateReviewInput) error
	GetAll(productID int) ([]model.Review, error)
	GetReviewIDByProductIDUserID(productID, userID int) (int, error)
}

type ReviewPostgresqlRepository struct {
	DB *sqlx.DB
}

func NewReviewPostgresqlRepo(db *sqlx.DB) *ReviewPostgresqlRepository {
	return &ReviewPostgresqlRepository{DB: db}
}

func (repo *ReviewPostgresqlRepository) Create(review model.Review) (int, error) {
	var reviewId int
	query := fmt.Sprintf("INSERT INTO %s (created_at, updated_at, product_id, user_id, text, category) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id", reviewsTable)

	row := repo.DB.QueryRow(query, review.CreatedAt, review.UpdatedAt, review.ProductID, review.UserID, review.Text, review.Category)
	err := row.Scan(&reviewId)
	if err != nil {
		return 0, ParsePostgresError(err)
	}

	return reviewId, nil
}

// проверка, что есть права
func (repo *ReviewPostgresqlRepository) Delete(reviewID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", reviewsTable)

	_, err := repo.DB.Exec(query, reviewID)
	if err != nil {
		return ParsePostgresError(err)
	}
	return nil
}

// проверка, что есть права
func (repo *ReviewPostgresqlRepository) Update(userID, productID int, input model.UpdateReviewInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	if input.Text != nil {
		setValues = append(setValues, fmt.Sprintf("text=$%d", argId))
		args = append(args, *input.Text)
		argId++
	}

	if input.Category != nil {
		setValues = append(setValues, fmt.Sprintf("category=$%d", argId))
		args = append(args, input.Category)
		argId++
	}

	if input.UpdatedAt != nil {
		setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argId))
		args = append(args, input.UpdatedAt)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE (user_id = $%d AND product_id = $%d)", reviewsTable, setQuery, argId, argId+1)
	args = append(args, userID, productID)

	_, err := repo.DB.Exec(query, args...)
	if err != nil {
		return ParsePostgresError(err)
	}
	return nil
}

func (repo *ReviewPostgresqlRepository) GetAll(productID int) ([]model.Review, error) {
	var rewiews []model.Review
	query := fmt.Sprintf("SELECT * FROM %s WHERE product_id = $1", reviewsTable)

	if err := repo.DB.Select(&rewiews, query, productID); err != nil {
		return nil, ParsePostgresError(err)
	}

	return rewiews, nil
}

func (repo *ReviewPostgresqlRepository) GetReviewIDByProductIDUserID(productID, userID int) (int, error) {
	var id int
	query := fmt.Sprintf("SELECT id FROM %s WHERE product_id = $1 AND user_id = $2", reviewsTable)

	if err := repo.DB.Get(&id, query, productID, userID); err != nil {
		return 0, ParsePostgresError(err)
	}

	return id, nil
}
