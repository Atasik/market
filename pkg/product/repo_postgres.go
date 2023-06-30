package product

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ProductPostgresqlRepository struct {
	DB *sqlx.DB
}

func NewPostgresqlRepo(db *sqlx.DB) *ProductPostgresqlRepository {
	return &ProductPostgresqlRepository{DB: db}
}

func (repo *ProductPostgresqlRepository) GetAll(orderBy string) ([]Product, error) {
	var products []Product
	var setValue string

	if orderBy != "" {
		setValue = fmt.Sprintf("ORDER BY %s DESC", orderBy)
	}

	query := fmt.Sprintf("SELECT * FROM products %s", setValue)

	if err := repo.DB.Select(&products, query); err != nil {
		print(err.Error())
		return nil, err
	}

	return products, nil
}

func (repo *ProductPostgresqlRepository) GetByID(productId int) (Product, error) {
	var product Product
	query := "SELECT * FROM products WHERE id = $1"

	if err := repo.DB.Get(&product, query, productId); err != nil {
		return Product{}, err
	}

	return product, nil
}

func (repo *ProductPostgresqlRepository) Create(product Product) (int, error) {
	var productId int
	query := "INSERT INTO products (title, price, tag, type, description, count, creation_date, views, image_url) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id"

	row := repo.DB.QueryRow(query, product.Title, product.Price, product.Tag, product.Type, product.Description, product.Count, product.CreationDate, product.Views, product.ImageURL)
	err := row.Scan(&productId)
	if err != nil {
		return 0, err
	}

	return productId, nil
}

func (repo *ProductPostgresqlRepository) Update(productId int, input UpdateProductInput) (bool, error) {
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

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE products SET %s WHERE id = $%d`, setQuery, argId)
	args = append(args, productId)

	_, err := repo.DB.Exec(query, args...)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo *ProductPostgresqlRepository) Delete(productId int) (bool, error) {
	query := "DELETE FROM products WHERE id = $1"

	_, err := repo.DB.Exec(query, productId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo *ProductPostgresqlRepository) GetByType(productType string, limit int) ([]Product, error) {
	var products []Product
	query := "SELECT * FROM products WHERE type = $1"

	if err := repo.DB.Select(&products, query, productType); err != nil {
		print(err.Error())
		return nil, err
	}

	return products, nil
}
