package service

import (
	"errors"
	"market/pkg/model"
	"market/pkg/repository"
)

var (
	ErrAddDuplicate  = errors.New("product already added to cart")
	ErrInvalidAmount = errors.New("invalid amount")
)

type Cart interface {
	CreateCart(userID int) (int, error)
	AddProduct(userID, CartID, productID, amountToPurcahse int) (int, error)
	GetByUserID(userID int) (model.Cart, error)
	GetProducts(userID, CartID int) ([]model.Product, error)
	DeleteProduct(userID, CartID, productID int) (bool, error)
	DeleteAll(userID, CartID int) (bool, error)
}

type CartService struct {
	cartRepo    repository.CartRepo
	productRepo repository.ProductRepo
	userRepo    repository.UserRepo
}

func NewCartService(cartRepo repository.CartRepo, userRepo repository.UserRepo, productRepo repository.ProductRepo) *CartService {
	return &CartService{cartRepo: cartRepo, userRepo: userRepo, productRepo: productRepo}
}

func (s *CartService) CreateCart(userID int) (int, error) {
	return s.cartRepo.CreateCart(userID)
}

func (s *CartService) AddProduct(userID, CartID, productID, amountToPurcahse int) (int, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return 0, err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		product, err := s.productRepo.GetByID(productID)
		if err != nil {
			switch err {
			case repository.ErrNotFound:
				return 0, ErrNoProduct
			default:
				return 0, err
			}
		}
		_, err = s.cartRepo.GetProductByID(CartID, productID)
		if err != nil {
			switch err {
			case repository.ErrNotFound:
				if product.Amount < amountToPurcahse {
					print("kek", product.Amount)
					return 0, ErrInvalidAmount
				}
				return s.cartRepo.AddProduct(CartID, productID, amountToPurcahse)
			default:
				return 0, err
			}
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
		return s.cartRepo.GetByUserID(userID)
	}
	return model.Cart{}, ErrPermissionDenied
}

func (s *CartService) GetProducts(userID, CartID int) ([]model.Product, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return []model.Product{}, err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		return s.cartRepo.GetProducts(CartID)
	}

	return []model.Product{}, ErrPermissionDenied
}

func (s *CartService) DeleteProduct(userID, CartID, productID int) (bool, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return false, err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		return s.cartRepo.DeleteProduct(CartID, productID)
	}

	return false, ErrPermissionDenied
}

func (s *CartService) DeleteAll(userID, CartID int) (bool, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return false, err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		return s.cartRepo.DeleteAll(CartID)
	}

	return false, ErrPermissionDenied
}
