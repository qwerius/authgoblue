package session

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	sessionPrefix          = "github.com/qwerius/authgoblue:session:"
	sessionUserIndexPrefix = "github.com/qwerius/authgoblue:session:user:"
)

type RedisStore struct {
	client *redis.Client

	prefix string
}

func NewRedisStore(
	client *redis.Client,
) *RedisStore {

	return &RedisStore{
		client: client,
		prefix: sessionPrefix,
	}
}

func (s *RedisStore) context() context.Context {
	return context.Background()
}

func (s *RedisStore) key(
	id string,
) string {

	return s.prefix + id
}

func (s *RedisStore) userKey(
	userID string,
) string {

	return sessionUserIndexPrefix + userID
}

// Create menyimpan session.
// Jika session ID sudah ada, data lama akan diganti.
func (s *RedisStore) Create(
	session Session,
) error {

	ctx := s.context()

	data, err :=
		json.Marshal(
			session,
		)

	if err != nil {
		return err
	}

	ttl :=
		time.Until(
			session.ExpiresAt,
		)

	if ttl <= 0 {
		return nil
	}

	err =
		s.client.Set(
			ctx,
			s.key(session.ID),
			data,
			ttl,
		).Err()

	if err != nil {
		return err
	}

	if session.UserID != "" {

		err =
			s.client.SAdd(
				ctx,
				s.userKey(session.UserID),
				session.ID,
			).Err()

		if err != nil {

			_ = s.client.Del(
				ctx,
				s.key(session.ID),
			).Err()

			return err
		}

		// Bersihkan index user ketika semua session expired.
		_ = s.client.Expire(
			ctx,
			s.userKey(session.UserID),
			ttl,
		).Err()
	}

	return nil
}

// Get mengambil session berdasarkan ID.
func (s *RedisStore) Get(
	id string,
) (Session, error) {

	ctx := s.context()

	data, err :=
		s.client.Get(
			ctx,
			s.key(id),
		).Bytes()

	if err != nil {

		if err == redis.Nil {
			return Session{},
				ErrSessionNotFound
		}

		return Session{}, err
	}

	var session Session

	err =
		json.Unmarshal(
			data,
			&session,
		)

	if err != nil {
		return Session{}, err
	}

	return session, nil
}

// Revoke mencabut satu session.
func (s *RedisStore) Revoke(
	id string,
) error {

	session, err :=
		s.Get(
			id,
		)

	if err != nil {
		return err
	}

	if session.Revoked {
		return nil
	}

	session.Revoked = true

	return s.Create(
		session,
	)
}

// DeleteExpired tidak melakukan apa-apa.
//
// Redis menghapus session otomatis berdasarkan TTL.
func (s *RedisStore) DeleteExpired(
	now time.Time,
) error {

	return nil
}

// GetByUserID mengambil semua session milik user.
func (s *RedisStore) GetByUserID(
	userID string,
) ([]Session, error) {

	ctx := s.context()

	memberIDs, err :=
		s.client.SMembers(
			ctx,
			s.userKey(userID),
		).Result()

	if err != nil {
		return nil, err
	}

	result :=
		make(
			[]Session,
			0,
			len(memberIDs),
		)

	for _, memberID := range memberIDs {

		data, err :=
			s.client.Get(
				ctx,
				s.key(memberID),
			).Bytes()

		if err != nil {

			if err == redis.Nil {

				_ = s.client.SRem(
					ctx,
					s.userKey(userID),
					memberID,
				).Err()

				continue
			}

			return nil, err
		}

		var sess Session

		err =
			json.Unmarshal(
				data,
				&sess,
			)

		if err != nil {
			continue
		}

		result = append(
			result,
			sess,
		)
	}

	return result, nil
}

// RevokeAll mencabut seluruh session milik user.
func (s *RedisStore) RevokeAll(
	userID string,
) error {

	sessions, err :=
		s.GetByUserID(
			userID,
		)

	if err != nil {
		return err
	}

	for _, sess := range sessions {

		err =
			s.Revoke(
				sess.ID,
			)

		if err != nil {
			return err
		}
	}

	return nil
}
