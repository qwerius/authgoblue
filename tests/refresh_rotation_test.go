package authgoblue_test

import (
	"testing"

	"github.com/qwerius/authgoblue/claims"
)

func TestRefreshTokenRotation(t *testing.T) {

	agb := newTestAuthGoBlue(t)

	sess, err :=
		agb.Session.Create(
			"user-123",
		)

	if err != nil {
		t.Fatal(err)
	}

	oldRefresh, err :=
		agb.Token.GenerateRefreshToken(
			claims.Claims{
				UserID:    "user-123",
				Role:      "admin",
				SessionID: sess.ID,
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	newAccess,
		newRefresh,
		refreshClaims,
		accessExpiresAt,
		refreshExpiresAt,
		err :=
		agb.Refresh.Rotate(
			oldRefresh,
		)

	if err != nil {
		t.Fatal(err)
	}

	if newAccess == "" {

		t.Fatal(
			"expected new access token",
		)

	}

	if newRefresh == "" {

		t.Fatal(
			"expected new refresh token",
		)

	}

	if accessExpiresAt <= 0 {

		t.Fatal(
			"expected access expires at",
		)

	}

	if refreshExpiresAt <= 0 {

		t.Fatal(
			"expected refresh expires at",
		)

	}

	if refreshClaims.SessionID != sess.ID {

		t.Fatal(
			"expected same session after rotation",
		)

	}

	newClaims, err :=
		agb.Token.ParseRefreshToken(
			newRefresh,
		)

	if err != nil {
		t.Fatal(err)
	}

	if newClaims.SessionID == "" {

		t.Fatal(
			"expected refresh token session id",
		)

	}

	if newClaims.SessionID != sess.ID {

		t.Fatal(
			"expected same session after rotation",
		)

	}

	_, _, _, _, _, err =
		agb.Refresh.Rotate(
			oldRefresh,
		)

	if err == nil {

		t.Fatal(
			"expected old refresh token reuse error",
		)

	}
}
