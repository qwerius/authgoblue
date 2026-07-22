package revoke

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
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
		prefix: "revoke:",
	}
}

func (s *RedisStore) key(
	tokenID string,
) string {

	return s.prefix + tokenID
}

func (s *RedisStore) Revoke(
	tokenID string,
	expireAt time.Time,
) error {

	ttl :=
		time.Until(
			expireAt,
		)

	if ttl <= 0 {
		return nil
	}

	return s.client.Set(
		context.Background(),

		s.key(tokenID),

		"1",

		ttl,
	).Err()
}

func (s *RedisStore) IsRevoked(
	tokenID string,
) (bool, error) {

	ctx := context.Background()

	count, err :=
		s.client.Exists(
			ctx,
			s.key(tokenID),
		).Result()

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *RedisStore) Consume(
	tokenID string,
	expireAt time.Time,
) (bool, error) {

	ttl :=
		time.Until(
			expireAt,
		)

	if ttl <= 0 {
		return false, nil
	}

	ok, err :=
		s.client.SetNX(
			context.Background(),
			s.key(tokenID),
			"1",
			ttl,
		).Result()

	if err != nil {
		return false, err
	}

	return ok, nil
}
