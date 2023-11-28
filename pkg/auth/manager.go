package auth

import (
	"context"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const TokenKey tokenKey = "token"

type tokenKey string

type tokenClaims struct {
	jwt.StandardClaims
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}

var ErrNoAuth = errors.New("no token found")

type TokenManager interface {
	NewJWT(userID int, username string, ttl time.Duration) (string, error)
	Parse(accessToken string) (*Token, error)
}

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}
	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) NewJWT(userID int, username string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userID,
		username,
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) Parse(accessToken string) (*Token, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return &Token{}, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return &Token{}, errors.New("token claims are not of type *tokenClaims")
	}

	return &Token{
		Username: claims.Username,
		UserID:   claims.UserID,
	}, nil
}

type Token struct {
	Username string
	UserID   int
}

func TokenFromContext(ctx context.Context) (*Token, error) {
	sess, ok := ctx.Value(TokenKey).(*Token)
	if !ok || sess == nil {
		return nil, ErrNoAuth
	}
	return sess, nil
}

func ContextWithToken(ctx context.Context, sess *Token) context.Context {
	return context.WithValue(ctx, TokenKey, sess)
}
