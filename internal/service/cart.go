package service

import (
	"errors"
	"market/internal/model"
	"market/internal/repository"
	"market/pkg/database/postgres"
)

var (
	ErrAddDuplicate  = errors.New("product already added to cart")
	ErrInvalidAmount = errors.New("invalid amount")
)

type CartService struct {
	cartRepo    repository.CartRepo
	productRepo repository.ProductRepo
	userRepo    repository.UserRepo
}

func NewCartService(cartRepo repository.CartRepo, userRepo repository.UserRepo, productRepo repository.ProductRepo) *CartService {
	return &CartService{cartRepo: cartRepo, userRepo: userRepo, productRepo: productRepo}
}

func (s *CartService) Create(userID int) (int, error) {
	return s.cartRepo.Create(userID)
}

func (s *CartService) AddProduct(userID, cartID, productID, amountToPurchase int) (int, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return 0, err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		product, err := s.productRepo.GetByID(productID)
		if err != nil {
			return 0, err
		}
		_, err = s.cartRepo.GetProductByID(cartID, productID)
		if err != nil {
			switch err {
			case postgres.ErrNotFound:
				if product.Amount < amountToPurchase {
					return 0, ErrInvalidAmount
				}
				return s.cartRepo.AddProduct(cartID, productID, amountToPurchase)
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

func (s *CartService) GetAllProducts(userID, cartID int, q model.ProductQueryInput) ([]model.Product, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return []model.Product{}, err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		return s.cartRepo.GetAllProducts(cartID, q)
	}

	return []model.Product{}, ErrPermissionDenied
}

func (s *CartService) UpdateProductAmount(userID, cartID, productID, amountToPurchase int) error {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		product, err := s.cartRepo.GetProductByID(cartID, productID)
		if err != nil {
			return err
		}
		if product.Amount < amountToPurchase {
			return ErrInvalidAmount
		}
		return s.cartRepo.UpdateProductAmount(cartID, productID, amountToPurchase)
	}

	return ErrPermissionDenied
}

func (s *CartService) DeleteProduct(userID, cartID, productID int) error {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		return s.cartRepo.DeleteProduct(cartID, productID)
	}

	return ErrPermissionDenied
}

func (s *CartService) DeleteAllProducts(userID, cartID int) error {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		return s.cartRepo.DeleteAllProducts(cartID)
	}

	return ErrPermissionDenied
}
