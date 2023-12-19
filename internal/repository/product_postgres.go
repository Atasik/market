package repository

import (
	"fmt"
	"market/internal/model"
	"market/pkg/database/postgres"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ProductPostgresqlRepository struct {
	db *sqlx.DB
}

func NewProductPostgresqlRepo(db *sqlx.DB) *ProductPostgresqlRepository {
	return &ProductPostgresqlRepository{db: db}
}

func (repo *ProductPostgresqlRepository) Create(product model.Product) (int, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return 0, postgres.ParsePostgresError(err)
	}

	fmt.Println("goyda")
	fmt.Println(product)
	var productID int
	query := "INSERT INTO products (user_id, title, price, category, description, amount, created_at, updated_at, views, image_url, image_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id"

	row := tx.QueryRow(query, product.UserID, product.Title, product.Price, product.Category, product.Description, product.Amount, product.CreatedAt, product.UpdatedAt, product.Views, product.ImageURL, product.ImageID)
	if err = row.Scan(&productID); err != nil {
		return 0, postgres.ParsePostgresError(err)
	}

	for i := 0; i < len(product.Tags); i++ {
		var tagId int
		query = "INSERT INTO tags (name) VALUES ($1) ON CONFLICT(id) DO NOTHING RETURNING id"
		row := tx.QueryRow(query, product.Tags[i].Name)
		if err = row.Scan(&tagId); err != nil {
			return 0, postgres.ParsePostgresError(err)
		}

		query = "INSERT INTO products_tags (product_id, tag_id) VALUES($1, $2)"
		if _, err = tx.Exec(query, productID, tagId); err != nil {
			return 0, postgres.ParsePostgresError(err)
		}
	}

	return productID, postgres.ParsePostgresError(tx.Commit())
}

func (repo *ProductPostgresqlRepository) GetAll(q model.ProductQueryInput) ([]model.Product, error) {
	var products []model.Product

	query := fmt.Sprintf(`SELECT
    						p.id,
    						p.user_id,
    						p.title,
    						p.price,
    						json_agg(jsonb_build_object('name', t.name)) as tags,
    						p.category,
    						p.description,
    						p.amount,
   							p.created_at,
    						p.updated_at,
    						p.views,
						    p.image_url
						FROM
						    products p
						    LEFT JOIN products_tags pt ON pt.product_id = p.id
						    LEFT JOIN tags t ON pt.tag_id = t.id
						GROUP BY
						    p.id, p.user_id, p.title, p.price, p.category, p.description, p.amount, p.created_at, p.updated_at, p.views, p.image_url
						ORDER BY
						    %s %s
						LIMIT
							$1 
						OFFSET 
							$2`, q.SortBy, q.SortOrder)

	if err := repo.db.Select(&products, query, q.Limit, q.Offset); err != nil {
		return nil, postgres.ParsePostgresError(err)
	}
	return products, nil
}

func (repo *ProductPostgresqlRepository) GetProductsByUserID(userID int, q model.ProductQueryInput) ([]model.Product, error) {
	var products []model.Product
	var setValue string
	argID := 2
	args := make([]interface{}, 0)
	args = append(args, userID)
	if q.ProductID != 0 {
		setValue = fmt.Sprintf("AND id!=$%d", argID)
		args = append(args, q.ProductID)
		argID++
	}

	args = append(args, q.Limit, q.Offset)

	query := fmt.Sprintf(`SELECT
							p.id,
							p.user_id,
							p.title,
							p.price,
							json_agg(jsonb_build_object('name', t.name)) as tags,
							p.category,
							p.description,
							p.amount,
							p.created_at,
							p.updated_at,
							p.views,
							p.image_url
						FROM
							products p
							LEFT JOIN products_tags pt ON pt.product_id = p.id
							LEFT JOIN tags t ON pt.tag_id = t.id
						WHERE 
							p.user_id = $1 %s
						GROUP BY
						    p.id, p.user_id, p.title, p.price, p.category, p.description, p.amount, p.created_at, p.updated_at, p.views, p.image_url
					    ORDER BY 
							%s %s 
						LIMIT 
							$%d 
						OFFSET 
							$%d`, setValue, q.SortBy, q.SortOrder, argID, argID+1)

	if err := repo.db.Select(&products, query, args...); err != nil {
		return nil, postgres.ParsePostgresError(err)
	}
	return products, nil
}

