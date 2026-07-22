package refresh

import "github.com/qwerius/authgoblue/claims"

type Response struct {
	AccessToken string

	RefreshToken string

	Claims          claims.Claims
	AccessExpiresAt int64

	RefreshExpiresAt int64
}
