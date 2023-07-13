package repository

import (
	"fmt"
	"market/pkg/model"

	"github.com/jmoiron/sqlx"
)

type BasketRepo interface {
	AddProduct(userId, productId int) (int, error)
	GetByID(userId int) ([]model.Product, error)
	DeleteProduct(userId, productId int) (bool, error)
	DeleteAll(userId int) (bool, error)
}

type BasketPostgresqlRepository struct {
	DB *sqlx.DB
}

func NewBasketPostgresqlRepo(db *sqlx.DB) *BasketPostgresqlRepository {
	return &BasketPostgresqlRepository{DB: db}
}

// проверка, что есть права
func (repo *BasketPostgresqlRepository) AddProduct(userId, productId int) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (product_id, user_id, purchased_count) VALUES ($1, $2, $3) RETURNING id", ProductsUsersTable)

	row := repo.DB.QueryRow(query, productId, userId, 0)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// проверка, что есть права
func (repo *BasketPostgresqlRepository) GetByID(userId int) ([]model.Product, error) {
	var products []model.Product
	query := fmt.Sprintf(`SELECT p.id, p.title, p.price, p.tag, p.type, p.description, p.count, p.creation_date, p.views, p.image_url FROM %s p 
			  INNER JOIN %s pu on pu.product_id = p.id
			  INNER JOIN %s u on pu.user_id = u.id
			  WHERE u.id = $1`, productsTable, ProductsUsersTable, usersTable)

	if err := repo.DB.Select(&products, query, userId); err != nil {
		return []model.Product{}, err
	}

	return products, nil
}

// проверка, что есть права
func (repo *BasketPostgresqlRepository) DeleteProduct(userId, productId int) (bool, error) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id = $1 AND product_id = $2`, ProductsUsersTable)
	_, err := repo.DB.Exec(query, userId, productId)
	if err != nil {
		return false, err
	}
	return true, nil
}

// проверка, что есть права
func (repo *BasketPostgresqlRepository) DeleteAll(userId int) (bool, error) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id = $1`, ProductsUsersTable)
	_, err := repo.DB.Exec(query, userId)
	if err != nil {
		return false, err
	}
	return true, nil
}
