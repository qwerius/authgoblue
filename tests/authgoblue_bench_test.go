package authgoblue_test

import (
	"authgoblue"
	"authgoblue/claims"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"authgoblue/middleware"
	"authgoblue/session"

	"github.com/gofiber/fiber/v3"

	"github.com/redis/go-redis/v9"
)

func newBenchmarkAuthGoBlue() *authgoblue.AuthGoBlue {
	return authgoblue.New(authgoblue.Config{
		Secret:          "benchmark-secret-key",
		Issuer:          "benchmark-service",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		Header:          "Authorization",
		Prefix:          "Bearer",
		Cookie:          false,
		CookieName:      "authgoblue_token",
	})
}

func benchmarkClaims() claims.Claims {
	return claims.Claims{
		UserID:      "user-123",
		Username:    "benchmark-user",
		Email:       "benchmark@example.com",
		Role:        "admin",
		Permissions: []string{"read", "write", "delete"},
	}
}

func benchmarkAuthWithSession() (
	*authgoblue.AuthGoBlue,
	string,
) {

	agb := newBenchmarkAuthGoBlue()

	session, err :=
		agb.Session.Create(
			"user-123",
		)

	if err != nil {
		panic(err)
	}

	token, err :=
		agb.Token.GenerateAccessToken(
			claims.Claims{
				UserID:      "user-123",
				Username:    "benchmark-user",
				Email:       "benchmark@example.com",
				Role:        "admin",
				Permissions: []string{"read", "write", "delete"},
				SessionID:   session.ID,
			},
		)

	if err != nil {
		panic(err)
	}

	return agb, token
}

// Benchmark AuthGoBlue initialization
func BenchmarkNewAuthGoBlue(b *testing.B) {

	for i := 0; i < b.N; i++ {

		_ = newBenchmarkAuthGoBlue()

	}

}

// Benchmark Access Token Generation
func BenchmarkGenerateAccessToken(b *testing.B) {

	agb, _ := benchmarkAuthWithSession()

	input := benchmarkClaims()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		_, err := agb.Token.GenerateAccessToken(input)

		if err != nil {
			b.Fatal(err)
		}

	}

}

// Benchmark Access Token Parsing
func BenchmarkParseAccessToken(b *testing.B) {

	agb := newBenchmarkAuthGoBlue()

	token, err :=
		agb.Token.GenerateAccessToken(
			benchmarkClaims(),
		)

	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		_, err :=
			agb.Token.ParseAccessToken(token)

		if err != nil {
			b.Fatal(err)
		}

	}

}

// Benchmark Access Token Validation
func BenchmarkValidateAccessToken(b *testing.B) {

	agb := newBenchmarkAuthGoBlue()

	token, err :=
		agb.Token.GenerateAccessToken(
			benchmarkClaims(),
		)

	if err != nil {
		b.Fatal(err)
	}

	parsed, err :=
		agb.Token.ParseAccessToken(token)

	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		err :=
			agb.Token.ValidateAccessToken(
				parsed,
			)

		if err != nil {
			b.Fatal(err)
		}

	}

}

// Benchmark Refresh Token Generation
func BenchmarkGenerateRefreshToken(b *testing.B) {

	agb := newBenchmarkAuthGoBlue()

	input := benchmarkClaims()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		_, err :=
			agb.Token.GenerateRefreshToken(input)

		if err != nil {
			b.Fatal(err)
		}

	}

}

// Benchmark Refresh Token Parsing
func BenchmarkParseRefreshToken(b *testing.B) {

	agb := newBenchmarkAuthGoBlue()

	token, err :=
		agb.Token.GenerateRefreshToken(
			benchmarkClaims(),
		)

	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		_, err :=
			agb.Token.ParseRefreshToken(token)

		if err != nil {
			b.Fatal(err)
		}

	}

}

