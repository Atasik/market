package order

import (
	"market/pkg/product"
	"time"
)

type Order struct {
	ID           int       `db:"id"`
	CreationDate time.Time `db:"creation_date"`
	DeliveryDate time.Time `db:"delivery_date"`
}

type OrderRepo interface {
	Create(userID int, order Order, products []product.Product) (int, error)
	GetAll(userID int) ([]Order, error)
	GetByID(orderID int) (Order, error)
}
