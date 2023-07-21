package service

import (
	"errors"
	"market/pkg/model"
	"market/pkg/repository"
)

var (
	ErrReviewExists = errors.New("review exists")
	ErrNoReview     = errors.New("review doesn't exists")
)

type Review interface {
	Create(review model.Review) (int, error)
	Delete(userID, reviewID int) error
	Update(userID, productID int, input model.UpdateReviewInput) error
	GetAll(productID int) ([]model.Review, error)
}

type ReviewService struct {
	reviewRepo  repository.ReviewRepo
	userRepo    repository.UserRepo
	productRepo repository.ProductRepo
}

func NewReviewService(reviewRepo repository.ReviewRepo, userRepo repository.UserRepo, productRepo repository.ProductRepo) *ReviewService {
	return &ReviewService{reviewRepo: reviewRepo, userRepo: userRepo, productRepo: productRepo}
}

func (s *ReviewService) Create(review model.Review) (int, error) {
	_, err := s.reviewRepo.GetReviewIDByProductIDUserID(review.ProductID, review.UserID)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return s.reviewRepo.Create(review)
		default:
			return 0, err
		}
	}
	return 0, ErrReviewExists
}

func (s *ReviewService) Delete(userID, reviewID int) error {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		return s.reviewRepo.Delete(reviewID)
	}

	return ErrPermissionDenied
}

func (s *ReviewService) Update(userID, productID int, input model.UpdateReviewInput) error {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		if err := input.Validate(); err != nil {
			return err
		}
		return s.reviewRepo.Update(userID, productID, input)
	}

	return ErrPermissionDenied
}

func (s *ReviewService) GetAll(productID int) ([]model.Review, error) {
	return s.reviewRepo.GetAll(productID)
}
