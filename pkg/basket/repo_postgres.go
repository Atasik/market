package basket

import (
	"market/pkg/product"

	"github.com/jmoiron/sqlx"
)

type BasketPostgresqlRepository struct {
	DB *sqlx.DB
}

func NewPostgresqlRepo(db *sqlx.DB) *BasketPostgresqlRepository {
	return &BasketPostgresqlRepository{DB: db}
}

func (repo *BasketPostgresqlRepository) AddProduct(userId, productId int) (int, error) {
	tx, err := repo.DB.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	query := "INSERT INTO products_users (product_id, user_id, purchased_count) VALUES ($1, $2, $3) RETURNING id"
	row := tx.QueryRow(query, productId, userId, 0)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (repo *BasketPostgresqlRepository) GetByID(userId int) ([]product.Product, error) {
	var products []product.Product
	query := `SELECT p.id, p.title, p.price, p.tag, p.type, p.description, p.count, p.creation_date, p.views, p.image_url FROM products p 
			  INNER JOIN products_users pu on pu.product_id = p.id
			  INNER JOIN users u on pu.user_id = u.id
			  WHERE u.id = $1`

	if err := repo.DB.Select(&products, query, userId); err != nil {
		return []product.Product{}, err
	}

	return products, nil
}

func (repo *BasketPostgresqlRepository) DeleteProduct(userId, productId int) (bool, error) {
	query := `DELETE FROM products_users WHERE user_id = $1 AND product_id = $2`
	_, err := repo.DB.Exec(query, userId, productId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo *BasketPostgresqlRepository) DeleteAll(userId int) (bool, error) {
	query := `DELETE FROM products_users WHERE user_id = $1`
	_, err := repo.DB.Exec(query, userId)
	if err != nil {
		return false, err
	}
	return true, nil
}
