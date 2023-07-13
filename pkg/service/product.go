package service

import (
	"market/pkg/model"
	"market/pkg/repository"
)

type Product interface {
	GetAll(orderBy string) ([]model.Product, error)
	GetByID(productId int) (model.Product, error)
	Create(product model.Product) (int, error)
	Update(productId int, input model.UpdateProductInput) (bool, error)
	Delete(productId int) (bool, error)
	GetByType(productType string, limit int) ([]model.Product, error)
}

type ProductService struct {
	productRepo repository.ProductRepo
}

func NewProductService(productRepo repository.ProductRepo) *ProductService {
	return &ProductService{productRepo: productRepo}
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

func (s *ProductService) Update(productId int, input model.UpdateProductInput) (bool, error) {
	if err := input.Validate(); err != nil {
		return false, err
	}

	return s.productRepo.Update(productId, input)
}

func (s *ProductService) Delete(productId int) (bool, error) {
	return s.productRepo.Delete(productId)
}

func (s *ProductService) GetByType(productType string, limit int) ([]model.Product, error) {
	return s.productRepo.GetByType(productType, limit)
}
