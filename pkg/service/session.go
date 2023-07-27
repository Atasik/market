package service

import (
	"context"
	"errors"
)

var (
	ErrNoAuth = errors.New("no session found")
)

type Session struct {
	Username string
	UserID   int
}

type sessKey string

const (
	SessionKey sessKey = "token"
)

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(SessionKey).(*Session)
	if !ok || sess == nil {
		return nil, ErrNoAuth
	}
	return sess, nil
}

func ContextWithSession(ctx context.Context, sess *Session) context.Context {
	return context.WithValue(ctx, SessionKey, sess)
}
