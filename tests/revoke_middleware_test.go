package authgoblue_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"authgoblue"
	"authgoblue/claims"

	"github.com/gofiber/fiber/v3"
)

func TestRevokedSessionRejectedByMiddleware(t *testing.T) {

	agb := authgoblue.New(
		authgoblue.Config{
			Secret:          "test-secret",
			Issuer:          "test",
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 7 * 24 * time.Hour,
			Header:          "Authorization",
			Prefix:          "Bearer",
		},
	)

	app := fiber.New()

	app.Use(
		"/protected",
		agb.Middleware.RequireAuth(),
	)

	app.Get(
		"/protected/me",
		func(c fiber.Ctx) error {
			return c.SendStatus(
				fiber.StatusOK,
			)
		},
	)

	// buat session
	sess, err := agb.Session.Create(
		"user-123",
	)

	if err != nil {
		t.Fatalf(
			"create session error: %v",
			err,
		)
	}

	token, err := agb.Token.GenerateAccessToken(
		claims.Claims{
			UserID:    "user-123",
			Username:  "alice",
			Email:     "alice@example.com",
			Role:      "admin",
			SessionID: sess.ID,
		},
	)

	if err != nil {
		t.Fatalf(
			"generate token error: %v",
			err,
		)
	}

	// revoke session
	err = agb.Session.Revoke(
		sess.ID,
	)

	if err != nil {
		t.Fatalf(
			"revoke session error: %v",
			err,
		)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/protected/me",
		nil,
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf(
			"request error: %v",
			err,
		)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {

		t.Fatalf(
			"expected 401 for revoked session, got %d",
			resp.StatusCode,
		)
	}
}
