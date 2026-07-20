package authgoblue_test

import (
	"context"
	"testing"

	"github.com/qwerius/authgoblue/hooks"
)

func TestHookRegistry(t *testing.T) {

	registry :=
		hooks.NewRegistry()

	called := false

	registry.Register(
		hooks.EventAfterLogin,

		func(
			ctx context.Context,
			payload hooks.Payload,
		) error {

			called = true

			if payload.UserID != "user-1" {

				t.Fail()
			}

			return nil
		},
	)

	err :=
		registry.Fire(
			context.Background(),

			hooks.EventAfterLogin,

			hooks.Payload{
				UserID: "user-1",
			},
		)

	if err != nil {

		t.Fatal(err)
	}

	if !called {

		t.Fatal("hook not executed")
	}
}
