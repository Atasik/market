package repository

import (
	"fmt"
	"market/internal/model"
	"market/pkg/database/postgres"

	"github.com/jmoiron/sqlx"
)

type CartPostgresqlRepository struct {
	DB *sqlx.DB
}

func NewCartPostgresqlRepo(db *sqlx.DB) *CartPostgresqlRepository {
	return &CartPostgresqlRepository{DB: db}
}

func (repo *CartPostgresqlRepository) Create(userID int) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (user_id) VALUES ($1) RETURNING id", cartsTable)

	row := repo.DB.QueryRow(query, userID)
	err := row.Scan(&id)
	if err != nil {
		return 0, postgres.ParsePostgresError(err)
	}

	return id, nil
}

func (repo *CartPostgresqlRepository) AddProduct(cartID, productID, amount int) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (product_id, cart_id, purchased_amount) VALUES ($1, $2, $3) RETURNING id", productsCartsTable)

	row := repo.DB.QueryRow(query, productID, cartID, amount)
	err := row.Scan(&id)
	if err != nil {
		return 0, postgres.ParsePostgresError(err)
	}

	return id, nil
}

func (repo *CartPostgresqlRepository) GetByUserID(userID int) (model.Cart, error) {
	var Cart model.Cart
	query := fmt.Sprintf(`SELECT * FROM %s WHERE user_id = $1`, cartsTable)

	if err := repo.DB.Get(&Cart, query, userID); err != nil {
		return model.Cart{}, postgres.ParsePostgresError(err)
	}

	return Cart, nil
}

func (repo *CartPostgresqlRepository) GetAllProducts(cartID int, q model.ProductQueryInput) ([]model.Product, error) {
	var products []model.Product
	var limitValue string
	argId := 2
	args := make([]interface{}, 0)
	args = append(args, cartID)
	if q.Limit != 0 {
		limitValue = fmt.Sprintf("LIMIT $%d", argId)
		args = append(args, q.Limit)
		argId++
	}

	args = append(args, q.Offset)

	query := fmt.Sprintf(`SELECT p.id, p.user_id, p.title, p.price, p.tag, p.category, p.description, p.amount, pc.purchased_amount, p.created_at, p.updated_at, p.views, p.image_url FROM %s p 
			  INNER JOIN %s pc on pc.product_id = p.id
			  INNER JOIN %s c on pc.cart_id = c.id
			  WHERE c.id = $1 ORDER BY %s %s %s OFFSET $%d`, productsTable, productsCartsTable, cartsTable, q.SortBy, q.SortOrder, limitValue, argId)

	if err := repo.DB.Select(&products, query, args...); err != nil {
		return []model.Product{}, postgres.ParsePostgresError(err)
	}

	return products, nil
}

func (repo *CartPostgresqlRepository) GetProductByID(cartID, productID int) (model.Product, error) {
	var product model.Product
	query := fmt.Sprintf(`SELECT p.id, p.user_id, p.title, p.price, p.tag, p.category, p.description, p.amount, pc.purchased_amount, p.created_at, p.updated_at, p.views, p.image_url FROM %s p 
			  INNER JOIN %s pc on pc.product_id = p.id
			  INNER JOIN %s c on pc.cart_id = c.id
			  WHERE c.id = $1 AND p.id = $2`, productsTable, productsCartsTable, cartsTable)

	if err := repo.DB.Get(&product, query, cartID, productID); err != nil {
		return model.Product{}, postgres.ParsePostgresError(err)
	}

	return product, nil
}

func (repo *CartPostgresqlRepository) UpdateProductAmount(cartID, productID, amount int) error {
	query := fmt.Sprintf(`UPDATE %s SET purchased_amount = $1 WHERE cart_id = $2`, productsCartsTable)
	if _, err := repo.DB.Exec(query, amount, cartID); err != nil {
		return postgres.ParsePostgresError(err)
	}
	return nil
}

func (repo *CartPostgresqlRepository) DeleteProduct(cartID, productID int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE cart_id = $1 AND product_id = $2`, productsCartsTable)
	_, err := repo.DB.Exec(query, cartID, productID)
	if err != nil {
		return postgres.ParsePostgresError(err)
	}
	return nil
}

func (repo *CartPostgresqlRepository) DeleteAllProducts(cartID int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE cart_id = $1`, productsCartsTable)
	_, err := repo.DB.Exec(query, cartID)
	if err != nil {
		return postgres.ParsePostgresError(err)
	}
	return nil
}