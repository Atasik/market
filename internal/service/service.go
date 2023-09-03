package service

import (
	"context"
	"market/internal/model"
	"market/internal/repository"
	"market/pkg/auth"
	"market/pkg/hash"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
)

type User interface {
	CreateUser(model.User) (int, error)
	GenerateToken(username, password string) (string, error)
}

type Review interface {
	Create(review model.Review) (int, error)
	GetAll(productID int, q model.ReviewQueryInput) ([]model.Review, error)
	Update(userID, productID int, input model.UpdateReviewInput) error
	Delete(userID, reviewID int) error
}

type Product interface {
	Create(product model.Product) (int, error)
	GetAll(q model.ProductQueryInput) ([]model.Product, error)
	GetProductsByUserID(userID int, q model.ProductQueryInput) ([]model.Product, error)
	GetProductsByCategory(productCategory string, q model.ProductQueryInput) ([]model.Product, error)
	GetByID(productID int) (model.Product, error)
	Update(userID, productID int, input model.UpdateProductInput) error
	IncreaseViewsCounter(productID int) error
	Delete(userID, productID int) error
}

type Order interface {
	Create(userID int, order model.Order) (int, error)
	GetAll(userID int, q model.OrderQueryInput) ([]model.Order, error)
	GetByID(userID, orderID int) (model.Order, error)
	GetProductsByOrderID(userID, orderID int, q model.ProductQueryInput) ([]model.Product, error)
}

type Image interface {
	Upload(ctx context.Context, file multipart.File) (ImageData, error)
	Delete(ctx context.Context, imageID string) error
}

type Cart interface {
	Create(userID int) (int, error)
	AddProduct(userID, cartID, productID, amountToPurchase int) (int, error)
	GetByUserID(userID int) (model.Cart, error)
	GetAllProducts(userID, CartID int, q model.ProductQueryInput) ([]model.Product, error)
	UpdateProductAmount(userID, CartID, productID, amountToPurchase int) error
	DeleteProduct(userID, cartID, productID int) error
	DeleteAllProducts(userID, cartID int) error
}

type Service struct {
	Product
	Cart
	Order
	Review
	User
	Image
}

func NewService(repos *repository.Repository, cloudinary *cloudinary.Cloudinary, hasher hash.PasswordHasher, tokenManager auth.TokenManager) *Service {
	return &Service{
		Product: NewProductService(repos.ProductRepo, repos.UserRepo),
		Cart:    NewCartService(repos.CartRepo, repos.UserRepo, repos.ProductRepo),
		Order:   NewOrderService(repos.OrderRepo, repos.CartRepo, repos.UserRepo),
		Review:  NewReviewService(repos.ReviewRepo, repos.UserRepo, repos.ProductRepo),
		User:    NewUserService(repos.UserRepo, hasher, tokenManager),
		Image:   NewImageServiceCloudinary(cloudinary),
	}
}
