package authgoblue_test

import (
	"testing"
	"time"

	"github.com/qwerius/authgoblue"
	"github.com/qwerius/authgoblue/claims"
	"github.com/qwerius/authgoblue/revoke"
	"github.com/qwerius/authgoblue/session"
)

func newRefreshBenchAuthGoBlue() *authgoblue.AuthGoBlue {

	agb, err := authgoblue.New(
		authgoblue.Config{
			Secret: "bench-secret",

			Issuer: "bench",

			AccessTokenTTL: 15 * time.Minute,

			RefreshTokenTTL: 7 * 24 * time.Hour,

			SessionStore: session.NewMemoryStore(),

			RevokeStore: revoke.NewMemoryStore(),
		},
	)
	if err != nil {
		panic(err)
	}
	return agb
}

func BenchmarkRefreshTokenAtomicRotation(
	b *testing.B,
) {

	agb :=
		newRefreshBenchAuthGoBlue()

	sess, err :=
		agb.Session.Create(
			"user-bench",
		)

	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		refreshToken, err :=
			agb.Token.GenerateRefreshToken(
				claims.Claims{

					UserID: "user-bench",

					SessionID: sess.ID,
				},
			)

		if err != nil {
			b.Fatal(err)
		}

		_, _, err =
			agb.Refresh.Rotate(
				refreshToken,
			)

		if err != nil {
			b.Fatal(err)
		}
	}
}