// Benchmark Middleware Authentication
func BenchmarkRequireAuthMiddleware(b *testing.B) {

	agb, token := benchmarkAuthWithSession()

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

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		req :=
			httptest.NewRequest(
				http.MethodGet,
				"/protected/me",
				nil,
			)

		req.Header.Set(
			"Authorization",
			"Bearer "+token,
		)

		resp, err :=
			app.Test(req)

		if err != nil {
			b.Fatal(err)
		}

		if resp.StatusCode != fiber.StatusOK {

			b.Fatalf(
				"expected 200 got %d",
				resp.StatusCode,
			)

		}

	}

}

func BenchmarkParse(b *testing.B) {

	agb := newBenchmarkAuthGoBlue()

	token, err := agb.Token.GenerateAccessToken(
		benchmarkClaims(),
	)

	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		_, err := agb.Token.Parse(token)

		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRefreshFlow(b *testing.B) {

	agb := newBenchmarkAuthGoBlue()

	input := benchmarkClaims()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		refresh, err := agb.Token.GenerateRefreshToken(input)

		if err != nil {
			b.Fatal(err)
		}

		parsed, err := agb.Token.ParseRefreshToken(refresh)

		if err != nil {
			b.Fatal(err)
		}

		_, err = agb.Token.GenerateAccessToken(parsed)

		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFullAuthentication(b *testing.B) {

	agb := newBenchmarkAuthGoBlue()

	input := benchmarkClaims()

	token, _ := agb.Token.GenerateAccessToken(input)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		parsed, err := agb.Token.ParseAccessToken(token)

		if err != nil {
			b.Fatal(err)
		}

		if err := agb.Token.ValidateAccessToken(parsed); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGenerateAccessTokenParallel(b *testing.B) {

	agb := newBenchmarkAuthGoBlue()

	input := benchmarkClaims()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {

		for pb.Next() {

			_, err := agb.Token.GenerateAccessToken(input)

			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkParseAccessTokenParallel(b *testing.B) {

	agb := newBenchmarkAuthGoBlue()

	token, _ := agb.Token.GenerateAccessToken(
		benchmarkClaims(),
	)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {

		for pb.Next() {

			_, err := agb.Token.ParseAccessToken(token)

			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkRequireAuthMiddlewareParallel(b *testing.B) {

	agb, token := benchmarkAuthWithSession()

	app := fiber.New()

	app.Use("/protected", agb.Middleware.RequireAuth())

	app.Get("/protected/me", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {

		for pb.Next() {

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
				b.Fatal(err)
			}

			if resp.StatusCode != fiber.StatusOK {
				b.Fatal(resp.StatusCode)
			}
		}
	})
}

func BenchmarkJSONClaims(b *testing.B) {

	c := claims.Claims{
		UserID:    "123",
		Email:     "test@example.com",
		Role:      "admin",
		TokenType: claims.TokenTypeAccess,
	}

	data, _ := json.Marshal(c)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		var out claims.Claims

		json.Unmarshal(
			data,
			&out,
		)
	}
}

func BenchmarkRequireAuthMiddlewareReuseRequest(b *testing.B) {

	agb, token := benchmarkAuthWithSession()

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

	req := httptest.NewRequest(
		http.MethodGet,
		"/protected/me",
		nil,
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		resp, err := app.Test(req)

		if err != nil {
			b.Fatal(err)
		}

		if resp.StatusCode != fiber.StatusOK {
			b.Fatal(resp.StatusCode)
		}

	}

}

func benchmarkSession() session.Session {

	return session.Session{
		ID: "session-benchmark-1",

		UserID: "user-123",

		ExpiresAt: time.Now().
			Add(
				time.Hour,
			),
	}
}

// Benchmark Memory Session Store Create
func BenchmarkMemorySessionCreate(
	b *testing.B,
) {

	store :=
		session.NewMemoryStore()

	svc :=
		session.NewService(
			store,
		)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		_, err :=
			svc.Create(
				"user-123",
			)

		if err != nil {
			b.Fatal(err)
		}

	}

}

// Benchmark Memory Session Get
func BenchmarkMemorySessionGet(
	b *testing.B,
) {

	store :=
		session.NewMemoryStore()

	svc :=
		session.NewService(
			store,
		)

	s, err :=
		svc.Create(
			"user-123",
		)

	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		_, err :=
			svc.Get(
				s.ID,
			)

		if err != nil {
			b.Fatal(err)
		}

	}

}

// Benchmark Redis Session Create
func BenchmarkRedisSessionCreate(
	b *testing.B,
) {

	client :=
		redis.NewClient(
			&redis.Options{
				Addr: "localhost:6379",
			},
		)

	ctx :=
		context.Background()

	if err :=
		client.Ping(ctx).Err(); err != nil {

		b.Skip(
			"redis not running",
		)
	}

	store :=
		session.NewRedisStore(
			client,
		)

	s := benchmarkSession()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		s.ID =
			"session-bench-create"

		err :=
			store.Create(
				s,
			)

		if err != nil {
			b.Fatal(err)
		}

	}

}

// Benchmark Redis Session Get
func BenchmarkRedisSessionGet(
	b *testing.B,
) {

	client :=
		redis.NewClient(
			&redis.Options{
				Addr: "localhost:6379",
			},
		)

	ctx :=
		context.Background()

	if err :=
		client.Ping(ctx).Err(); err != nil {

		b.Skip(
			"redis not running",
		)
	}

	store :=
		session.NewRedisStore(
			client,
		)

	s :=
		benchmarkSession()

	err :=
		store.Create(
			s,
		)

	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		_, err :=
			store.Get(
				s.ID,
			)

		if err != nil {
			b.Fatal(err)
		}

	}

}

// Benchmark Middleware Authentication With Redis Session
func BenchmarkRequireAuthMiddlewareRedis(
	b *testing.B,
) {

	client :=
		redis.NewClient(
			&redis.Options{
				Addr: "localhost:6379",
			},
		)

	ctx :=
		context.Background()

	if err :=
		client.Ping(ctx).Err(); err != nil {

		b.Skip(
			"redis not running",
		)
	}

	store :=
		session.NewRedisStore(
			client,
		)

	agb :=
		authgoblue.New(
			authgoblue.Config{
				Secret:       "benchmark-secret-key",
				Issuer:       "benchmark-service",
				SessionStore: store,
			},
		)

	s, err :=
		agb.Session.Create(
			"user-123",
		)

	if err != nil {
		b.Fatal(err)
	}

	token, err :=
		agb.Token.GenerateAccessToken(
			claims.Claims{
				UserID:    "user-123",
				Username:  "benchmark-user",
				Email:     "benchmark@example.com",
				Role:      "admin",
				SessionID: s.ID,
				Permissions: []string{
					"read",
					"write",
				},
			},
		)

	if err != nil {
		b.Fatal(err)
	}

	app :=
		fiber.New()

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

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		req :=
			httptest.NewRequest(
				http.MethodGet,
				"/protected/me",
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
			b.Fatal(err)
		}

		if resp.StatusCode != fiber.StatusOK {

			b.Fatalf(
				"expected 200 got %d",
				resp.StatusCode,
			)

		}

	}

}

func BenchmarkRequireAuthMiddlewareNoSession(b *testing.B) {

	agb := newBenchmarkAuthGoBlue()

	app := fiber.New()

	app.Use(
		"/protected",
		agb.Middleware.RequireAuth_With(
			middleware.Options{
				SkipSessionCheck: true,
			},
		),
	)

	app.Get(
		"/protected/me",
		func(c fiber.Ctx) error {
			return c.SendStatus(
				fiber.StatusOK,
			)
		},
	)

	token, err := agb.Token.GenerateAccessToken(
		benchmarkClaims(),
	)

	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

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
			b.Fatal(err)
		}

		if resp.StatusCode != fiber.StatusOK {
			b.Fatal(resp.StatusCode)
		}
	}
}
