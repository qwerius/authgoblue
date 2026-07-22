package authgoblue_test

import (
	"sync"
	"testing"
	"time"

	"github.com/qwerius/authgoblue"
	"github.com/qwerius/authgoblue/claims"
)

func TestAuthGoBlueConcurrentTokenGeneration(
	t *testing.T,
) {

	agb, err :=
		authgoblue.New(
			authgoblue.Config{

				Secret: "race-secret",

				Issuer: "race-service",

				AccessTokenTTL: 15 * time.Minute,
				Provider:       &mockProvider{},
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {

		wg.Add(1)

		go func(id int) {

			defer wg.Done()

			_, err :=
				agb.Token.GenerateAccessToken(
					claims.Claims{

						UserID: "user",

						Username: "worker",
					},
				)

			if err != nil {

				t.Errorf(
					"goroutine %d failed: %v",
					id,
					err,
				)
			}

		}(i)

	}

	wg.Wait()

}
