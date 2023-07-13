package session

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"market/pkg/model"
)

type Session struct {
	ID            string
	UserID        int
	UserName      string
	UserType      model.Role
	BasketStorage map[int]struct{}
}

func (s *Session) IsPurchased(id int) bool {
	_, ok := s.BasketStorage[id]
	return ok
}

func (s *Session) AddPurchase(id int) {
	s.BasketStorage[id] = struct{}{}
}

func (s *Session) DeletePurchase(id int) {
	delete(s.BasketStorage, id)
}

func (s *Session) PurgeBasket() {
	s.BasketStorage = make(map[int]struct{})
}

func NewSession(userID int, userName string, userType model.Role) *Session {
	randID := make([]byte, 16)
	rand.Read(randID)

	return &Session{
		ID:            fmt.Sprintf("%x", randID),
		UserID:        userID,
		UserName:      userName,
		UserType:      userType,
		BasketStorage: map[int]struct{}{},
	}
}

var (
	ErrNoAuth = errors.New("no session found")
)

type sessKey string

var SessionKey sessKey = "sessionKey"

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(SessionKey).(*Session)
	if !ok || sess == nil {
		return &Session{}, ErrNoAuth
	}
	return sess, nil
}

func ContextWithSession(ctx context.Context, sess *Session) context.Context {
	return context.WithValue(ctx, SessionKey, sess)
}
