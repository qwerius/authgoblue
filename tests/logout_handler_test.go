package authgoblue_test

import (
	"authgoblue"
	"authgoblue/claims"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestLogoutHandlerRevokesSession(
	t *testing.T,
) {

	agb :=
		authgoblue.New(
			authgoblue.Config{
				Secret: "secret",
				Issuer: "test",
			},
		)

	app :=
		fiber.New()

	app.Post(
		"/logout",
		agb.LogoutHandler.Logout(),
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
				UserID:    "1",
				SessionID: sess.ID,
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	req :=
		httptest.NewRequest(
			"POST",
			"/logout",
			nil,
		)

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	resp, err :=
		app.Test(
			req,
		)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 204 {

		t.Fatalf(
			"expected 204 got %d",
			resp.StatusCode,
		)
	}

	// cek session sudah revoked
	session, err :=
		agb.Session.Get(
			sess.ID,
		)

	if err != nil {
		t.Fatal(err)
	}

	if !session.Revoked {

		t.Fatal(
			"session not revoked",
		)
	}
}
