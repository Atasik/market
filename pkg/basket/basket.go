package basket

import "market/pkg/product"

type BasketRepo interface {
	AddProduct(userId, productId int) (int, error)
	GetByID(buserId int) ([]product.Product, error)
	DeleteProduct(userId, productId int) (bool, error)
	DeleteAll(userId int) (bool, error)
}
