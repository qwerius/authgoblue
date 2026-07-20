package authgoblue_test

import (
	"testing"
	"time"

	"github.com/qwerius/authgoblue"
	"github.com/qwerius/authgoblue/claims"
)

func newRefreshFlowAuthGoBlue() *authgoblue.AuthGoBlue {

	return authgoblue.New(authgoblue.Config{
		Secret:          "refresh-secret-key",
		Issuer:          "refresh-service",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
	})
}

func TestRefreshTokenFlowCreatesNewAccessToken(t *testing.T) {

	agb := newRefreshFlowAuthGoBlue()

	userClaims :=
		claims.Claims{
			UserID:      "user-123",
			Username:    "alice",
			Email:       "alice@example.com",
			Role:        "admin",
			Permissions: []string{"read"},
		}

	refreshToken, err :=
		agb.Token.GenerateRefreshToken(
			userClaims,
		)

	if err != nil {
		t.Fatal(err)
	}

	parsedRefresh, err :=
		agb.Token.ParseRefreshToken(
			refreshToken,
		)

	if err != nil {
		t.Fatal(err)
	}

	err =
		agb.Token.ValidateRefreshToken(
			parsedRefresh,
		)

	if err != nil {
		t.Fatal(err)
	}

	newAccessToken, err :=
		agb.Token.GenerateAccessToken(
			claims.Claims{
				UserID:      parsedRefresh.UserID,
				Username:    parsedRefresh.Username,
				Email:       parsedRefresh.Email,
				Role:        parsedRefresh.Role,
				Permissions: parsedRefresh.Permissions,
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	parsedAccess, err :=
		agb.Token.ParseAccessToken(
			newAccessToken,
		)

	if err != nil {
		t.Fatal(err)
	}

	err =
		agb.Token.ValidateAccessToken(
			parsedAccess,
		)

	if err != nil {
		t.Fatal(err)
	}

	if parsedAccess.UserID != userClaims.UserID {

		t.Fatalf(
			"expected user id %s got %s",
			userClaims.UserID,
			parsedAccess.UserID,
		)
	}

	if parsedAccess.Role != userClaims.Role {

		t.Fatalf(
			"expected role %s got %s",
			userClaims.Role,
			parsedAccess.Role,
		)
	}

}
