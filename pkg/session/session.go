package session

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
)

type Session struct {
	ID        string
	UserID    int
	UserName  string
	UserType  string
	Purchases map[int]struct{}
}

func (s *Session) IsPurchased(id int) bool {
	_, ok := s.Purchases[id]
	return ok
}

// мб race-condition
func (s *Session) AddPurchase(id int) {
	s.Purchases[id] = struct{}{}
}

func (s *Session) DeletePurchase(id int) {
	delete(s.Purchases, id)
}

func NewSession(userID int, userName, userType string) *Session {
	// лучше генерировать из заданного алфавита, но так писать меньше и для учебного примера ОК
	randID := make([]byte, 16)
	rand.Read(randID)

	return &Session{
		ID:        fmt.Sprintf("%x", randID),
		UserID:    userID,
		UserName:  userName,
		UserType:  userType,
		Purchases: map[int]struct{}{},
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
