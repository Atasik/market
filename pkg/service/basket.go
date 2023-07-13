package service

import (
	"market/pkg/model"
	"market/pkg/repository"
)

type Basket interface {
	CreateBasket(userID int) (int, error)
	AddProduct(basketID, productID int) (int, error)
	GetByUserID(userID int) (model.Basket, error)
	GetProducts(basketID int) ([]model.Product, error)
	DeleteProduct(basketID, productID int) (bool, error)
	DeleteAll(basketID int) (bool, error)
}

type BasketService struct {
	basketRepo repository.BasketRepo
}

func NewBasketService(basketRepo repository.BasketRepo) *BasketService {
	return &BasketService{basketRepo: basketRepo}
}

func (s *BasketService) CreateBasket(userID int) (int, error) {
	return s.basketRepo.CreateBasket(userID)
}

func (s *BasketService) AddProduct(basketID, productID int) (int, error) {
	return s.basketRepo.AddProduct(basketID, productID)
}

func (s *BasketService) GetByUserID(userID int) (model.Basket, error) {
	return s.basketRepo.GetByUserID(userID)
}

func (s *BasketService) GetProducts(basketID int) ([]model.Product, error) {
	return s.basketRepo.GetProducts(basketID)
}

func (s *BasketService) DeleteProduct(basketID, productID int) (bool, error) {
	return s.basketRepo.DeleteProduct(basketID, productID)
}

func (s *BasketService) DeleteAll(basketID int) (bool, error) {
	return s.basketRepo.DeleteAll(basketID)
}
