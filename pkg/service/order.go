package service

import (
	"errors"
	"market/pkg/model"
	"market/pkg/repository"
)

var (
	ErrNoOrder    = errors.New("order doesn't exists")
	ErrNoProducts = errors.New("no products in cart")
)

type Order interface {
	Create(userID int, order model.Order) (int, error)
	GetAll(userID int, q model.OrderQueryInput) ([]model.Order, error)
	GetByID(userID, orderID int) (model.Order, error)
	GetProductsByOrderID(userID, orderID int, q model.ProductQueryInput) ([]model.Product, error)
}

type OrderService struct {
	orderRepo repository.OrderRepo
	cartRepo  repository.CartRepo
	userRepo  repository.UserRepo
}

func NewOrderService(orderRepo repository.OrderRepo, cartRepo repository.CartRepo, userRepo repository.UserRepo) *OrderService {
	return &OrderService{orderRepo: orderRepo, cartRepo: cartRepo, userRepo: userRepo}
}

func (s *OrderService) Create(userID int, order model.Order) (int, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return 0, err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		cart, err := s.cartRepo.GetByUserID(userID)
		if err != nil {
			return 0, err
		}

		q := model.ProductQueryInput{
			QueryInput: model.QueryInput{
				Limit:     0,
				Offset:    0,
				SortBy:    model.SortByDate,
				SortOrder: model.DESCENDING,
			},
		}

		order.Products, err = s.cartRepo.GetAllProducts(cart.ID, q)
		if err != nil {
			return 0, err
		}

		if order.Products == nil {
			return 0, ErrNoProducts
		}

		lastID, err := s.orderRepo.Create(cart.ID, userID, order)
		if err != nil {
			return 0, err
		}
		return lastID, nil
	}

	return 0, ErrPermissionDenied
}

func (s *OrderService) GetAll(userID int, q model.OrderQueryInput) ([]model.Order, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return []model.Order{}, err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		return s.orderRepo.GetAll(userID, q)
	}

	return []model.Order{}, ErrPermissionDenied
}

func (s *OrderService) GetByID(userID, orderID int) (model.Order, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return model.Order{}, err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		order, err := s.orderRepo.GetByID(orderID)
		if err != nil {
			switch err {
			case repository.ErrNotFound:
				return model.Order{}, ErrNoOrder
			default:
				return model.Order{}, err
			}
		}
		return order, nil
	}

	return model.Order{}, ErrPermissionDenied
}

func (s *OrderService) GetProductsByOrderID(userID, orderID int, q model.ProductQueryInput) ([]model.Product, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return []model.Product{}, err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		return s.orderRepo.GetProductsByOrderID(orderID, q)
	}

	return []model.Product{}, ErrPermissionDenied
}
