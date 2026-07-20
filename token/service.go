package token

import "time"

type Service struct {
	secret []byte

	issuer string

	accessTokenTTL time.Duration

	refreshTokenTTL time.Duration
}

func NewService(
	secret string,
	issuer string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *Service {

	return &Service{

		secret: []byte(secret),

		issuer: issuer,

		accessTokenTTL: accessTokenTTL,

		refreshTokenTTL: refreshTokenTTL,
	}
}
