package authgoblue_test

import (
	"context"
	"testing"
	"time"

	"github.com/qwerius/authgoblue/session"

	"github.com/redis/go-redis/v9"
)

func TestRedisSessionStore(
	t *testing.T,
) {

	client :=
		redis.NewClient(
			&redis.Options{
				Addr: "localhost:6379",
			},
		)

	err :=
		client.Ping(
			context.Background(),
		).Err()

	if err != nil {

		t.Skip(
			"redis not running",
		)
	}

	store :=
		session.NewRedisStore(
			client,
		)

	s :=
		session.Session{

			ID: "session-test-1",

			UserID: "user-1",

			ExpiresAt: time.Now().
				Add(
					time.Hour,
				),
		}

	err =
		store.Create(
			s,
		)

	if err != nil {
		t.Fatal(err)
	}

	result, err :=
		store.Get(
			"session-test-1",
		)

	if err != nil {
		t.Fatal(err)
	}

	if result.UserID != "user-1" {

		t.Fatal(
			"invalid session",
		)
	}
}
