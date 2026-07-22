package authgoblue_test

import (
	"context"
	"testing"
	"time"

	"github.com/qwerius/authgoblue/auth/login"
	"github.com/qwerius/authgoblue/auth/refresh"
	"github.com/qwerius/authgoblue/hooks"
	coreRefresh "github.com/qwerius/authgoblue/refresh"
	"github.com/qwerius/authgoblue/revoke"
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

	revokeStore :=
		revoke.NewMemoryStore()

	revokeService :=
		revoke.NewService(
			revokeStore,
		)

	provider :=
		&mockProvider{}

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

	oldSessionID :=
		loginResult.Result.SessionID

	rotateService :=
		coreRefresh.NewService(
			tokenService,
			revokeService,
			sessionService,
		)

	refreshService :=
		refresh.New(
			tokenService,
			rotateService,
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

	currentSession, err :=
		sessionService.Get(
			oldSessionID,
		)

	if err != nil {
		t.Fatal(err)
	}

	if currentSession.Revoked {

		t.Fatal(
			"expected session still active after refresh",
		)

	}

	oldClaims, err :=
		tokenService.ParseRefreshToken(
			oldRefresh,
		)

	if err != nil {
		t.Fatal(err)
	}

	revoked, err :=
		revokeService.IsRevoked(
			oldClaims.TokenID,
		)

	if err != nil {
		t.Fatal(err)
	}

	if !revoked {

		t.Fatal(
			"expected old refresh token revoked",
		)

	}

	newClaims, err :=
		tokenService.ParseRefreshToken(
			result.RefreshToken,
		)

	if err != nil {
		t.Fatal(err)
	}

	if newClaims.SessionID != oldSessionID {

		t.Fatal(
			"expected same session id after rotation",
		)

	}

	if result.AccessExpiresAt <= 0 {

		t.Fatal(
			"expected access expires at",
		)

	}

	if result.RefreshExpiresAt <= 0 {

		t.Fatal(
			"expected refresh expires at",
		)

	}
}
