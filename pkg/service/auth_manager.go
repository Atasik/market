package service

import (
	"context"
	"errors"
)

var (
	ErrNoAuth = errors.New("no token found")
)

type Token struct {
	Username string
	UserID   int
}

type tokenKey string

const (
	TokenKey tokenKey = "token"
)

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
