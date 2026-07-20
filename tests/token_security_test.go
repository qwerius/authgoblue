package authgoblue_test

import (
	"testing"
	"time"

	"authgoblue"
	"authgoblue/claims"
)

func newSecurityTestAuthGoBlue() *authgoblue.AuthGoBlue {
	return authgoblue.New(authgoblue.Config{
		Secret:          "security-secret-key",
		Issuer:          "security-service",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
	})
}

func securityClaims() claims.Claims {
	return claims.Claims{
		UserID:      "user-123",
		Username:    "alice",
		Email:       "alice@example.com",
		Role:        "admin",
		Permissions: []string{"read", "write"},
	}
}

// Token rusak harus ditolak
func TestInvalidAccessTokenRejected(t *testing.T) {

	agb := newSecurityTestAuthGoBlue()

	_, err := agb.Token.ParseAccessToken(
		"invalid-token",
	)

	if err == nil {
		t.Fatal("expected invalid token error")
	}
}

// Secret berbeda harus ditolak
func TestAccessTokenRejectedWithDifferentSecret(t *testing.T) {

	agb1 := authgoblue.New(authgoblue.Config{
		Secret:          "secret-one",
		Issuer:          "service",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
	})

	agb2 := authgoblue.New(authgoblue.Config{
		Secret:          "secret-two",
		Issuer:          "service",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
	})

	token, err :=
		agb1.Token.GenerateAccessToken(
			securityClaims(),
		)

	if err != nil {
		t.Fatal(err)
	}

	_, err =
		agb2.Token.ParseAccessToken(token)

	if err == nil {
		t.Fatal("expected invalid signature error")
	}
}

// Issuer berbeda harus ditolak saat validation
func TestAccessTokenRejectedWithDifferentIssuer(t *testing.T) {

	agb1 := authgoblue.New(authgoblue.Config{
		Secret:          "same-secret",
		Issuer:          "issuer-one",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
	})

	agb2 := authgoblue.New(authgoblue.Config{
		Secret:          "same-secret",
		Issuer:          "issuer-two",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
	})

	token, err :=
		agb1.Token.GenerateAccessToken(
			securityClaims(),
		)

	if err != nil {
		t.Fatal(err)
	}

	parsed, err :=
		agb2.Token.ParseAccessToken(token)

	if err != nil {
		t.Fatal(err)
	}

	err =
		agb2.Token.ValidateAccessToken(parsed)

	if err == nil {
		t.Fatal("expected invalid issuer error")
	}
}

// Access token tidak boleh dipakai sebagai refresh token
func TestAccessTokenCannotBeUsedAsRefreshToken(t *testing.T) {

	agb := newSecurityTestAuthGoBlue()

	token, err :=
		agb.Token.GenerateAccessToken(
			securityClaims(),
		)

	if err != nil {
		t.Fatal(err)
	}

	parsed, err :=
		agb.Token.ParseAccessToken(token)

	if err != nil {
		t.Fatal(err)
	}

	err =
		agb.Token.ValidateRefreshToken(parsed)

	if err == nil {
		t.Fatal("expected access token rejected as refresh token")
	}
}

// Refresh token tidak boleh dipakai sebagai access token
func TestRefreshTokenCannotBeUsedAsAccessToken(t *testing.T) {

	agb := newSecurityTestAuthGoBlue()

	token, err :=
		agb.Token.GenerateRefreshToken(
			securityClaims(),
		)

	if err != nil {
		t.Fatal(err)
	}

	parsed, err :=
		agb.Token.ParseRefreshToken(token)

	if err != nil {
		t.Fatal(err)
	}

	err =
		agb.Token.ValidateAccessToken(parsed)

	if err == nil {
		t.Fatal("expected refresh token rejected as access token")
	}
}

// Expired claim harus ditolak
func TestExpiredAccessTokenRejected(t *testing.T) {

	agb := newSecurityTestAuthGoBlue()

	expiredClaims := claims.Claims{
		UserID:      "user-123",
		Username:    "alice",
		Email:       "alice@example.com",
		Role:        "admin",
		Permissions: []string{"read"},
		TokenType:   claims.TokenTypeAccess,
		Issuer:      "security-service",
		ExpiresAt:   time.Now().Add(-time.Hour).Unix(),
	}

	err :=
		agb.Token.ValidateAccessToken(
			expiredClaims,
		)

	if err == nil {
		t.Fatal("expected expired token error")
	}
}
