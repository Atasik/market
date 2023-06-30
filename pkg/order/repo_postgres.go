package order

import (
	"market/pkg/product"

	"github.com/jmoiron/sqlx"
)

type OrderPostgresqlRepository struct {
	DB *sqlx.DB
}

func NewPostgresqlRepo(db *sqlx.DB) *OrderPostgresqlRepository {
	return &OrderPostgresqlRepository{DB: db}
}

func (repo *OrderPostgresqlRepository) Create(userID int, order Order, products []product.Product) (int, error) {
	tx, err := repo.DB.Begin()
	if err != nil {
		return 0, err
	}

	query := "INSERT INTO orders (creation_date, delivery_date) VALUES ($1, $2) RETURNING id"
	row := tx.QueryRow(query, order.CreationDate, order.DeliveryDate)
	if err := row.Scan(&order.ID); err != nil {
		tx.Rollback()
		return 0, err
	}

	query = "INSERT INTO orders_users (order_id, user_id) VALUES ($1, $2)"
	if _, err := tx.Exec(query, order.ID, userID); err != nil {
		tx.Rollback()
		return 0, err
	}

	query = "INSERT INTO products_orders (order_id, product_id) VALUES ($1, $2)"
	for _, product := range products {
		if _, err := tx.Exec(query, order.ID, product.ID); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	return order.ID, tx.Commit()
}

func (repo *OrderPostgresqlRepository) GetAll(userID int) ([]Order, error) {
	var orders []Order
	query := `SELECT o.id, o.creation_date, o.delivery_date FROM orders o
			  INNER JOIN orders_users ou on ou.order_id = o.id
			  INNER JOIN users u on ou.user_id = u.id
			  WHERE u.id = $1`

	if err := repo.DB.Select(&orders, query, userID); err != nil {
		return []Order{}, err
	}

	return orders, nil
}

func (repo *OrderPostgresqlRepository) GetByID(orderID int) (Order, error) {
	var order Order
	query := "SELECT * FROM orders WHERE id = $1"

	if err := repo.DB.Get(&order, query, orderID); err != nil {
		return Order{}, err
	}

	return order, nil
}
