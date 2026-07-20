package authgoblue_test

import (
	"context"
	"testing"
	"time"

	"github.com/qwerius/authgoblue/hooks"
	"github.com/qwerius/authgoblue/login"
	"github.com/qwerius/authgoblue/password"
	"github.com/qwerius/authgoblue/providers"
	"github.com/qwerius/authgoblue/session"
	"github.com/qwerius/authgoblue/token"
)

type benchLoginProvider struct {
	passwordHash string
}

func (b *benchLoginProvider) FindByIdentifier(
	ctx context.Context,
	identifier string,
) (*providers.User, error) {

	return &providers.User{

		ID: "user-001",

		Username: "admin",

		Email: "admin@test.com",

		PasswordHash: b.passwordHash,

		Role: "admin",

		Permissions: []string{
			"user.read",
		},
	}, nil
}

func newBenchLoginService(
	b *testing.B,
	withHook bool,
) *login.Service {

	passwordService :=
		password.NewService()

	hash, err :=
		passwordService.Hash(
			"password123",
		)

	if err != nil {

		b.Fatal(err)
	}

	hookRegistry :=
		hooks.NewRegistry()

	if withHook {

		hookRegistry.Register(

			hooks.EventAfterLogin,

			func(
				ctx context.Context,
				payload hooks.Payload,
			) error {

				return nil
			},
		)
	}

	return login.NewService(

		&benchLoginProvider{

			passwordHash: hash,
		},

		passwordService,

		token.NewService(

			"bench-secret",

			"github.com/qwerius/authgoblue",

			15*time.Minute,

			7*24*time.Hour,
		),

		session.NewService(

			session.NewMemoryStore(),
		),

		hookRegistry,
	)
}

func BenchmarkLoginWithoutHook(
	b *testing.B,
) {

	service :=
		newBenchLoginService(
			b,
			false,
		)

	ctx :=
		context.Background()

	req :=
		login.Request{

			Identifier: "admin@test.com",

			Password: "password123",
		}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		_, err :=
			service.Login(
				ctx,
				req,
			)

		if err != nil {

			b.Fatal(err)
		}
	}
}

func BenchmarkLoginWithHook(
	b *testing.B,
) {

	service :=
		newBenchLoginService(
			b,
			true,
		)

	ctx :=
		context.Background()

	req :=
		login.Request{

			Identifier: "admin@test.com",

			Password: "password123",
		}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		_, err :=
			service.Login(
				ctx,
				req,
			)

		if err != nil {

			b.Fatal(err)
		}
	}
}

func BenchmarkHookRegistryFire(
	b *testing.B,
) {

	registry :=
		hooks.NewRegistry()

	registry.Register(

		hooks.EventAfterLogin,

		func(
			ctx context.Context,
			payload hooks.Payload,
		) error {

			return nil
		},
	)

	payload :=
		hooks.Payload{

			UserID: "user-001",

			SessionID: "session-001",
		}

	ctx :=
		context.Background()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		err :=
			registry.Fire(

				ctx,

				hooks.EventAfterLogin,

				payload,
			)

		if err != nil {

			b.Fatal(err)
		}
	}
}
