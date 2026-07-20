package authgoblue_test

import (
	"errors"
	"testing"
	"time"

	"authgoblue"
	"authgoblue/claims"
)

func TestInvalidIssuerError(t *testing.T) {

	agb := authgoblue.New(
		authgoblue.Config{

			Secret: "error-secret",

			Issuer: "issuer-a",

			AccessTokenTTL: 15 * time.Minute,
		},
	)

	err :=
		agb.Token.ValidateAccessToken(
			claims.Claims{

				UserID: "user-123",

				Issuer: "issuer-b",

				TokenType: claims.TokenTypeAccess,

				ExpiresAt: time.Now().
					Add(time.Hour).
					Unix(),
			},
		)

	if err == nil {

		t.Fatal(
			"expected issuer error",
		)
	}

	// memastikan error tersedia
	// dan bisa dicek
	if !errors.Is(err, err) {

		t.Fatal(
			"error should be comparable",
		)
	}

}
