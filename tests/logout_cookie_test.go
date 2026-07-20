package authgoblue_test

import (
	"authgoblue"
	"authgoblue/claims"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestLogoutClearsCookie(
	t *testing.T,
) {

	agb :=
		authgoblue.New(
			authgoblue.Config{
				Secret: "secret",
				Issuer: "test",

				Cookie: true,

				CookieName: "authgoblue_token",
			},
		)

	app :=
		fiber.New()

	app.Post(
		"/logout",
		agb.LogoutHandler.Logout(),
	)

	sess, err :=
		agb.Session.Create(
			"1",
		)

	if err != nil {
		t.Fatal(err)
	}

	token, err :=
		agb.Token.GenerateAccessToken(
			claims.Claims{
				UserID: "1",

				SessionID: sess.ID,
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	req :=
		httptest.NewRequest(
			http.MethodPost,
			"/logout",
			nil,
		)

	req.AddCookie(
		&http.Cookie{
			Name:  "authgoblue_token",
			Value: token,
		},
	)

	resp, err :=
		app.Test(
			req,
		)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusNoContent {

		t.Fatalf(
			"expected 204 got %d",
			resp.StatusCode,
		)
	}

	cookies :=
		resp.Cookies()

	found := false

	for _, cookie := range cookies {

		if cookie.Name ==
			"authgoblue_token" {

			found = true

			if cookie.MaxAge != -1 {

				t.Fatalf(
					"expected cookie deletion",
				)
			}
		}
	}

	if !found {

		t.Fatal(
			"logout cookie not cleared",
		)
	}
}
