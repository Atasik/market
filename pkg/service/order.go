package service

import (
	"market/pkg/model"
	"market/pkg/repository"
)

type Order interface {
	Create(userID int, order model.Order) (int, error)
	GetAll(userID int) ([]model.Order, error)
	GetByID(orderID int) (model.Order, error)
}

type OrderService struct {
	orderRepo  repository.OrderRepo
	basketRepo repository.BasketRepo
}

func NewOrderService(orderRepo repository.OrderRepo, basketRepo repository.BasketRepo) *OrderService {
	return &OrderService{orderRepo: orderRepo, basketRepo: basketRepo}
}

func (s *OrderService) Create(userID int, order model.Order) (int, error) {
	products, err := s.basketRepo.GetByID(userID)
	if err != nil {
		return 0, err
	}

	lastID, err := s.orderRepo.Create(userID, order, products)
	if err != nil {
		return 0, err
	}

	_, err = s.basketRepo.DeleteAll(userID)
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

func (s *OrderService) GetAll(userID int) ([]model.Order, error) {
	return s.orderRepo.GetAll(userID)
}

func (s *OrderService) GetByID(orderID int) (model.Order, error) {
	return s.orderRepo.GetByID(orderID)
}
