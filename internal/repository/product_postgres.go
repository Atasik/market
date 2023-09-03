package repository

import (
	"fmt"
	"market/internal/model"
	"market/pkg/database/postgres"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ProductPostgresqlRepository struct {
	DB *sqlx.DB
}

func NewProductPostgresqlRepo(db *sqlx.DB) *ProductPostgresqlRepository {
	return &ProductPostgresqlRepository{DB: db}
}

func (repo *ProductPostgresqlRepository) Create(product model.Product) (int, error) {
	var productId int
	query := fmt.Sprintf("INSERT INTO %s (user_id, title, price, tag, category, description, amount, created_at, updated_at, views, image_url, image_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id", productsTable)

	row := repo.DB.QueryRow(query, product.UserID, product.Title, product.Price, product.Tag, product.Category, product.Description, product.Amount, product.CreatedAt, product.UpdatedAt, product.Views, product.ImageURL, product.ImageID)
	err := row.Scan(&productId)
	if err != nil {
		return 0, postgres.ParsePostgresError(err)
	}

	return productId, nil
}

func (repo *ProductPostgresqlRepository) GetAll(q model.ProductQueryInput) ([]model.Product, error) {
	var products []model.Product

	query := fmt.Sprintf("SELECT * FROM %s ORDER BY %s %s LIMIT $1 OFFSET $2", productsTable, q.SortBy, q.SortOrder)

	if err := repo.DB.Select(&products, query, q.Limit, q.Offset); err != nil {
		return nil, postgres.ParsePostgresError(err)
	}

	return products, nil
}

func (repo *ProductPostgresqlRepository) GetProductsByUserID(userID int, q model.ProductQueryInput) ([]model.Product, error) {
	var products []model.Product
	var setValue string
	argId := 2
	args := make([]interface{}, 0)
	args = append(args, userID)
	if q.ProductID != 0 {
		setValue = fmt.Sprintf("AND id!=$%d", argId)
		args = append(args, q.ProductID)
		argId++
	}

	args = append(args, q.Limit, q.Offset)

	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 %s ORDER BY %s %s LIMIT $%d OFFSET $%d", productsTable, setValue, q.SortBy, q.SortOrder, argId, argId+1)

	if err := repo.DB.Select(&products, query, args...); err != nil {
		return nil, postgres.ParsePostgresError(err)
	}

	return products, nil
}

func (repo *ProductPostgresqlRepository) GetByID(productID int) (model.Product, error) {
	var product model.Product
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", productsTable)

	if err := repo.DB.Get(&product, query, productID); err != nil {
		return model.Product{}, postgres.ParsePostgresError(err)
	}

	return product, nil
}

func (repo *ProductPostgresqlRepository) GetProductsByCategory(productCategory string, q model.ProductQueryInput) ([]model.Product, error) {
	var products []model.Product
	var setValue string
	argId := 2
	args := make([]interface{}, 0)
	args = append(args, productCategory)
	if q.ProductID != 0 {
		setValue = fmt.Sprintf("AND id!=$%d", argId)
		args = append(args, q.ProductID)
		argId++
	}

	args = append(args, q.Limit, q.Offset)

	query := fmt.Sprintf("SELECT * FROM %s WHERE category = $1 %s ORDER BY %s %s LIMIT $%d OFFSET $%d", productsTable, setValue, q.SortBy, q.SortOrder, argId, argId+1)

	if err := repo.DB.Select(&products, query, args...); err != nil {
		return nil, postgres.ParsePostgresError(err)
	}

	return products, nil
}

func (repo *ProductPostgresqlRepository) Update(productID int, input model.UpdateProductInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Price != nil {
		setValues = append(setValues, fmt.Sprintf("price=$%d", argId))
		args = append(args, *input.Price)
		argId++
	}

	if input.Tag != nil {
		setValues = append(setValues, fmt.Sprintf("tag=$%d", argId))
		args = append(args, *input.Tag)
		argId++
	}

	if input.Type != nil {
		setValues = append(setValues, fmt.Sprintf("category=$%d", argId))
		args = append(args, *input.Type)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.Amount != nil {
		setValues = append(setValues, fmt.Sprintf("amount=$%d", argId))
		args = append(args, *input.Amount)
		argId++
	}

	if input.Views != nil {
		setValues = append(setValues, "views=views+1")
	}

	if input.ImageURL != nil {
		setValues = append(setValues, fmt.Sprintf("image_url=$%d", argId))
		args = append(args, *input.ImageURL)
		argId++
	}

	if input.ImageURL != nil {
		setValues = append(setValues, fmt.Sprintf("image_id=$%d", argId))
		args = append(args, *input.ImageID)
		argId++
	}

	if input.UpdatedAt != nil {
		setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argId))
		args = append(args, *input.UpdatedAt)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id = $%d`, productsTable, setQuery, argId)
	args = append(args, productID)

	_, err := repo.DB.Exec(query, args...)
	if err != nil {
		return postgres.ParsePostgresError(err)
	}
	return nil
}

func (repo *ProductPostgresqlRepository) Delete(productId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", productsTable)

	_, err := repo.DB.Exec(query, productId)
	if err != nil {
		return postgres.ParsePostgresError(err)
	}
	return nil
}
