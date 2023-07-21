package repository

import (
	"fmt"
	"market/pkg/model"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ProductRepo interface {
	GetAll(orderBy string) ([]model.Product, error)
	GetByID(productID int) (model.Product, error)
	Create(product model.Product) (int, error)
	Update(productID int, input model.UpdateProductInput) error
	Delete(productID int) error
	GetByType(productType string, productID, limit int) ([]model.Product, error)
}

type ProductPostgresqlRepository struct {
	DB *sqlx.DB
}

func NewProductPostgresqlRepo(db *sqlx.DB) *ProductPostgresqlRepository {
	return &ProductPostgresqlRepository{DB: db}
}

func (repo *ProductPostgresqlRepository) GetAll(orderBy string) ([]model.Product, error) {
	var products []model.Product
	var setValue string

	if orderBy != "" {
		setValue = fmt.Sprintf("ORDER BY %s DESC", orderBy)
	}

	query := fmt.Sprintf("SELECT * FROM %s %s", productsTable, setValue)

	if err := repo.DB.Select(&products, query); err != nil {
		return nil, ParsePostgresError(err)
	}

	return products, nil
}

func (repo *ProductPostgresqlRepository) GetByID(productID int) (model.Product, error) {
	var product model.Product
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", productsTable)

	if err := repo.DB.Get(&product, query, productID); err != nil {
		return model.Product{}, ParsePostgresError(err)
	}

	return product, nil
}

// проверка, что есть права
func (repo *ProductPostgresqlRepository) Create(product model.Product) (int, error) {
	var productId int
	query := fmt.Sprintf("INSERT INTO %s (user_id, title, price, tag, category, description, amount, created_at, updated_at, views, image_url, image_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id", productsTable)

	row := repo.DB.QueryRow(query, product.UserID, product.Title, product.Price, product.Tag, product.Category, product.Description, product.Amount, product.CreatedAt, product.UpdatedAt, product.Views, product.ImageURL, product.ImageID)
	err := row.Scan(&productId)
	if err != nil {
		return 0, ParsePostgresError(err)
	}

	return productId, nil
}

// проверка, что есть права
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
		return ParsePostgresError(err)
	}
	return nil
}

// проверка, что есть права
func (repo *ProductPostgresqlRepository) Delete(productId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", productsTable)

	_, err := repo.DB.Exec(query, productId)
	if err != nil {
		return ParsePostgresError(err)
	}
	return nil
}

func (repo *ProductPostgresqlRepository) GetByType(productType string, productID, limit int) ([]model.Product, error) {
	var products []model.Product
	var setValue string
	argId := 2
	if productID != 0 {
		setValue = fmt.Sprintf("AND id!=$%d", argId)
		argId++
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE category = $1 %s LIMIT $%d", productsTable, setValue, argId)

	if err := repo.DB.Select(&products, query, productType, productID, limit); err != nil {
		return nil, ParsePostgresError(err)
	}

	return products, nil
}
