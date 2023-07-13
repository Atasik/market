package repository

import (
	"fmt"
	"market/pkg/model"

	"github.com/jmoiron/sqlx"
)

type BasketRepo interface {
	CreateBasket(userID int) (int, error)
	AddProduct(basketID, productID int) (int, error)
	GetByUserID(userID int) (model.Basket, error)
	GetProducts(basketID int) ([]model.Product, error)
	DeleteProduct(basketID, productID int) (bool, error)
	DeleteAll(basketID int) (bool, error)
}

type BasketPostgresqlRepository struct {
	DB *sqlx.DB
}

func NewBasketPostgresqlRepo(db *sqlx.DB) *BasketPostgresqlRepository {
	return &BasketPostgresqlRepository{DB: db}
}

func (repo *BasketPostgresqlRepository) CreateBasket(userID int) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (user_id) VALUES ($1) RETURNING id", basketsTable)

	row := repo.DB.QueryRow(query, userID)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// проверка, что есть права
func (repo *BasketPostgresqlRepository) AddProduct(basketID, productID int) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (product_id, basket_id, purchased_count) VALUES ($1, $2, $3) RETURNING id", ProductsBasketsTable)

	row := repo.DB.QueryRow(query, productID, basketID, 0)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *BasketPostgresqlRepository) GetByUserID(userID int) (model.Basket, error) {
	var basket model.Basket
	query := fmt.Sprintf(`SELECT * FROM %s WHERE user_id = $1`, basketsTable)

	if err := repo.DB.Get(&basket, query, userID); err != nil {
		return model.Basket{}, err
	}

	return basket, nil
}

// проверка, что есть права
func (repo *BasketPostgresqlRepository) GetProducts(basketID int) ([]model.Product, error) {
	var products []model.Product
	query := fmt.Sprintf(`SELECT p.id, p.title, p.price, p.tag, p.type, p.description, p.count, p.creation_date, p.views, p.image_url FROM %s p 
			  INNER JOIN %s pb on pb.product_id = p.id
			  INNER JOIN %s b on pb.basket_id = b.id
			  WHERE b.id = $1`, productsTable, ProductsBasketsTable, basketsTable)

	if err := repo.DB.Select(&products, query, basketID); err != nil {
		return []model.Product{}, err
	}

	return products, nil
}

// проверка, что есть права
func (repo *BasketPostgresqlRepository) DeleteProduct(basketID, productID int) (bool, error) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE basket_id = $1 AND product_id = $2`, ProductsBasketsTable)
	_, err := repo.DB.Exec(query, basketID, productID)
	if err != nil {
		return false, err
	}
	return true, nil
}

// проверка, что есть права
func (repo *BasketPostgresqlRepository) DeleteAll(basketID int) (bool, error) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE basket_id = $1`, ProductsBasketsTable)
	_, err := repo.DB.Exec(query, basketID)
	if err != nil {
		return false, err
	}
	return true, nil
}
