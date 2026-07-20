package authgoblue_test

import (
	"testing"

	"authgoblue"
	"authgoblue/revoke"

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

	agb :=
		authgoblue.New(
			authgoblue.Config{

				Secret: "test-secret",

				Issuer: "test",

				RevokeStore: store,
			},
		)

	if agb.Revoke == nil {

		t.Fatal(
			"expected revoke service",
		)

	}

}

func TestAuthGoBlueDefaultUsesMemoryRevokeStore(t *testing.T) {

	agb :=
		authgoblue.New(
			authgoblue.Config{

				Secret: "test-secret",
				Issuer: "test",
			},
		)
	if agb.Revoke == nil {

		t.Fatal(
			"expected revoke service",
		)
	}
}
