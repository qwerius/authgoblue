package authgoblue_test

import (
	"context"
	"testing"
	"time"

	"authgoblue/revoke"

	"github.com/redis/go-redis/v9"
)

func TestRedisRevokeStore(
	t *testing.T,
) {

	client :=
		redis.NewClient(
			&redis.Options{
				Addr: "localhost:6379",
			},
		)

	ctx := context.Background()

	err :=
		client.Ping(
			ctx,
		).Err()

	if err != nil {

		t.Skip(
			"redis not running",
		)

	}

	store :=
		revoke.NewRedisStore(
			client,
		)

	err =
		store.Revoke(
			"test-jti-123",
			time.Now().Add(
				time.Minute,
			),
		)

	if err != nil {

		t.Fatal(err)

	}

	revoked, err :=
		store.IsRevoked(
			"test-jti-123",
		)

	if err != nil {

		t.Fatal(err)

	}

	if !revoked {

		t.Fatal(
			"expected revoked token",
		)

	}

}
