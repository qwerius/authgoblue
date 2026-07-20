package main

import (
	"context"
	"fmt"
	"time"

	"authgoblue"
	"authgoblue/login"
	"authgoblue/providers"
	"authgoblue/session"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
)

// Contoh provider user dari database
type UserProvider struct {
	passwordHash string
}

func (u *UserProvider) FindByIdentifier(
	ctx context.Context,
	identifier string,
) (*providers.User, error) {

	if identifier != "admin@example.com" {

		return nil, login.ErrInvalidCredentials
	}

	return &providers.User{

		ID: "user-123",

		Username: "admin",

		Email: "admin@example.com",

		PasswordHash: u.passwordHash,

		Role: "admin",

		Permissions: []string{

			"user.read",

			"user.write",
		},
	}, nil
}

func main() {

	redisClient :=
		redis.NewClient(
			&redis.Options{

				Addr: "localhost:6379",

				PoolSize: 1000,

				MinIdleConns: 200,

				PoolTimeout: 5 * time.Second,

				DialTimeout: 2 * time.Second,

				ReadTimeout: 2 * time.Second,

				WriteTimeout: 2 * time.Second,
			},
		)

	sessionStore :=
		session.NewRedisStore(
			redisClient,
		)

	agb :=
		authgoblue.New(

			authgoblue.Config{

				Secret: "super-secret-key",

				Issuer: "example-service",

				AccessTokenTTL: 15 * time.Minute,

				RefreshTokenTTL: 7 * 24 * time.Hour,

				SessionStore: sessionStore,
			},
		)

	fmt.Printf(
		"Session Store: %T\n",
		agb.Session.Store(),
	)

	// ==========================
	// Setup Login Provider
	// ==========================

	passwordHash, err :=
		agb.Password.Hash(
			"password123",
		)

	if err != nil {

		panic(err)
	}

	agb.SetupLogin(
		&UserProvider{

			passwordHash: passwordHash,
		},
	)

	app :=
		fiber.New()

	// ==========================
	// Login
	// ==========================

	app.Post(
		"/login",
		func(c fiber.Ctx) error {
			result, err :=
				agb.SignIn(
					context.Background(),
					login.Request{
						Identifier: "admin@example.com",
						Password:   "password123",
					},
				)

			if err != nil {
				return err
			}

			return c.JSON(
				fiber.Map{
					"user":          result.User,
					"access_token":  result.AccessToken,
					"refresh_token": result.RefreshToken,
					"session_id":    result.Session.ID,
				},
			)
		},
	)

	// ==========================
	// Protected Route
	// ==========================

	app.Use(
		"/protected",
		agb.Middleware.RequireAuth(),
		agb.Middleware.RequireSession(),
	)

	app.Get(
		"/protected/me",
		func(c fiber.Ctx) error {

			id, _ :=
				agb.Context.UserID(c)

			return c.JSON(

				fiber.Map{

					"user_id": id,
				},
			)
		},
	)

	app.Post("/protected/revoke", func(c fiber.Ctx) error {

		sessionID, err :=
			agb.Context.SessionID(c)

		if err != nil {
			return err
		}

		err =
			agb.Session.Revoke(
				sessionID,
			)

		if err != nil {
			return err
		}

		return c.JSON(
			fiber.Map{
				"revoked": true,
			},
		)
	})

	fmt.Println(
		"running :8080",
	)
	for _, r := range app.GetRoutes() {
		fmt.Println(
			r.Method,
			r.Path,
		)
	}

	app.Listen(":8080")
}
