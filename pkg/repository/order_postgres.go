package repository

import (
	"fmt"
	"market/pkg/model"

	"github.com/jmoiron/sqlx"
)

type OrderRepo interface {
	Create(cartID, userID int, order model.Order, products []model.Product) (int, error)
	GetAll(userID int) ([]model.Order, error)
	GetByID(orderID int) (model.Order, error)
	GetProductsByOrderID(orderID int) ([]model.Product, error)
}

type OrderPostgresqlRepository struct {
	DB *sqlx.DB
}

func NewOrderPostgresqlRepo(db *sqlx.DB) *OrderPostgresqlRepository {
	return &OrderPostgresqlRepository{DB: db}
}

func (repo *OrderPostgresqlRepository) Create(cartID, userID int, order model.Order, products []model.Product) (int, error) {
	tx, err := repo.DB.Begin()
	if err != nil {
		return 0, ParsePostgresError(err)
	}

	query := fmt.Sprintf("INSERT INTO %s (created_at, delivered_at, user_id) VALUES ($1, $2, $3) RETURNING id", ordersTable)
	row := tx.QueryRow(query, order.CreatedAt, order.DeliveredAt, userID)
	if err := row.Scan(&order.ID); err != nil {
		tx.Rollback()
		return 0, ParsePostgresError(err)
	}

	query = fmt.Sprintf("INSERT INTO %s (order_id, product_id) VALUES ($1, $2)", productsOrdersTable)
	for _, product := range products {
		if _, err := tx.Exec(query, order.ID, product.ID); err != nil {
			tx.Rollback()
			return 0, ParsePostgresError(err)
		}
	}

	query = fmt.Sprintf(`DELETE FROM %s WHERE cart_id = $1`, productsCartsTable)
	if _, err = tx.Exec(query, cartID); err != nil {
		tx.Rollback()
		return 0, ParsePostgresError(err)
	}

	return order.ID, ParsePostgresError(tx.Commit())
}

func (repo *OrderPostgresqlRepository) GetAll(userID int) ([]model.Order, error) {
	var orders []model.Order
	query := fmt.Sprintf(`SELECT o.id, o.created_at, o.delivered_at FROM %s o
			  INNER JOIN %s u on o.user_id = u.id
			  WHERE u.id = $1`, ordersTable, usersTable)

	if err := repo.DB.Select(&orders, query, userID); err != nil {
		return []model.Order{}, ParsePostgresError(err)
	}

	return orders, nil
}

func (repo *OrderPostgresqlRepository) GetByID(orderID int) (model.Order, error) {
	var order model.Order
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", ordersTable)

	if err := repo.DB.Get(&order, query, orderID); err != nil {
		return model.Order{}, ParsePostgresError(err)
	}

	return order, nil
}

func (repo *OrderPostgresqlRepository) GetProductsByOrderID(orderID int) ([]model.Product, error) {
	var products []model.Product
	query := fmt.Sprintf(`SELECT p.id, p.user_id, p.title, p.price, p.tag, p.category, p.description, p.amount, p.created_at, p.updated_at, p.views, p.image_url FROM %s p 
			  INNER JOIN %s po on po.product_id = p.id
			  INNER JOIN %s o on po.order_id = o.id
			  WHERE o.id = $1`, productsTable, productsOrdersTable, ordersTable)

	if err := repo.DB.Select(&products, query, orderID); err != nil {
		return []model.Product{}, ParsePostgresError(err)
	}

	return products, nil
}
