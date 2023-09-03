package repository

import (
	"market/internal/model"

	"github.com/jmoiron/sqlx"
)

const (
	usersTable          = "users"
	productsTable       = "products"
	ordersTable         = "orders"
	reviewsTable        = "reviews"
	cartsTable          = "carts"
	productsUsersTable  = "products_users"
	productsOrdersTable = "products_orders"
	productsCartsTable  = "products_carts"
)

type ProductRepo interface {
	Create(product model.Product) (int, error)
	GetAll(q model.ProductQueryInput) ([]model.Product, error)
	GetProductsByUserID(userID int, q model.ProductQueryInput) ([]model.Product, error)
	GetByID(productID int) (model.Product, error)
	GetProductsByCategory(productCategory string, q model.ProductQueryInput) ([]model.Product, error)
	Update(productID int, input model.UpdateProductInput) error
	Delete(productID int) error
}

type OrderRepo interface {
	Create(cartID, userID int, order model.Order) (int, error)
	GetAll(userID int, q model.OrderQueryInput) ([]model.Order, error)
	GetByID(orderID int) (model.Order, error)
	GetProductsByOrderID(orderID int, q model.ProductQueryInput) ([]model.Product, error)
}

type ReviewRepo interface {
	Create(review model.Review) (int, error)
	Delete(reviewID int) error
	Update(userID, productID int, input model.UpdateReviewInput) error
	GetAll(productID int, q model.ReviewQueryInput) ([]model.Review, error)
	GetReviewIDByProductIDUserID(productID, userID int) (int, error)
}

type CartRepo interface {
	Create(userID int) (int, error)
	AddProduct(cartID, productID, amount int) (int, error)
	GetByUserID(userID int) (model.Cart, error)
	GetProductByID(cartID, productID int) (model.Product, error)
	GetAllProducts(cartID int, q model.ProductQueryInput) ([]model.Product, error)
	UpdateProductAmount(cartID, productID, amount int) error
	DeleteProduct(cartID, productID int) error
	DeleteAllProducts(cartID int) error
}

type UserRepo interface {
	GetUser(login string) (model.User, error)
	GetUserById(userID int) (model.User, error)
	CreateUser(model.User) (int, error)
}

type Repository struct {
	CartRepo
	OrderRepo
	ProductRepo
	UserRepo
	ReviewRepo
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		CartRepo:    NewCartPostgresqlRepo(db),
		OrderRepo:   NewOrderPostgresqlRepo(db),
		ProductRepo: NewProductPostgresqlRepo(db),
		UserRepo:    NewUserPostgresqlRepo(db),
		ReviewRepo:  NewReviewPostgresqlRepo(db),
	}
}
