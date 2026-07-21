package authgoblue_test

import (
	"testing"
	"time"

	"github.com/qwerius/authgoblue/hooks"
	"github.com/qwerius/authgoblue/session"
)

func TestSessionCreate(t *testing.T) {

	store := session.NewMemoryStore()

	service := session.NewService(
		store,
		hooks.NewRegistry(),
	)

	s, err := service.Create(
		"user-1",
	)

	if err != nil {
		t.Fatal(err)
	}

	if s.ID == "" {
		t.Fatal("expected session id")
	}

	if s.UserID != "user-1" {
		t.Fatal("unexpected user id")
	}
}

func TestSessionGet(t *testing.T) {

	store := session.NewMemoryStore()

	service := session.NewService(
		store,
		hooks.NewRegistry(),
	)

	s, err := service.Create(
		"user-1",
	)

	if err != nil {
		t.Fatal(err)
	}

	result, err := service.Get(
		s.ID,
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.ID != s.ID {
		t.Fatal("session id mismatch")
	}
}

func TestSessionRevoke(t *testing.T) {

	store := session.NewMemoryStore()

	service := session.NewService(
		store,
		hooks.NewRegistry(),
	)

	s, err := service.Create(
		"user-1",
	)

	if err != nil {
		t.Fatal(err)
	}

	if err := service.Revoke(
		s.ID,
	); err != nil {

		t.Fatal(err)
	}

	result, err := service.Get(
		s.ID,
	)

	if err != nil {
		t.Fatal(err)
	}

	if !result.Revoked {
		t.Fatal("expected revoked session")
	}
}

func TestSessionDeleteExpired(
	t *testing.T,
) {

	store := session.NewMemoryStore()

	service := session.NewService(
		store,
		hooks.NewRegistry(),
	)

	expired := session.Session{

		ID: "session-1",

		UserID: "user-1",

		ExpiresAt: time.Now().
			Add(
				-time.Hour,
			),

		CreatedAt: time.Now().
			Add(
				-2 * time.Hour,
			),
	}

	if err := store.Create(
		expired,
	); err != nil {

		t.Fatal(err)
	}

	if err := service.DeleteExpired(
		time.Now(),
	); err != nil {

		t.Fatal(err)
	}

	_, err := service.Get(
		expired.ID,
	)

	if err == nil {

		t.Fatal(
			"expected session deleted",
		)
	}
}
