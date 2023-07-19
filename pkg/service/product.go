package service

import (
	"errors"
	"market/pkg/model"
	"market/pkg/repository"
)

var (
	ErrPermissionDenied = errors.New("you have no access")
	ErrNoProduct        = errors.New("product doesn't exists")
	ErrProductExists    = errors.New("product already exists")
)

type Product interface {
	GetAll(orderBy string) ([]model.Product, error)
	GetByID(productID int) (model.Product, error)
	Create(product model.Product) (int, error)
	Update(userID, productID int, input model.UpdateProductInput) (bool, error)
	Delete(userID, productID int) (bool, error)
	GetByType(productType string, productID, limit int) ([]model.Product, error)
	IncreaseViewsCounter(productID int) (bool, error)
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

func (s *ProductService) GetByID(productID int) (model.Product, error) {
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return model.Product{}, ErrNoProduct
		default:
			return model.Product{}, err
		}
	}
	return product, nil
}

func (s *ProductService) Create(product model.Product) (int, error) {
	id, err := s.productRepo.Create(product)
	if err != nil {
		switch err {
		case repository.ErrAlreadyExists:
			return 0, ErrProductExists
		default:
			return 0, err
		}
	}
	return id, nil
}

func (s *ProductService) Update(userID, productID int, input model.UpdateProductInput) (bool, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return false, err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		if err := input.Validate(); err != nil {
			return false, err
		}
		return s.productRepo.Update(productID, input)
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

func (s *ProductService) GetByType(productType string, productID, limit int) ([]model.Product, error) {
	return s.productRepo.GetByType(productType, productID, limit)
}

func (s *ProductService) IncreaseViewsCounter(productID int) (bool, error) {
	views := 1
	input := model.UpdateProductInput{
		Views: &views,
	}
	if err := input.Validate(); err != nil {
		return false, err
	}
	return s.productRepo.Update(productID, input)
}
