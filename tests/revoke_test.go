package authgoblue_test

import (
	"testing"
	"time"

	"authgoblue/revoke"
)

func TestTokenRevocation(t *testing.T) {

	store :=
		revoke.NewMemoryStore()

	service :=
		revoke.NewService(
			store,
		)

	err :=
		service.Revoke(
			"token-123",
			time.Now().
				Add(time.Hour),
		)

	if err != nil {
		t.Fatal(err)
	}

	revoked, err :=
		service.IsRevoked(
			"token-123",
		)

	if err != nil {
		t.Fatal(err)
	}

	if !revoked {

		t.Fatal(
			"expected token revoked",
		)
	}

	revoked, err =
		service.IsRevoked(
			"token-456",
		)

	if err != nil {
		t.Fatal(err)
	}

	if revoked {

		t.Fatal(
			"unexpected revoked token",
		)
	}

}
