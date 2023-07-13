package repository

import (
	"fmt"
	"market/pkg/model"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ProductRepo interface {
	GetAll(orderBy string) ([]model.Product, error)
	GetByID(productId int) (model.Product, error)
	Create(product model.Product) (int, error)
	Update(productId int, input model.UpdateProductInput) (bool, error)
	Delete(productId int) (bool, error)
	GetByType(productType string, limit int) ([]model.Product, error)
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
		return nil, err
	}

	return products, nil
}

func (repo *ProductPostgresqlRepository) GetByID(productId int) (model.Product, error) {
	var product model.Product
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", productsTable)

	if err := repo.DB.Get(&product, query, productId); err != nil {
		return model.Product{}, err
	}

	return product, nil
}

// проверка, что есть права
func (repo *ProductPostgresqlRepository) Create(product model.Product) (int, error) {
	var productId int
	query := fmt.Sprintf("INSERT INTO %s (title, price, tag, type, description, count, creation_date, views, image_url, image_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id", productsTable)

	row := repo.DB.QueryRow(query, product.Title, product.Price, product.Tag, product.Type, product.Description, product.Count, product.CreationDate, product.Views, product.ImageURL, product.ImageID)
	err := row.Scan(&productId)
	if err != nil {
		return 0, err
	}

	return productId, nil
}

// проверка, что есть права
func (repo *ProductPostgresqlRepository) Update(productId int, input model.UpdateProductInput) (bool, error) {
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
		setValues = append(setValues, fmt.Sprintf("type=$%d", argId))
		args = append(args, *input.Type)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.Count != nil {
		setValues = append(setValues, fmt.Sprintf("count=$%d", argId))
		args = append(args, *input.Count)
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

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id = $%d`, productsTable, setQuery, argId)
	args = append(args, productId)

	_, err := repo.DB.Exec(query, args...)
	if err != nil {
		return false, err
	}
	return true, nil
}

// проверка, что есть права
func (repo *ProductPostgresqlRepository) Delete(productId int) (bool, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", productsTable)

	_, err := repo.DB.Exec(query, productId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo *ProductPostgresqlRepository) GetByType(productType string, limit int) ([]model.Product, error) {
	var products []model.Product
	query := fmt.Sprintf("SELECT * FROM %s WHERE type = $1", productsTable)

	if err := repo.DB.Select(&products, query, productType); err != nil {
		return nil, err
	}

	return products, nil
}
