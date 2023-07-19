package repository

import (
	"fmt"
	"market/pkg/model"

	"github.com/jmoiron/sqlx"
)

type CartRepo interface {
	CreateCart(userID int) (int, error)
	AddProduct(CartID, productID int) (int, error)
	GetByUserID(userID int) (model.Cart, error)
	GetProductByID(CartID, productID int) (model.Product, error)
	GetProducts(CartID int) ([]model.Product, error)
	DeleteProduct(CartID, productID int) (bool, error)
	DeleteAll(CartID int) (bool, error)
}

type CartPostgresqlRepository struct {
	DB *sqlx.DB
}

func NewCartPostgresqlRepo(db *sqlx.DB) *CartPostgresqlRepository {
	return &CartPostgresqlRepository{DB: db}
}

func (repo *CartPostgresqlRepository) CreateCart(userID int) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (user_id) VALUES ($1) RETURNING id", cartsTable)

	row := repo.DB.QueryRow(query, userID)
	err := row.Scan(&id)
	if err != nil {
		return 0, ParsePostgresError(err)
	}

	return id, nil
}

// проверка, что есть права
func (repo *CartPostgresqlRepository) AddProduct(CartID, productID int) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (product_id, cart_id, purchased_amount) VALUES ($1, $2, $3) RETURNING id", productsCartsTable)

	row := repo.DB.QueryRow(query, productID, CartID, 0)
	err := row.Scan(&id)
	if err != nil {
		return 0, ParsePostgresError(err)
	}

	return id, nil
}

func (repo *CartPostgresqlRepository) GetByUserID(userID int) (model.Cart, error) {
	var Cart model.Cart
	query := fmt.Sprintf(`SELECT * FROM %s WHERE user_id = $1`, cartsTable)

	if err := repo.DB.Get(&Cart, query, userID); err != nil {
		return model.Cart{}, ParsePostgresError(err)
	}

	return Cart, nil
}

// проверка, что есть права
func (repo *CartPostgresqlRepository) GetProducts(CartID int) ([]model.Product, error) {
	var products []model.Product
	query := fmt.Sprintf(`SELECT p.id, p.user_id, p.title, p.price, p.tag, p.category, p.description, p.amount, p.created_at, p.updated_at, p.views, p.image_url FROM %s p 
			  INNER JOIN %s pb on pb.product_id = p.id
			  INNER JOIN %s b on pb.cart_id = b.id
			  WHERE b.id = $1`, productsTable, productsCartsTable, cartsTable)

	if err := repo.DB.Select(&products, query, CartID); err != nil {
		return []model.Product{}, ParsePostgresError(err)
	}

	return products, nil
}

func (repo *CartPostgresqlRepository) GetProductByID(CartID, productID int) (model.Product, error) {
	var product model.Product
	query := fmt.Sprintf(`SELECT p.id, p.user_id, p.title, p.price, p.tag, p.category, p.description, p.amount, p.created_at, p.updated_at, p.views, p.image_url FROM %s p 
			  INNER JOIN %s pb on pb.product_id = p.id
			  INNER JOIN %s b on pb.cart_id = b.id
			  WHERE b.id = $1 AND p.id = $2`, productsTable, productsCartsTable, cartsTable)

	if err := repo.DB.Get(&product, query, CartID, productID); err != nil {
		return model.Product{}, ParsePostgresError(err)
	}

	return product, nil
}

// проверка, что есть права
func (repo *CartPostgresqlRepository) DeleteProduct(CartID, productID int) (bool, error) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE cart_id = $1 AND product_id = $2`, productsCartsTable)
	_, err := repo.DB.Exec(query, CartID, productID)
	if err != nil {
		return false, ParsePostgresError(err)
	}
	return true, nil
}

// проверка, что есть права
func (repo *CartPostgresqlRepository) DeleteAll(CartID int) (bool, error) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE cart_id = $1`, productsCartsTable)
	_, err := repo.DB.Exec(query, CartID)
	if err != nil {
		return false, ParsePostgresError(err)
	}
	return true, nil
}
