package service

import (
	"errors"
	"market/internal/model"
	"market/internal/repository"
	"market/pkg/auth"
	"market/pkg/hash"
	"time"
)

var (
	ErrBadPass    = errors.New("wrong password")
	ErrNoUser     = errors.New("no user found")
	ErrUserExists = errors.New("user already exists")
)

type UserService struct {
	userRepo     repository.UserRepo
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager

	accessTokenTTL time.Duration
}

func NewUserService(userRepo repository.UserRepo, hasher hash.PasswordHasher, tokenManager auth.TokenManager, accessTTL time.Duration) *UserService {
	return &UserService{userRepo: userRepo, hasher: hasher, tokenManager: tokenManager, accessTokenTTL: accessTTL}
}

func (s *UserService) CreateUser(user model.User) (int, error) {
	password, err := s.hasher.Hash(user.Password)
	if err != nil {
		return 0, err
	}

	user.Password = password
	return s.userRepo.CreateUser(user)
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
	return s.tokenManager.NewJWT(user.ID, username, s.accessTokenTTL)
}
