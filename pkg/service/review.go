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
	Delete(userID, reviewID int) (bool, error)
	Update(userID, productID int, input model.UpdateReviewInput) (bool, error)
	GetAll(productID int, orderBy string) ([]model.Review, error)
	GetReviewIDByProductIDUserID(productID, userID int) (int, error)
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
	_, err := s.productRepo.GetByID(review.ProductID)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return 0, ErrNoProduct
		default:
			return 0, err
		}
	}
	_, err = s.reviewRepo.GetReviewIDByProductIDUserID(review.ProductID, review.UserID)
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

func (s *ReviewService) Delete(userID, reviewID int) (bool, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return false, err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		return s.reviewRepo.Delete(reviewID)
	}

	return false, ErrPermissionDenied
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
