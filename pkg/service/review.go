package service

import (
	"database/sql"
	"errors"
	"log"
	"market/pkg/model"
	"market/pkg/repository"
)

var (
	ErrReviewExists = errors.New("review exists")
)

type Review interface {
	Create(review model.Review) (int, error)
	Delete(reviewID int) (bool, error)
	Update(userID, productID int, text string) (bool, error)
	GetAll(productID int, orderBy string) ([]model.Review, error)
	GetReviewIDByProductIDUserID(productID, userID int) (int, error)
}

type ReviewService struct {
	reviewRepo repository.ReviewRepo
}

func NewReviewService(reviewRepo repository.ReviewRepo) *ReviewService {
	return &ReviewService{reviewRepo: reviewRepo}
}

func (s *ReviewService) Create(review model.Review) (int, error) {
	id, err := s.reviewRepo.GetReviewIDByProductIDUserID(review.ProductID, review.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("review doesnt't exist")
		} else {
			return 0, err
		}
	}
	if id != 0 {
		return 0, ErrReviewExists
	}

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

func (s *ReviewService) GetReviewIDByProductIDUserID(productID, userID int) (int, error) {
	return s.reviewRepo.GetReviewIDByProductIDUserID(productID, userID)
}
