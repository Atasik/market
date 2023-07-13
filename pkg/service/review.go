package service

import (
	"market/pkg/model"
	"market/pkg/repository"
)

type Review interface {
	Create(review model.Review) (int, error)
	Delete(reviewID int) (bool, error)
	Update(userID, productID int, text string) (bool, error)
	GetAll(productID int, orderBy string) ([]model.Review, error)
}

type ReviewService struct {
	reviewRepo repository.ReviewRepo
}

func NewReviewService(reviewRepo repository.ReviewRepo) *ReviewService {
	return &ReviewService{reviewRepo: reviewRepo}
}

func (s *ReviewService) Create(review model.Review) (int, error) {
	return s.reviewRepo.Create(review)
}

func (s *ReviewService) Delete(reviewID int) (bool, error) {
	return s.reviewRepo.Delete(reviewID)
}

func (s *ReviewService) Update(userID, productID int, text string) (bool, error) {
	return s.reviewRepo.Update(userID, productID, text)
}

func (s *ReviewService) GetAll(productID int, orderBy string) ([]model.Review, error) {
	return s.reviewRepo.GetAll(productID, orderBy)
}
