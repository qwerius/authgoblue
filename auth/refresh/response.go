package refresh

import "github.com/qwerius/authgoblue/claims"

type Response struct {
	AccessToken string

	RefreshToken string

	Claims claims.Claims
}
