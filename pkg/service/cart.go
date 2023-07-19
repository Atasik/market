package service

import (
	"database/sql"
	"errors"
	"market/pkg/model"
	"market/pkg/repository"
)

var (
	ErrAddDuplicate = errors.New("product already added to cart")
)

type Cart interface {
	CreateCart(userID int) (int, error)
	AddProduct(userID, CartID, productID int) (int, error)
	GetByUserID(userID int) (model.Cart, error)
	GetProducts(userID, CartID int) ([]model.Product, error)
	DeleteProduct(userID, CartID, productID int) (bool, error)
	DeleteAll(userID, CartID int) (bool, error)
}

type CartService struct {
	CartRepo repository.CartRepo
	userRepo repository.UserRepo
}

func NewCartService(CartRepo repository.CartRepo, userRepo repository.UserRepo) *CartService {
	return &CartService{CartRepo: CartRepo, userRepo: userRepo}
}

func (s *CartService) CreateCart(userID int) (int, error) {
	return s.CartRepo.CreateCart(userID)
}

func (s *CartService) AddProduct(userID, CartID, productID int) (int, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return 0, err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		_, err := s.CartRepo.GetProductByID(CartID, productID)
		if err != nil {
			if err == sql.ErrNoRows {
				return s.CartRepo.AddProduct(CartID, productID)
			}
			return 0, err
		}
		return 0, ErrAddDuplicate
	}

	return 0, ErrPermissionDenied
}

func (s *CartService) GetByUserID(userID int) (model.Cart, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return model.Cart{}, err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		return s.CartRepo.GetByUserID(userID)
	}
	return model.Cart{}, ErrPermissionDenied
}

func (s *CartService) GetProducts(userID, CartID int) ([]model.Product, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return []model.Product{}, err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		return s.CartRepo.GetProducts(CartID)
	}

	return []model.Product{}, ErrPermissionDenied
}

func (s *CartService) DeleteProduct(userID, CartID, productID int) (bool, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return false, err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		return s.CartRepo.DeleteProduct(CartID, productID)
	}

	return false, ErrPermissionDenied
}

func (s *CartService) DeleteAll(userID, CartID int) (bool, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return false, err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		return s.CartRepo.DeleteAll(CartID)
	}

	return false, ErrPermissionDenied
}
