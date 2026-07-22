package authgoblue_test

import (
	"testing"

	"github.com/qwerius/authgoblue"
	"github.com/qwerius/authgoblue/revoke"

	"github.com/redis/go-redis/v9"
)

func TestAuthGoBlueUsesCustomRevokeStore(t *testing.T) {

	redisClient :=
		redis.NewClient(
			&redis.Options{
				Addr: "127.0.0.1:6379",
			},
		)

	store :=
		revoke.NewRedisStore(
			redisClient,
		)

	agb, err :=
		authgoblue.New(
			authgoblue.Config{

				Secret: "test-secret",

				Issuer: "test",

				Provider: &mockProvider{},

				RevokeStore: store,
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	if agb.Revoke == nil {

		t.Fatal(
			"expected revoke service",
		)

	}

}

func TestAuthGoBlueDefaultUsesMemoryRevokeStore(t *testing.T) {

	agb, err :=
		authgoblue.New(
			authgoblue.Config{

				Secret: "test-secret",

				Issuer: "test",

				Provider: &mockProvider{},
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	if agb.Revoke == nil {

		t.Fatal(
			"expected revoke service",
		)
	}
}
