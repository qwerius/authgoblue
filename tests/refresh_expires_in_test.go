package authgoblue_test

import (
	"context"
	"testing"
	"time"

	"github.com/qwerius/authgoblue"
	authRefresh "github.com/qwerius/authgoblue/auth/refresh"
	"github.com/qwerius/authgoblue/claims"
	coreRefresh "github.com/qwerius/authgoblue/refresh"
)

func TestRefreshExpiresInCalculation(t *testing.T) {

	agb, err :=
		authgoblue.New(
			authgoblue.Config{
				Secret:          "test-secret",
				Issuer:          "test",
				Provider:        &mockProvider{},
				AccessTokenTTL:  15 * time.Minute,
				RefreshTokenTTL: 7 * 24 * time.Hour,
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	session, err :=
		agb.Session.Create(
			"user-test",
		)

	if err != nil {
		t.Fatal(err)
	}

	refreshToken, err :=
		agb.Token.GenerateRefreshToken(
			claims.Claims{
				UserID:    "user-test",
				SessionID: session.ID,
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	rotateService :=
		coreRefresh.NewService(
			agb.Token,
			agb.Revoke,
			agb.Session,
		)

	refreshService :=
		authRefresh.New(
			rotateService,
			agb.Hooks,
		)

	before := time.Now().Unix()

	result, err :=
		refreshService.Execute(
			context.Background(),
			authRefresh.Request{
				RefreshToken: refreshToken,
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	after := time.Now().Unix()

	t.Log("before =", before)
	t.Log("after  =", after)

	t.Log("access expires at =", result.AccessExpiresAt)
	t.Log("access expires in =", result.AccessExpiresIn)

	t.Log("refresh expires at =", result.RefreshExpiresAt)
	t.Log("refresh expires in =", result.RefreshExpiresIn)

	if result.AccessExpiresIn < 850 ||
		result.AccessExpiresIn > 950 {

		t.Fatalf(
			"wrong access expires in: %d",
			result.AccessExpiresIn,
		)
	}

	if result.RefreshExpiresIn < 604700 ||
		result.RefreshExpiresIn > 604900 {

		t.Fatalf(
			"wrong refresh expires in: %d",
			result.RefreshExpiresIn,
		)
	}

	if result.AccessExpiresAt <= result.AccessExpiresIn {

		t.Fatal(
			"access expires at should be timestamp, not duration",
		)
	}

	if result.RefreshExpiresAt <= result.RefreshExpiresIn {

		t.Fatal(
			"refresh expires at should be timestamp, not duration",
		)
	}
}
