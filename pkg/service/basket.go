package service

import (
	"market/pkg/model"
	"market/pkg/repository"
)

type Basket interface {
	AddProduct(userId, productId int) (int, error)
	GetByID(userId int) ([]model.Product, error)
	DeleteProduct(userId, productId int) (bool, error)
	DeleteAll(userId int) (bool, error)
}

type BasketService struct {
	basketRepo repository.BasketRepo
}

func NewBasketService(basketRepo repository.BasketRepo) *BasketService {
	return &BasketService{basketRepo: basketRepo}
}

func (s *BasketService) AddProduct(userId, productId int) (int, error) {
	return s.basketRepo.AddProduct(userId, productId)
}

func (s *BasketService) GetByID(userId int) ([]model.Product, error) {
	return s.basketRepo.GetByID(userId)
}

func (s *BasketService) DeleteProduct(userId, productId int) (bool, error) {
	return s.basketRepo.DeleteProduct(userId, productId)
}

func (s *BasketService) DeleteAll(userId int) (bool, error) {
	return s.basketRepo.DeleteAll(userId)
}
