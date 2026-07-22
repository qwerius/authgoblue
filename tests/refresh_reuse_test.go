package authgoblue_test

import (
	"testing"

	"github.com/qwerius/authgoblue"
	"github.com/qwerius/authgoblue/claims"
	"github.com/qwerius/authgoblue/refresh"

	"github.com/stretchr/testify/require"
)

func TestRefreshTokenReuseDetection(t *testing.T) {

	agb, err :=
		authgoblue.New(
			authgoblue.Config{
				Secret: "test-secret",
				Issuer: "test",
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	session, err :=
		agb.Session.Create(
			"user-1",
		)

	require.NoError(
		t,
		err,
	)

	refreshToken, err :=
		agb.Token.GenerateRefreshToken(
			claims.Claims{
				UserID:    "user-1",
				SessionID: session.ID,
			},
		)

	require.NoError(
		t,
		err,
	)

	// pertama berhasil
	_, _, _, err =
		agb.Refresh.Rotate(
			refreshToken,
		)

	require.NoError(
		t,
		err,
	)

	// token lama dipakai lagi
	_, _, _, err =
		agb.Refresh.Rotate(
			refreshToken,
		)

	require.ErrorIs(
		t,
		err,
		refresh.ErrRefreshTokenReuse,
	)

}