func (repo *ProductPostgresqlRepository) GetByID(productID int) (model.Product, error) {
	var product model.Product
	query := `SELECT 
				  p.id,
				  p.user_id,
				  p.title,
				  p.price,
				  json_agg(jsonb_build_object('name', t.name)) as tags,
				  p.category,
				  p.description,
				  p.amount,
				  p.created_at,
				  p.updated_at,
				  p.views,
				  p.image_url
			  FROM
				  products p
				  LEFT JOIN products_tags pt ON pt.product_id = p.id
				  LEFT JOIN tags t ON pt.tag_id = t.id
			  WHERE 
				  p.id = $1
			  GROUP BY
				  p.id, p.user_id, p.title, p.price, p.category, p.description, p.amount, p.created_at, p.updated_at, p.views, p.image_url`

	if err := repo.db.Get(&product, query, productID); err != nil {
		return model.Product{}, postgres.ParsePostgresError(err)
	}
	return product, nil
}

func (repo *ProductPostgresqlRepository) GetProductsByCategory(productCategory string, q model.ProductQueryInput) ([]model.Product, error) {
	var products []model.Product
	var setValue string
	argID := 2
	args := make([]interface{}, 0)
	args = append(args, productCategory)
	if q.ProductID != 0 {
		setValue = fmt.Sprintf("AND p.id!=$%d", argID)
		args = append(args, q.ProductID)
		argID++
	}

	args = append(args, q.Limit, q.Offset)

	query := fmt.Sprintf(`SELECT 
							p.id,
							p.user_id,
							p.title,
							p.price,
							json_agg(jsonb_build_object('name', t.name)) as tags,
							p.category,
							p.description,
							p.amount,
							p.created_at,
							p.updated_at,
							p.views,
							p.image_url
						FROM
							products p
							LEFT JOIN products_tags pt ON pt.product_id = p.id
							LEFT JOIN tags t ON pt.tag_id = t.id
						WHERE 
							category = $1 %s 
						GROUP BY
							p.id, p.user_id, p.title, p.price, p.category, p.description, p.amount, p.created_at, p.updated_at, p.views, p.image_url
						ORDER BY 
							%s %s
						LIMIT 
							$%d 
						OFFSET 
							$%d`, setValue, q.SortBy, q.SortOrder, argID, argID+1)

	if err := repo.db.Select(&products, query, args...); err != nil {
		return nil, postgres.ParsePostgresError(err)
	}
	return products, nil
}

func (repo *ProductPostgresqlRepository) Update(productID int, input model.UpdateProductInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argID))
		args = append(args, *input.Title)
		argID++
	}

	if input.Price != nil {
		setValues = append(setValues, fmt.Sprintf("price=$%d", argID))
		args = append(args, *input.Price)
		argID++
	}

	if input.Type != nil {
		setValues = append(setValues, fmt.Sprintf("category=$%d", argID))
		args = append(args, *input.Type)
		argID++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argID))
		args = append(args, *input.Description)
		argID++
	}

	if input.Amount != nil {
		setValues = append(setValues, fmt.Sprintf("amount=$%d", argID))
		args = append(args, *input.Amount)
		argID++
	}

	if input.Views != nil {
		setValues = append(setValues, "views=views+1")
	}

	if input.ImageURL != nil {
		setValues = append(setValues, fmt.Sprintf("image_url=$%d", argID))
		args = append(args, *input.ImageURL)
		argID++
	}

	if input.ImageURL != nil {
		setValues = append(setValues, fmt.Sprintf("image_id=$%d", argID))
		args = append(args, *input.ImageID)
		argID++
	}

	if input.UpdatedAt != nil {
		setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argID))
		args = append(args, *input.UpdatedAt)
		argID++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE products SET %s WHERE id = $%d`, setQuery, argID)
	args = append(args, productID)

	if _, err := repo.db.Exec(query, args...); err != nil {
		return postgres.ParsePostgresError(err)
	}
	return nil
}

func (repo *ProductPostgresqlRepository) Delete(productID int) error {
	query := "DELETE FROM products WHERE id = $1"

	if _, err := repo.db.Exec(query, productID); err != nil {
		return postgres.ParsePostgresError(err)
	}
	return nil
}
