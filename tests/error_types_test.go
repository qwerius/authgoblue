package authgoblue_test

import (
	"errors"
	"testing"
	"time"

	"authgoblue"
	"authgoblue/claims"
	"authgoblue/token"
)

func TestErrorTypeInvalidIssuer(t *testing.T) {

	agb := authgoblue.New(
		authgoblue.Config{
			Secret: "test-secret",
			Issuer: "service-a",

			AccessTokenTTL: 15 * time.Minute,
		},
	)

	err := agb.Token.ValidateAccessToken(
		claims.Claims{

			UserID: "user-123",

			Issuer: "service-b",

			TokenType: claims.TokenTypeAccess,

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

	agb := authgoblue.New(
		authgoblue.Config{
			Secret: "test-secret",
			Issuer: "service-a",

			AccessTokenTTL: 15 * time.Minute,
		},
	)

	err := agb.Token.ValidateAccessToken(
		claims.Claims{

			UserID: "user-123",

			Issuer: "service-a",

			TokenType: claims.TokenTypeAccess,

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

	agb := authgoblue.New(
		authgoblue.Config{
			Secret: "test-secret",
			Issuer: "service-a",
		},
	)

	err := agb.Token.ValidateAccessToken(
		claims.Claims{

			UserID: "user-123",

			Issuer: "service-a",

			TokenType: claims.TokenTypeAccess,

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
