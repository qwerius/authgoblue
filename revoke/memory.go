package revoke

import (
	"sync"
	"time"
)

type MemoryStore struct {
	mu sync.RWMutex

	data map[string]time.Time
}

func NewMemoryStore() *MemoryStore {

	return &MemoryStore{
		data: make(
			map[string]time.Time,
		),
	}
}

func (s *MemoryStore) Revoke(
	tokenID string,
	expireAt time.Time,
) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[tokenID] = expireAt

	return nil
}

func (s *MemoryStore) IsRevoked(
	tokenID string,
) (bool, error) {

	s.mu.RLock()

	expireAt, exists :=
		s.data[tokenID]

	s.mu.RUnlock()

	if !exists {
		return false, nil
	}

	if time.Now().After(expireAt) {

		s.mu.Lock()

		delete(
			s.data,
			tokenID,
		)

		s.mu.Unlock()

		return false, nil
	}

	return true, nil
}

func (s *MemoryStore) Consume(
	tokenID string,
	expireAt time.Time,
) (bool, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	if old, exists := s.data[tokenID]; exists {

		if time.Now().Before(old) {
			return false, nil
		}

		delete(
			s.data,
			tokenID,
		)
	}

	s.data[tokenID] = expireAt

	return true, nil
}
