package authgoblue_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qwerius/authgoblue"
	"github.com/qwerius/authgoblue/claims"

	"github.com/gofiber/fiber/v3"
)

func newMiddlewareOptionTestAuthGoBlue() *authgoblue.AuthGoBlue {

	agb, err :=
		authgoblue.New(
			authgoblue.Config{

				Secret: "test-secret",

				Issuer: "test",

				Provider: &mockProvider{},

				Header: "X-Auth-Token",

				Prefix: "Token",

				Cookie: false,
			},
		)

	if err != nil {
		panic(err)
	}

	return agb
}

func TestRequireAuthWithCustomHeaderAndPrefix(t *testing.T) {

	agb := newMiddlewareOptionTestAuthGoBlue()

	app := fiber.New()

	app.Use(
		"/private",
		agb.Middleware.RequireAuth(),
	)

	app.Get(
		"/private",
		func(c fiber.Ctx) error {

			return c.SendStatus(
				fiber.StatusOK,
			)
		},
	)

	sess, err :=
		agb.Session.Create(
			"user-123",
		)

	if err != nil {
		t.Fatal(err)
	}

	token, err :=
		agb.Token.GenerateAccessToken(
			claims.Claims{
				UserID:    "user-123",
				Role:      "admin",
				SessionID: sess.ID,
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	req :=
		httptest.NewRequest(
			http.MethodGet,
			"/private",
			nil,
		)

	req.Header.Set(
		"X-Auth-Token",
		"Token "+token,
	)

	resp, err :=
		app.Test(req)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusOK {

		t.Fatalf(
			"expected 200 got %d",
			resp.StatusCode,
		)
	}

}

func TestRequireAuthRejectsWrongPrefix(t *testing.T) {

	agb := newMiddlewareOptionTestAuthGoBlue()

	app := fiber.New()

	app.Use(
		"/private",
		agb.Middleware.RequireAuth(),
	)

	app.Get(
		"/private",
		func(c fiber.Ctx) error {

			return c.SendStatus(
				fiber.StatusOK,
			)
		},
	)

	sess, err :=
		agb.Session.Create(
			"user-123",
		)

	if err != nil {
		t.Fatal(err)
	}

	token, err :=
		agb.Token.GenerateAccessToken(
			claims.Claims{
				UserID:    "user-123",
				SessionID: sess.ID,
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	req :=
		httptest.NewRequest(
			http.MethodGet,
			"/private",
			nil,
		)

	// Salah prefix
	req.Header.Set(
		"X-Auth-Token",
		"Bearer "+token,
	)

	resp, err :=
		app.Test(req)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {

		t.Fatalf(
			"expected 401 got %d",
			resp.StatusCode,
		)
	}

}

func TestRequireAuthRejectsMissingHeader(t *testing.T) {

	agb := newMiddlewareOptionTestAuthGoBlue()

	app := fiber.New()

	app.Use(
		"/private",
		agb.Middleware.RequireAuth(),
	)

	app.Get(
		"/private",
		func(c fiber.Ctx) error {

			return c.SendStatus(
				fiber.StatusOK,
			)
		},
	)

	req :=
		httptest.NewRequest(
			http.MethodGet,
			"/private",
			nil,
		)

	resp, err :=
		app.Test(req)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {

		t.Fatalf(
			"expected 401 got %d",
			resp.StatusCode,
		)
	}

}
