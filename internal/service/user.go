package service

import (
	"errors"
	"market/internal/model"
	"market/internal/repository"
	auth "market/pkg/auth"
	"market/pkg/hash"
	"time"
)

var (
	ErrBadPass    = errors.New("wrong password")
	ErrNoUser     = errors.New("no user found")
	ErrUserExists = errors.New("user already exists")
)

const (
	tokenTTL = 12 * time.Hour
)

type UserService struct {
	userRepo     repository.UserRepo
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
}

func NewUserService(userRepo repository.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(user model.User) (int, error) {
	password, err := s.hasher.Hash(user.Password)
	if err != nil {
		return 0, err
	}

	user.Password = password

	id, err := s.userRepo.CreateUser(user)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *UserService) GenerateToken(username, password string) (string, error) {
	user, err := s.userRepo.GetUser(username)
	if err != nil {
		return "", err
	}

	match, err := s.hasher.Verify(password, user.Password)
	if err != nil {
		return "", err
	}

	if !match {
		return "", ErrBadPass
	}

	return s.tokenManager.NewJWT(user.ID, username, tokenTTL)
}
