package service

import (
	"market/pkg/repository"

	"github.com/cloudinary/cloudinary-go/v2"
)

type Service struct {
	Product
	Basket
	Order
	Review
	User
	Image
}

func NewService(repos *repository.Repository, cloudinary *cloudinary.Cloudinary) *Service {
	return &Service{
		Product: NewProductService(repos.ProductRepo),
		Basket:  NewBasketService(repos.BasketRepo),
		Order:   NewOrderService(repos.OrderRepo, repos.BasketRepo),
		Review:  NewReviewService(repos.ReviewRepo),
		User:    NewUserService(repos.UserRepo),
		Image:   NewImageServiceCloudinary(cloudinary),
	}
}
