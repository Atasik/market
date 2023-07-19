package repository

import (
	"fmt"
	"market/pkg/model"

	"github.com/jmoiron/sqlx"
)

type OrderRepo interface {
	Create(userID int, order model.Order, products []model.Product) (int, error)
	GetAll(userID int) ([]model.Order, error)
	GetByID(orderID int) (model.Order, error)
}

type OrderPostgresqlRepository struct {
	DB *sqlx.DB
}

func NewOrderPostgresqlRepo(db *sqlx.DB) *OrderPostgresqlRepository {
	return &OrderPostgresqlRepository{DB: db}
}

// проверка, что есть права
func (repo *OrderPostgresqlRepository) Create(userID int, order model.Order, products []model.Product) (int, error) {
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

	query = fmt.Sprintf("INSERT INTO %s (order_id, product_id) VALUES ($1, $2)", ProductsOrdersTable)
	for _, product := range products {
		if _, err := tx.Exec(query, order.ID, product.ID); err != nil {
			tx.Rollback()
			return 0, ParsePostgresError(err)
		}
	}

	return order.ID, ParsePostgresError(tx.Commit())
}

// проверка, что есть права
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

// проверка, что есть права
func (repo *OrderPostgresqlRepository) GetByID(orderID int) (model.Order, error) {
	var order model.Order
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", ordersTable)

	if err := repo.DB.Get(&order, query, orderID); err != nil {
		return model.Order{}, ParsePostgresError(err)
	}

	return order, nil
}
