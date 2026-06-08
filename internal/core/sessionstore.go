package core

import (
	"errors"
	"time"
)

type Session struct {
	State        string
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

type SessionStore interface {
	Set(key string, sess *Session)
	Get(key string) (*Session, error)
}

var (
	ErrSessionStoreKeyNotFound = errors.New("key provided was not found in session store")
)
