package authgoblue_test

import (
	"errors"
	"testing"
	"time"

	"github.com/qwerius/authgoblue"
	"github.com/qwerius/authgoblue/claims"
	"github.com/qwerius/authgoblue/token"
)

func newTestAuthGoBlue(t *testing.T) *authgoblue.AuthGoBlue {

	agb, err := authgoblue.New(
		authgoblue.Config{
			Secret: "test-secret",
			Issuer: "service-a",

			AccessTokenTTL: 15 * time.Minute,

			Provider: &mockProvider{},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	return agb
}

func TestErrorTypeInvalidIssuer(t *testing.T) {

	agb := newTestAuthGoBlue(t)

	err := agb.Token.ValidateAccessToken(
		claims.Claims{
			UserID: "user-123",

			Issuer: "service-b",

			TokenType: claims.TokenTypeAccess,

			IssuedAt: time.Now().Unix(),

			ExpiresAt: time.Now().
				Add(time.Hour).
				Unix(),
		},
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(
		err,
		token.ErrInvalidIssuer,
	) {
		t.Fatalf(
			"expected ErrInvalidIssuer, got %v",
			err,
		)
	}
}

func TestErrorTypeExpiredToken(t *testing.T) {

	agb := newTestAuthGoBlue(t)

	err := agb.Token.ValidateAccessToken(
		claims.Claims{
			UserID: "user-123",

			Issuer: "service-a",

			TokenType: claims.TokenTypeAccess,

			IssuedAt: time.Now().
				Add(-2 * time.Hour).
				Unix(),

			ExpiresAt: time.Now().
				Add(-time.Hour).
				Unix(),
		},
	)

	if !errors.Is(
		err,
		token.ErrTokenExpired,
	) {
		t.Fatalf(
			"expected ErrTokenExpired, got %v",
			err,
		)
	}
}

func TestErrorTypeMissingExpiration(t *testing.T) {

	agb := newTestAuthGoBlue(t)

	err := agb.Token.ValidateAccessToken(
		claims.Claims{
			UserID: "user-123",

			Issuer: "service-a",

			TokenType: claims.TokenTypeAccess,

			IssuedAt: time.Now().Unix(),

			ExpiresAt: 0,
		},
	)

	if !errors.Is(
		err,
		token.ErrMissingExpiration,
	) {
		t.Fatalf(
			"expected ErrMissingExpiration, got %v",
			err,
		)
	}
}
