package authgoblue_test

import (
	"context"
	"testing"

	"github.com/qwerius/authgoblue/hooks"
)

func TestLoginFiresSessionCreatedHook(
	t *testing.T,
) {

	registry := hooks.NewRegistry()

	called := false

	registry.Register(
		hooks.EventSessionCreated,
		func(
			ctx context.Context,
			payload hooks.Payload,
		) error {

			called = true

			if payload.UserID == "" {
				t.Error("missing user id")
			}

			if payload.SessionID == "" {
				t.Error("missing session id")
			}

			return nil
		},
	)

	err := registry.Fire(
		context.Background(),
		hooks.EventSessionCreated,
		hooks.Payload{
			UserID:    "user-1",
			SessionID: "session-1",
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if !called {
		t.Error(
			"session created hook was not called",
		)
	}
}
