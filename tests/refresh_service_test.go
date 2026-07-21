package authgoblue_test

import (
	"context"
	"testing"
	"time"

	"github.com/qwerius/authgoblue/auth/login"
	"github.com/qwerius/authgoblue/auth/refresh"
	"github.com/qwerius/authgoblue/hooks"
	"github.com/qwerius/authgoblue/session"
	"github.com/qwerius/authgoblue/token"
)

func TestRefreshServiceRotation(t *testing.T) {

	tokenService :=
		token.NewService(
			"secret-test",
			"test",
			15*time.Minute,
			7*24*time.Hour,
		)

	hookRegistry :=
		hooks.NewRegistry()

	sessionStore :=
		session.NewMemoryStore()

	sessionService :=
		session.NewService(
			sessionStore,
			hookRegistry,
		)

	provider :=
		&mockProvider{}

	// Login membuat session pertama
	loginService :=
		login.New(
			provider,
			tokenService,
			sessionService,
			hookRegistry,
			5,
		)

	loginResult, err :=
		loginService.Execute(
			context.Background(),
			login.Request{
				Email:    "user@test.com",
				Password: "password",
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	oldRefresh :=
		loginResult.Result.RefreshToken

	oldClaims :=
		loginResult.Result.Claims

	oldSessionID :=
		oldClaims.SessionID

	refreshService :=
		refresh.New(
			tokenService,
			sessionService,
			hookRegistry,
		)

	result, err :=
		refreshService.Execute(
			context.Background(),
			refresh.Request{
				RefreshToken: oldRefresh,
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	if result.AccessToken == "" {

		t.Fatal(
			"expected new access token",
		)
	}

	if result.RefreshToken == "" {

		t.Fatal(
			"expected new refresh token",
		)
	}

	if result.RefreshToken == oldRefresh {

		t.Fatal(
			"expected refresh token rotation",
		)
	}

	// session lama harus revoked

	oldSession, err :=
		sessionService.Get(
			oldSessionID,
		)

	if err != nil {
		t.Fatal(err)
	}

	if !oldSession.Revoked {

		t.Fatal(
			"expected old session revoked",
		)
	}

	// refresh token baru harus memakai session baru

	newClaims, err :=
		tokenService.ParseRefreshToken(
			result.RefreshToken,
		)

	if err != nil {
		t.Fatal(err)
	}

	if newClaims.SessionID == "" {

		t.Fatal(
			"expected new session id",
		)
	}

	if newClaims.SessionID == oldSessionID {

		t.Fatal(
			"expected rotated session id",
		)
	}

}
