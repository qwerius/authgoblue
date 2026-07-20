package session

import (
	"sync"
	"time"
)

type MemoryStore struct {
	mu sync.RWMutex

	data map[string]Session
}

func NewMemoryStore() *MemoryStore {

	return &MemoryStore{
		data: make(
			map[string]Session,
		),
	}
}

// Create menyimpan session.
// Jika ID sudah ada, data lama akan diganti.
func (s *MemoryStore) Create(
	session Session,
) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[session.ID] = session

	return nil
}

// Get mengambil session berdasarkan ID.
func (s *MemoryStore) Get(
	id string,
) (Session, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok :=
		s.data[id]

	if !ok {
		return Session{},
			ErrSessionNotFound
	}

	return session, nil
}

// GetByUserID mengambil seluruh session milik user.
func (s *MemoryStore) GetByUserID(
	userID string,
) ([]Session, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	result :=
		make(
			[]Session,
			0,
		)

	now := time.Now()

	for _, sess := range s.data {

		if sess.UserID != userID {
			continue
		}

		// Jangan tampilkan session expired.
		if !sess.ExpiresAt.After(now) {
			continue
		}

		result = append(
			result,
			sess,
		)
	}

	return result, nil
}

// Revoke mencabut satu session.
func (s *MemoryStore) Revoke(
	id string,
) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok :=
		s.data[id]

	if !ok {
		return ErrSessionNotFound
	}

	if session.Revoked {
		return nil
	}

	session.Revoked = true

	s.data[id] = session

	return nil
}

// RevokeAll mencabut seluruh session milik user.
func (s *MemoryStore) RevokeAll(
	userID string,
) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	for id, sess := range s.data {

		if sess.UserID != userID {
			continue
		}

		if sess.Revoked {
			continue
		}

		sess.Revoked = true

		s.data[id] = sess
	}

	return nil
}

// DeleteExpired menghapus session yang sudah melewati masa aktif.
func (s *MemoryStore) DeleteExpired(
	now time.Time,
) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	for id, sess := range s.data {

		if !sess.ExpiresAt.After(now) {

			delete(
				s.data,
				id,
			)
		}
	}

	return nil
}
