package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/qwerius/authgoblue"
	"github.com/qwerius/authgoblue/session"
	"github.com/redis/go-redis/v9"
)

func main() {

	redisClient := redis.NewClient(
		&redis.Options{
			Addr:         "localhost:6379",
			PoolSize:     1000,
			MinIdleConns: 200,
			PoolTimeout:  5 * time.Second,
			DialTimeout:  2 * time.Second,
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 2 * time.Second,
		},
	)

	sessionStore := session.NewRedisStore(redisClient)

	agb := authgoblue.New(authgoblue.Config{
		Secret:          "super-secret-key",
		Issuer:          "example-service",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		SessionStore:    sessionStore,
	})

	fmt.Printf("Session Store: %T\n", agb.Session.Store())

	app := fiber.New()

	app.Post("/login", func(c fiber.Ctx) error {

		return fiber.ErrNotImplemented
	})

	// ==========================
	// Protected
	// ==========================

	app.Use(
		"/protected",
		agb.Middleware.RequireAuth(),
		agb.Middleware.RequireSession(),
	)

	app.Get("/protected/me", func(c fiber.Ctx) error {

		id, err := agb.Context.UserID(c)
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"user_id": id,
		})
	})

	app.Post("/protected/revoke", func(c fiber.Ctx) error {

		sessionID, err := agb.Context.SessionID(c)
		if err != nil {
			return err
		}

		if err := agb.Session.Revoke(sessionID); err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"revoked": true,
		})
	})

	fmt.Println("running :8080")

	for _, r := range app.GetRoutes() {
		fmt.Println(r.Method, r.Path)
	}

	app.Listen(":8080")
}
