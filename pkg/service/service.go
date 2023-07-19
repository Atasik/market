package service

import (
	"market/pkg/repository"

	"github.com/cloudinary/cloudinary-go/v2"
)

type Service struct {
	Product
	Cart
	Order
	Review
	User
	Image
}

func NewService(repos *repository.Repository, cloudinary *cloudinary.Cloudinary) *Service {
	return &Service{
		Product: NewProductService(repos.ProductRepo, repos.UserRepo),
		Cart:    NewCartService(repos.CartRepo, repos.UserRepo, repos.ProductRepo),
		Order:   NewOrderService(repos.OrderRepo, repos.CartRepo, repos.UserRepo),
		Review:  NewReviewService(repos.ReviewRepo, repos.UserRepo, repos.ProductRepo),
		User:    NewUserService(repos.UserRepo),
		Image:   NewImageServiceCloudinary(cloudinary),
	}
}
