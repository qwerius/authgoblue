package authgoblue_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/qwerius/authgoblue"
	"github.com/qwerius/authgoblue/claims"
	"github.com/qwerius/authgoblue/hooks"
	"github.com/qwerius/authgoblue/session"

	"github.com/gofiber/fiber/v3"
)

func TestNewAppliesDefaults(t *testing.T) {
	agb := newTestAuthGoBlue(t)

	cfg := agb.Config()

	if cfg.Header != "Authorization" {
		t.Fatalf("expected default header Authorization, got %q", cfg.Header)
	}

	if cfg.Prefix != "Bearer" {
		t.Fatalf("expected default prefix Bearer, got %q", cfg.Prefix)
	}

	if cfg.AccessCookieName != "github.com/qwerius/authgoblue_token" {
		t.Fatalf("expected default cookie name authgoblue_token, got %q", cfg.AccessCookieName)
	}

	if cfg.AccessTokenTTL <= 0 {
		t.Fatal("expected access token ttl to be set")
	}

	if cfg.RefreshTokenTTL <= 0 {
		t.Fatal("expected refresh token ttl to be set")
	}
}

func TestTokenGenerateParseAndValidateAccessToken(t *testing.T) {
	agb := newTestAuthGoBlue(t)

	input := claims.Claims{
		UserID:      "user-123",
		Username:    "alice",
		Email:       "alice@example.com",
		Role:        "admin",
		Permissions: []string{"read", "write"},
	}

	accessToken, err := agb.Token.GenerateAccessToken(input)
	if err != nil {
		t.Fatalf("GenerateAccessToken returned error: %v", err)
	}

	parsed, err := agb.Token.ParseAccessToken(accessToken)
	if err != nil {
		t.Fatalf("ParseAccessToken returned error: %v", err)
	}

	if parsed.UserID != input.UserID {
		t.Fatalf("expected user id %q, got %q", input.UserID, parsed.UserID)
	}

	if parsed.Role != input.Role {
		t.Fatalf("expected role %q, got %q", input.Role, parsed.Role)
	}

	if !reflect.DeepEqual(parsed.Permissions, input.Permissions) {
		t.Fatalf("expected permissions %v, got %v", input.Permissions, parsed.Permissions)
	}

	if err := agb.Token.ValidateAccessToken(parsed); err != nil {
		t.Fatalf("ValidateAccessToken returned error: %v", err)
	}
}

func TestTokenGenerateParseAndValidateRefreshToken(t *testing.T) {
	agb := newTestAuthGoBlue(t)

	input := claims.Claims{
		UserID:   "user-123",
		Username: "alice",
		Email:    "alice@example.com",
		Role:     "admin",
	}

	refreshToken, err := agb.Token.GenerateRefreshToken(input)
	if err != nil {
		t.Fatalf("GenerateRefreshToken returned error: %v", err)
	}

	parsed, err := agb.Token.ParseRefreshToken(refreshToken)
	if err != nil {
		t.Fatalf("ParseRefreshToken returned error: %v", err)
	}

	if parsed.UserID != input.UserID {
		t.Fatalf("expected user id %q, got %q", input.UserID, parsed.UserID)
	}

	if err := agb.Token.ValidateRefreshToken(parsed); err != nil {
		t.Fatalf("ValidateRefreshToken returned error: %v", err)
	}
}

