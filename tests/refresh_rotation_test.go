package authgoblue_test

import (
	"testing"

	"github.com/qwerius/authgoblue/claims"
)

func TestRefreshTokenRotation(t *testing.T) {

	agb := newTestAuthGoBlue()

	// buat session awal
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

	// refresh token baru harus memiliki session yang sama
	// karena rotation tidak membuat session baru
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

	// refresh token lama harus ditolak
	_, _, err =
		agb.Refresh.Rotate(
			oldRefresh,
		)

	if err == nil {

		t.Fatal(
			"expected old refresh token reuse error",
		)

	}
}
