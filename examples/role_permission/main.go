package main

import (
	"time"

	"github.com/qwerius/authgoblue"
	"github.com/qwerius/authgoblue/claims"

	"github.com/gofiber/fiber/v3"
)

func main() {
	agb := authgoblue.New(authgoblue.Config{
		Secret:          "super-secret-key",
		Issuer:          "example-service",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		Header:          "Authorization",
		Prefix:          "Bearer",
		Cookie:          false,
		CookieName:      "github.com/qwerius/authgoblue_token",
	})

	app := fiber.New()

	app.Get("/login", func(c fiber.Ctx) error {
		accessToken, err := agb.Token.GenerateAccessToken(claims.Claims{
			UserID:      "user-123",
			Username:    "admin-user",
			Email:       "admin@example.com",
			Role:        "admin",
			Permissions: []string{"read", "write", "delete"},
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"access_token": accessToken})
	})

	app.Use("/admin", agb.Middleware.RequireAuth())
	app.Use("/admin", agb.Middleware.RequireRole("admin"))
	app.Get("/admin/dashboard", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "welcome admin dashboard"})
	})

	app.Use("/reports", agb.Middleware.RequireAuth())
	app.Use("/reports", agb.Middleware.RequirePermission("read"))
	app.Get("/reports", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "report access granted"})
	})

	if err := app.Listen(":3001"); err != nil {
		panic(err)
	}
}
