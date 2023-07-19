package service

import (
	"market/pkg/model"
	"market/pkg/repository"
)

type Order interface {
	Create(userID int, order model.Order) (int, error)
	GetAll(userID int) ([]model.Order, error)
	GetByID(userID, orderID int) (model.Order, error)
}

type OrderService struct {
	orderRepo repository.OrderRepo
	CartRepo  repository.CartRepo
	userRepo  repository.UserRepo
}

func NewOrderService(orderRepo repository.OrderRepo, CartRepo repository.CartRepo, userRepo repository.UserRepo) *OrderService {
	return &OrderService{orderRepo: orderRepo, CartRepo: CartRepo, userRepo: userRepo}
}

func (s *OrderService) Create(userID int, order model.Order) (int, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return 0, err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		Cart, err := s.CartRepo.GetByUserID(userID)
		if err != nil {
			return 0, err
		}

		products, err := s.CartRepo.GetProducts(Cart.ID)
		if err != nil {
			return 0, err
		}

		lastID, err := s.orderRepo.Create(userID, order, products)
		if err != nil {
			return 0, err
		}

		_, err = s.CartRepo.DeleteAll(userID)
		if err != nil {
			return 0, err
		}
		return lastID, nil
	}

	return 0, ErrPermissionDenied
}

func (s *OrderService) GetAll(userID int) ([]model.Order, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return []model.Order{}, err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		return s.orderRepo.GetAll(userID)
	}

	return []model.Order{}, ErrPermissionDenied
}

func (s *OrderService) GetByID(userID, orderID int) (model.Order, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return model.Order{}, err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		return s.orderRepo.GetByID(orderID)
	}

	return model.Order{}, ErrPermissionDenied
}
