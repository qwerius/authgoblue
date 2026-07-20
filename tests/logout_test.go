package authgoblue_test

import (
	"authgoblue"
	"authgoblue/claims"
	"testing"
)

func TestLogoutRevokesSession(
	t *testing.T,
) {

	agb :=
		authgoblue.New(
			authgoblue.Config{
				Secret: "secret",
				Issuer: "test",
			},
		)

	// buat session
	sess, err :=
		agb.Session.Create(
			"user-1",
		)

	if err != nil {
		t.Fatal(err)
	}

	token, err :=
		agb.Token.GenerateAccessToken(
			claims.Claims{
				UserID:    "user-1",
				SessionID: sess.ID,
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	err =
		agb.Logout.Logout(
			token,
		)

	if err != nil {
		t.Fatal(err)
	}

	// ambil session setelah logout
	result, err :=
		agb.Session.Get(
			sess.ID,
		)

	if err != nil {
		t.Fatal(err)
	}

	if !result.Revoked {

		t.Fatal(
			"expected session revoked",
		)
	}
}