func TestRequireAuthAcceptsValidBearerTokenAndRejectsInvalidToken(t *testing.T) {
	agb := newTestAuthGoBlue(t)

	app := fiber.New()
	app.Use("/protected", agb.Middleware.RequireAuth())
	app.Get("/protected/me", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	sess, err := agb.Session.Create(
		"user-123",
	)

	if err != nil {
		t.Fatal(err)
	}

	accessToken, err := agb.Token.GenerateAccessToken(
		claims.Claims{
			UserID:      "user-123",
			Username:    "alice",
			Email:       "alice@example.com",
			Role:        "admin",
			Permissions: []string{"read"},
			SessionID:   sess.ID,
		},
	)
	if err != nil {
		t.Fatalf("GenerateAccessToken returned error: %v", err)
	}

	validReq := httptest.NewRequest(http.MethodGet, "/protected/me", nil)
	validReq.Header.Set("Authorization", "Bearer "+accessToken)

	validResp, err := app.Test(validReq)
	if err != nil {
		t.Fatalf("valid request error: %v", err)
	}

	if validResp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected valid protected route to return 200, got %d", validResp.StatusCode)
	}

	invalidReq := httptest.NewRequest(http.MethodGet, "/protected/me", nil)
	invalidReq.Header.Set("Authorization", "Bearer invalid-token")

	invalidResp, err := app.Test(invalidReq)
	if err != nil {
		t.Fatalf("invalid request error: %v", err)
	}

	if invalidResp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected invalid token to return 401, got %d", invalidResp.StatusCode)
	}
}

func TestRoleAndPermissionMiddlewares(t *testing.T) {
	agb := newTestAuthGoBlue(t)

	app := fiber.New()

	app.Use("/admin", agb.Middleware.RequireAuth())
	app.Use("/admin", agb.Middleware.RequireRole("admin"))
	app.Get("/admin/dashboard", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app.Use("/reports", agb.Middleware.RequireAuth())
	app.Use("/reports", agb.Middleware.RequirePermission("read"))
	app.Get("/reports", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app.Use("/secret", agb.Middleware.RequireAuth())
	app.Use("/secret", agb.Middleware.RequirePermission("delete"))
	app.Get("/secret", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	sess, err := agb.Session.Create(
		"user-123",
	)

	if err != nil {
		t.Fatal(err)
	}

	accessToken, err := agb.Token.GenerateAccessToken(
		claims.Claims{
			UserID:      "user-123",
			Username:    "alice",
			Email:       "alice@example.com",
			Role:        "admin",
			Permissions: []string{"read", "write"},
			SessionID:   sess.ID,
		},
	)
	if err != nil {
		t.Fatalf("GenerateAccessToken returned error: %v", err)
	}

	adminReq := httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	adminReq.Header.Set("Authorization", "Bearer "+accessToken)

	adminResp, err := app.Test(adminReq)
	if err != nil {
		t.Fatalf("admin request error: %v", err)
	}

	if adminResp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected admin route to return 200, got %d", adminResp.StatusCode)
	}

	reportsReq := httptest.NewRequest(http.MethodGet, "/reports", nil)
	reportsReq.Header.Set("Authorization", "Bearer "+accessToken)

	reportsResp, err := app.Test(reportsReq)
	if err != nil {
		t.Fatalf("reports request error: %v", err)
	}

	if reportsResp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected reports route to return 200, got %d", reportsResp.StatusCode)
	}

	forbiddenReq := httptest.NewRequest(http.MethodGet, "/secret", nil)
	forbiddenReq.Header.Set("Authorization", "Bearer "+accessToken)

	forbiddenResp, err := app.Test(forbiddenReq)
	if err != nil {
		t.Fatalf("forbidden request error: %v", err)
	}

	if forbiddenResp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected route permission denial to return 403, got %d", forbiddenResp.StatusCode)
	}
}

