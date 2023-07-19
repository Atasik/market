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
	Delete(userID, reviewID int) (bool, error)
	Update(userID, productID int, input model.UpdateReviewInput) (bool, error)
	GetAll(productID int, orderBy string) ([]model.Review, error)
	GetReviewIDByProductIDUserID(productID, userID int) (int, error)
}

type ReviewService struct {
	reviewRepo repository.ReviewRepo
	userRepo   repository.UserRepo
}

func NewReviewService(reviewRepo repository.ReviewRepo, userRepo repository.UserRepo) *ReviewService {
	return &ReviewService{reviewRepo: reviewRepo, userRepo: userRepo}
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

func (s *ReviewService) Delete(userID, reviewID int) (bool, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return false, err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		return s.reviewRepo.Delete(reviewID)
	}

	return s.reviewRepo.Delete(reviewID)
}

func (s *ReviewService) Update(userID, productID int, input model.UpdateReviewInput) (bool, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return false, err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		if err := input.Validate(); err != nil {
			return false, err
		}
		return s.reviewRepo.Update(userID, productID, input)
	}

	return false, ErrPermissionDenied
}

func (s *ReviewService) GetAll(productID int, orderBy string) ([]model.Review, error) {
	return s.reviewRepo.GetAll(productID, orderBy)
}

func (s *ReviewService) GetReviewIDByProductIDUserID(productID, userID int) (int, error) {
	return s.reviewRepo.GetReviewIDByProductIDUserID(productID, userID)
}
