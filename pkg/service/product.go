package service

import (
	"errors"
	"market/pkg/model"
	"market/pkg/repository"
)

var (
	ErrPermissionDenied = errors.New("you have no access")
)

type Product interface {
	GetAll(orderBy string) ([]model.Product, error)
	GetByID(productId int) (model.Product, error)
	Create(product model.Product) (int, error)
	Update(userID, productId int, input model.UpdateProductInput) (bool, error)
	Delete(userID, productId int) (bool, error)
	GetByType(productType string, limit int) ([]model.Product, error)
}

type ProductService struct {
	productRepo repository.ProductRepo
	userRepo    repository.UserRepo
}

func NewProductService(productRepo repository.ProductRepo, userRepo repository.UserRepo) *ProductService {
	return &ProductService{productRepo: productRepo, userRepo: userRepo}
}

func (s *ProductService) GetAll(orderBy string) ([]model.Product, error) {
	return s.productRepo.GetAll(orderBy)
}

func (s *ProductService) GetByID(productId int) (model.Product, error) {
	return s.productRepo.GetByID(productId)
}

func (s *ProductService) Create(product model.Product) (int, error) {
	return s.productRepo.Create(product)
}

func (s *ProductService) Update(userID, productId int, input model.UpdateProductInput) (bool, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return false, err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		if err := input.Validate(); err != nil {
			return false, err
		}
		return s.productRepo.Update(productId, input)
	}

	return false, ErrPermissionDenied
}

func (s *ProductService) Delete(userID, productID int) (bool, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return false, err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		return s.productRepo.Delete(productID)
	}

	return false, ErrPermissionDenied
}

func (s *ProductService) GetByType(productType string, limit int) ([]model.Product, error) {
	return s.productRepo.GetByType(productType, limit)
}
