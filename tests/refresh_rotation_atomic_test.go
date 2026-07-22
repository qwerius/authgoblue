package authgoblue_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/qwerius/authgoblue"
	"github.com/qwerius/authgoblue/claims"
	"github.com/qwerius/authgoblue/refresh"
)

func newRefreshTestAuthGoBlue() *authgoblue.AuthGoBlue {

	agb, err :=
		authgoblue.New(
			authgoblue.Config{
				Secret: "test-secret",

				Issuer: "test",

				Provider: &mockProvider{},

				AccessTokenTTL: 15 * time.Minute,

				RefreshTokenTTL: 7 * 24 * time.Hour,
			},
		)

	if err != nil {
		panic(err)
	}

	return agb
}

func TestRefreshTokenAtomicRotation(t *testing.T) {

	agb :=
		newRefreshTestAuthGoBlue()

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
				UserID: "user-1",

				SessionID: sess.ID,
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup

	var mu sync.Mutex

	success := 0

	reuseError := 0

	workers := 20

	for i := 0; i < workers; i++ {

		wg.Add(1)

		go func() {

			defer wg.Done()

			_, _, _, _, _, err :=
				agb.Refresh.Rotate(
					refreshToken,
				)

			mu.Lock()

			defer mu.Unlock()

			switch {

			case err == nil:

				success++

			case errors.Is(
				err,
				refresh.ErrRefreshTokenReuse,
			):

				reuseError++

			default:

				t.Errorf(
					"unexpected error: %v",
					err,
				)
			}

		}()
	}

	wg.Wait()

	if success != 1 {

		t.Fatalf(
			"expected exactly 1 successful rotation, got %d",
			success,
		)
	}

	if reuseError != workers-1 {

		t.Fatalf(
			"expected %d reuse errors, got %d",
			workers-1,
			reuseError,
		)
	}
}
