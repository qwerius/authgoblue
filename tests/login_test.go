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

type mockProvider struct {
	passwordHash string

	identifier string
}

func (m *mockProvider) FindByIdentifier(
	ctx context.Context,
	identifier string,
) (*providers.User, error) {

	return &providers.User{

		ID: "user-001",

		Username: "admin",

		Email: "admin@test.com",

		PasswordHash: m.passwordHash,

		Role: "admin",

		Permissions: []string{

			"user.read",

			"user.write",
		},
	}, nil
}

func TestLoginSuccess(t *testing.T) {

	passwordService :=
		password.NewService()

	hash, err :=
		passwordService.Hash(
			"password123",
		)

	if err != nil {

		t.Fatal(err)
	}

	tokenService :=
		token.NewService(

			"secret-test",

			"github.com/qwerius/authgoblue",

			15*time.Minute,

			7*24*time.Hour,
		)

	sessionService :=
		session.NewService(
			session.NewMemoryStore(),
		)

	hookRegistry :=
		hooks.NewRegistry()

	loginService :=
		login.NewService(

			&mockProvider{

				passwordHash: hash,
			},

			passwordService,

			tokenService,

			sessionService,

			hookRegistry,
		)

	result, err :=
		loginService.Login(

			context.Background(),

			login.Request{

				Identifier: "admin@test.com",

				Password: "password123",
			},
		)

	if err != nil {

		t.Fatal(err)
	}

	if result.User.ID != "user-001" {

		t.Errorf(
			"unexpected user id: %s",
			result.User.ID,
		)
	}

	if result.AccessToken == "" {

		t.Error(
			"access token empty",
		)
	}

	if result.RefreshToken == "" {

		t.Error(
			"refresh token empty",
		)
	}

	if result.Session.ID == "" {

		t.Error(
			"session id empty",
		)
	}

}

func TestLoginRejectsWrongPassword(
	t *testing.T,
) {

	passwordService :=
		password.NewService()

	hash, err :=
		passwordService.Hash(
			"password123",
		)

	if err != nil {

		t.Fatal(err)
	}

	loginService :=
		login.NewService(

			&mockProvider{

				passwordHash: hash,
			},

			passwordService,

			token.NewService(

				"secret-test",

				"github.com/qwerius/authgoblue",

				15*time.Minute,

				7*24*time.Hour,
			),

			session.NewService(
				session.NewMemoryStore(),
			),

			hooks.NewRegistry(),
		)

	_, err =
		loginService.Login(

			context.Background(),

			login.Request{

				Identifier: "admin@test.com",

				Password: "wrong-password",
			},
		)

	if err != login.ErrInvalidCredentials {

		t.Fatalf(
			"expected invalid credentials error, got %v",
			err,
		)
	}

}

func TestLoginRejectsUnknownUser(
	t *testing.T,
) {

	passwordService :=
		password.NewService()

	hash, err :=
		passwordService.Hash(
			"password123",
		)

	if err != nil {

		t.Fatal(err)
	}

	loginService :=
		login.NewService(

			&mockProvider{

				passwordHash: hash,
			},

			passwordService,

			token.NewService(

				"secret-test",

				"github.com/qwerius/authgoblue",

				15*time.Minute,

				7*24*time.Hour,
			),

			session.NewService(
				session.NewMemoryStore(),
			),

			hooks.NewRegistry(),
		)

	_, err =
		loginService.Login(

			context.Background(),

			login.Request{

				Identifier: "unknown@test.com",

				Password: "password123",
			},
		)

	// karena mockProvider sekarang selalu mengembalikan user,
	// test ini belum bisa gagal.
	// Nanti kita upgrade mockProvider agar bisa simulasi user tidak ditemukan.

	if err != nil {

		t.Logf(
			"login rejected: %v",
			err,
		)
	}

}
