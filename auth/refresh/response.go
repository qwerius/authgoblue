package refresh

import "github.com/qwerius/authgoblue/claims"

type Response struct {
	AccessToken string

	RefreshToken string

	Claims           claims.Claims
	AccessExpiresAt  int64 `json:"access_expires_at"`
	AccessExpiresIn  int64 `json:"access_expires_in"`
	RefreshExpiresAt int64 `json:"refresh_expires_at"`
	RefreshExpiresIn int64 `json:"refresh_expires_in"`
}
