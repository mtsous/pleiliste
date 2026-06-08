package session

import (
	"log/slog"
	"sync"
	"time"

	"github.com/mtsous/pleiliste/internal/core"
)

type Store struct {
	mu   sync.RWMutex
	data map[string]*core.Session
	ttl  time.Duration
}

func NewStore(ttl time.Duration) *Store {
	s := &Store{
		data: make(map[string]*core.Session),
		ttl:  ttl,
	}

	go s.cleanup()

	return s
}

func (s *Store) cleanup() {
	for range time.Tick(time.Minute) {
		s.mu.Lock()

		for id, sess := range s.data {
			if time.Now().After(sess.ExpiresAt.Add(s.ttl)) {
				delete(s.data, id)
			}
		}

		s.mu.Unlock()
	}
}

func (s *Store) Set(key string, sess *core.Session) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = sess
}

func (s *Store) Get(key string) (*core.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sess, ok := s.data[key]
	if !ok {
		slog.Error(
			"key not found in session store",
			"key", key,
		)
		return nil, core.ErrSessionStoreKeyNotFound
	}

	return sess, nil
}