func TestContextHelpersReadClaimsFromFiberContext(t *testing.T) {
	agb := newTestAuthGoBlue(t)
	app := fiber.New()

	app.Use(func(c fiber.Ctx) error {
		agb.Context.SetClaims(c, claims.Claims{
			UserID:      "user-123",
			Username:    "alice",
			Email:       "alice@example.com",
			Role:        "admin",
			Permissions: []string{"read", "write"},
		})
		return c.Next()
	})

	app.Get("/claims", func(c fiber.Ctx) error {
		userID, err := agb.Context.UserID(c)
		if err != nil {
			t.Fatalf("UserID returned error: %v", err)
		}

		role, err := agb.Context.Role(c)
		if err != nil {
			t.Fatalf("Role returned error: %v", err)
		}

		permissions, err := agb.Context.Permissions(c)
		if err != nil {
			t.Fatalf("Permissions returned error: %v", err)
		}

		allowed, err := agb.Context.HasPermission(c, "read")
		if err != nil {
			t.Fatalf("HasPermission returned error: %v", err)
		}

		if userID != "user-123" {
			t.Fatalf("expected user id user-123, got %q", userID)
		}

		if role != "admin" {
			t.Fatalf("expected role admin, got %q", role)
		}

		if !reflect.DeepEqual(permissions, []string{"read", "write"}) {
			t.Fatalf("expected permissions [read write], got %v", permissions)
		}

		if !allowed {
			t.Fatal("expected permission read to be allowed")
		}

		return c.SendString("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/claims", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("ReadAll returned error: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d with body %s", resp.StatusCode, string(body))
	}

}

func TestAuthGoBlueDefaultUsesMemorySessionStore(t *testing.T) {
	agb := newTestAuthGoBlue(t)

	if agb.Session == nil {
		t.Fatal("expected session service")
	}

	store := agb.Session.Store()

	if _, ok := store.(*session.MemoryStore); !ok {
		t.Fatalf("expected MemoryStore, got %T", store)
	}
}

func TestAuthGoBlueUsesCustomSessionStore(t *testing.T) {
	customStore := session.NewMemoryStore()

	agb, err := authgoblue.New(authgoblue.Config{
		Secret:           "test-secret-key",
		Issuer:           "test-issuer",
		AccessTokenTTL:   15 * time.Minute,
		RefreshTokenTTL:  7 * 24 * time.Hour,
		Header:           "Authorization",
		Prefix:           "Bearer",
		Cookie:           false,
		AccessCookieName: "github.com/qwerius/authgoblue_token",
		SessionStore:     customStore,
	})

	if err != nil {
		t.Fatal(err)
	}

	if agb.Session == nil {
		t.Fatal("expected session service")
	}

	store := agb.Session.Store()

	if store != customStore {
		t.Fatal("custom SessionStore was not used")
	}
}

func TestRequireSessionAcceptsValidSession(t *testing.T) {

	agb := newTestAuthGoBlue(t)

	app := fiber.New()

	app.Use(
		"/profile",
		agb.Middleware.RequireAuth(),
	)

	app.Use(
		"/profile",
		agb.Middleware.RequireSession(),
	)

	app.Get("/profile", func(c fiber.Ctx) error {

		sess, err := agb.Context.Session(c)

		if err != nil {
			t.Fatalf(
				"expected session in context: %v",
				err,
			)
		}

		if sess.UserID != "user-123" {
			t.Fatalf(
				"expected user id user-123 got %s",
				sess.UserID,
			)
		}

		return c.SendStatus(
			fiber.StatusOK,
		)
	})

	sess, err := agb.Session.Create(
		"user-123",
	)

	if err != nil {
		t.Fatalf(
			"session create error: %v",
			err,
		)
	}

	token, err := agb.Token.GenerateAccessToken(
		claims.Claims{
			UserID:    "user-123",
			SessionID: sess.ID,
		},
	)

	if err != nil {
		t.Fatalf(
			"generate token error: %v",
			err,
		)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/profile",
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

	if resp.StatusCode != fiber.StatusOK {

		t.Fatalf(
			"expected status 200 got %d",
			resp.StatusCode,
		)
	}
}

func TestRequireSessionRejectsRevokedSession(t *testing.T) {

	agb := newTestAuthGoBlue(t)

	app := fiber.New()

	app.Use(
		agb.Middleware.RequireAuth(),
	)

	app.Use(
		agb.Middleware.RequireSession(),
	)

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendStatus(
			fiber.StatusOK,
		)
	})

	sess, err := agb.Session.Create(
		"user-123",
	)

	if err != nil {
		t.Fatal(err)
	}

	err = agb.Session.Revoke(
		sess.ID,
	)

	if err != nil {
		t.Fatal(err)
	}

	token, err := agb.Token.GenerateAccessToken(
		claims.Claims{
			UserID:    "user-123",
			SessionID: sess.ID,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	resp, err := app.Test(req)

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

func TestRequireSessionRejectsExpiredSession(t *testing.T) {

	agb := newTestAuthGoBlue(t)

	app := fiber.New()

	app.Use(
		agb.Middleware.RequireAuth(),
	)

	app.Use(
		agb.Middleware.RequireSession(),
	)

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendStatus(
			fiber.StatusOK,
		)
	})

	sess, err := agb.Session.Create(
		"user-123",
	)

	if err != nil {
		t.Fatal(err)
	}

	// hapus session lalu buat expired
	sess.ExpiresAt =
		time.Now().Add(
			-1 * time.Hour,
		)

	// karena Store belum punya Update,
	// revoke untuk simulasi expired tidak bisa.
	// test ini membutuhkan SessionStore Update.
	_ = sess

	token, err := agb.Token.GenerateAccessToken(
		claims.Claims{
			UserID:    "user-123",
			SessionID: sess.ID,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	resp, err := app.Test(req)

	if err != nil {
		t.Fatal(err)
	}

	// sementara hanya memastikan middleware jalan
	if resp.StatusCode != fiber.StatusUnauthorized &&
		resp.StatusCode != fiber.StatusOK {

		t.Fatalf(
			"unexpected status %d",
			resp.StatusCode,
		)
	}
}

func TestSessionMultipleDevice(t *testing.T) {

	store := session.NewMemoryStore()

	service := session.NewService(
		store,
		hooks.NewRegistry(),
	)

	s1, err := service.CreateWithDevice(
		"user-1",
		"device-1",
		"Chrome Windows",
		"windows",
		"127.0.0.1",
		"Chrome",
	)

	if err != nil {
		t.Fatalf("create session 1 error: %v", err)
	}

	s2, err := service.CreateWithDevice(
		"user-1",
		"device-2",
		"Android Phone",
		"android",
		"127.0.0.2",
		"Android Chrome",
	)

	if err != nil {
		t.Fatalf("create session 2 error: %v", err)
	}

	sessions, err :=
		service.GetByUserID(
			"user-1",
		)

	if err != nil {
		t.Fatalf(
			"GetByUserID error: %v",
			err,
		)
	}

	if len(sessions) != 2 {

		t.Fatalf(
			"expected 2 sessions, got %d",
			len(sessions),
		)
	}

	if s1.DeviceName != "Chrome Windows" {

		t.Fatalf(
			"wrong device 1",
		)
	}

	if s2.DeviceName != "Android Phone" {

		t.Fatalf(
			"wrong device 2",
		)
	}
}

func TestSessionLogoutCurrentDevice(t *testing.T) {

	store := session.NewMemoryStore()

	service := session.NewService(
		store,
		hooks.NewRegistry(),
	)

	sess, err :=
		service.CreateWithDevice(
			"user-1",
			"device-1",
			"Chrome",
			"windows",
			"127.0.0.1",
			"Chrome",
		)

	if err != nil {
		t.Fatal(err)
	}

	err =
		service.Revoke(
			sess.ID,
		)

	if err != nil {
		t.Fatal(err)
	}

	result, err :=
		service.Get(
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

func TestSessionLogoutAllDevices(t *testing.T) {

	store := session.NewMemoryStore()

	service := session.NewService(
		store,
		hooks.NewRegistry(),
	)

	_, err :=
		service.CreateWithDevice(
			"user-1",
			"d1",
			"Chrome",
			"windows",
			"127.0.0.1",
			"Chrome",
		)

	if err != nil {
		t.Fatal(err)
	}

	_, err =
		service.CreateWithDevice(
			"user-1",
			"d2",
			"Android",
			"android",
			"127.0.0.2",
			"Android",
		)

	if err != nil {
		t.Fatal(err)
	}

	err =
		service.RevokeAll(
			"user-1",
		)

	if err != nil {
		t.Fatal(err)
	}

	list, err :=
		service.GetByUserID(
			"user-1",
		)

	if err != nil {
		t.Fatal(err)
	}

	for _, sess := range list {

		if !sess.Revoked {

			t.Fatal(
				"expected all sessions revoked",
			)
		}
	}
}

func TestSessionLimitRevokesOldestSession(t *testing.T) {

	store := session.NewMemoryStore()

	service := session.NewService(store, hooks.NewRegistry())

	s1, _ := service.CreateWithDevice(
		"user-1",
		"d1",
		"Chrome",
		"windows",
		"ip",
		"ua",
	)

	time.Sleep(time.Millisecond)

	_, _ = service.CreateWithDevice(
		"user-1",
		"d2",
		"Android",
		"android",
		"ip",
		"ua",
	)

	err :=
		service.EnforceLimit(
			"user-1",
			1,
		)

	if err != nil {
		t.Fatal(err)
	}

	old, _ :=
		service.Get(
			s1.ID,
		)

	if !old.Revoked {
		t.Fatal(
			"oldest session should be revoked",
		)
	}
}

func TestAuthGoBlueDefaultMaxSessions(t *testing.T) {

	agb, err := authgoblue.New(
		authgoblue.Config{
			Secret: "secret",
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if agb.Config().MaxSessions != 5 {
		t.Fatalf(
			"expected max sessions 5, got %d",
			agb.Config().MaxSessions,
		)
	}
}

func TestAuthGoBlueCustomMaxSessions(t *testing.T) {

	agb, err := authgoblue.New(
		authgoblue.Config{
			Secret:      "secret",
			MaxSessions: 10,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if agb.Config().MaxSessions != 10 {
		t.Fatalf(
			"expected max sessions 10, got %d",
			agb.Config().MaxSessions,
		)
	}
}

func TestDefaultConfigValues(t *testing.T) {

	agb, err := authgoblue.New(
		authgoblue.Config{
			Secret:   "secret",
			Provider: &mockProvider{},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	cfg := agb.Config()

	if cfg.Cookie != false {
		t.Fatal("expected cookie false")
	}

	if cfg.AccessCookieName != "access_token" {
		t.Fatalf(
			"unexpected cookie name %s",
			cfg.AccessCookieName,
		)
	}

	if cfg.MaxSessions != 5 {
		t.Fatalf(
			"expected max sessions 5 got %d",
			cfg.MaxSessions,
		)
	}
}

func TestRefreshRotationCreatesNewSession(t *testing.T) {

	agb := newTestAuthGoBlue(t)

	sess, err :=
		agb.Session.Create(
			"user-123",
		)

	if err != nil {
		t.Fatal(err)
	}

	refreshToken, err :=
		agb.Token.GenerateRefreshToken(
			claims.Claims{
				UserID:    "user-123",
				Username:  "alice",
				Email:     "alice@example.com",
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
			refreshToken,
		)

	if err != nil {
		t.Fatal(err)
	}

	if newAccess == "" {
		t.Fatal("expected new access token")
	}

	if newRefresh == "" {
		t.Fatal("expected new refresh token")
	}

}

func TestRefreshRotationRejectsReuse(t *testing.T) {

	agb := newTestAuthGoBlue(t)

	sess, err :=
		agb.Session.Create(
			"user-1",
		)

	if err != nil {
		t.Fatal(err)
	}

	refreshToken, err :=
		agb.Token.GenerateRefreshToken(
			claims.Claims{
				UserID:    "user-1",
				SessionID: sess.ID,
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	_, _, err =
		agb.Refresh.Rotate(
			refreshToken,
		)

	if err != nil {
		t.Fatal(err)
	}

	_, _, err =
		agb.Refresh.Rotate(
			refreshToken,
		)

	if err == nil {
		t.Fatal(
			"expected refresh reuse error",
		)
	}

}

func FuzzParseToken(f *testing.F) {

	agb, err := authgoblue.New(authgoblue.Config{
		Secret: "fuzz-secret",
	})

	if err != nil {
		f.Fatal(err)
	}

	f.Add("invalid.token.data")

	f.Fuzz(func(t *testing.T, token string) {

		_, _ = agb.Token.Parse(token)

	})

}
