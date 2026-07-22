package authgoblue_test

import (
	"errors"
	"testing"
	"time"

	"github.com/qwerius/authgoblue"
	"github.com/qwerius/authgoblue/claims"
	"github.com/qwerius/authgoblue/token"
)

func TestInvalidIssuerError(t *testing.T) {

	agb, err := authgoblue.New(
		authgoblue.Config{
			Secret: "error-secret",
			Issuer: "issuer-a",

			AccessTokenTTL: 15 * time.Minute,

			Provider: &mockProvider{},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	validateErr := agb.Token.ValidateAccessToken(
		claims.Claims{
			UserID: "user-123",

			Issuer: "issuer-b",

			TokenType: claims.TokenTypeAccess,

			IssuedAt: time.Now().Unix(),

			ExpiresAt: time.Now().
				Add(time.Hour).
				Unix(),
		},
	)

	if validateErr == nil {
		t.Fatal(
			"expected issuer error",
		)
	}

	if !errors.Is(
		validateErr,
		token.ErrInvalidIssuer,
	) {
		t.Fatalf(
			"expected invalid issuer error, got %v",
			validateErr,
		)
	}
}
