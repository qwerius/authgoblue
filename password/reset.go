package password

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type ResetTokenService struct {
	duration time.Duration
}

func NewResetTokenService(
	duration time.Duration,
) *ResetTokenService {

	return &ResetTokenService{
		duration: duration,
	}
}

func (s *ResetTokenService) Generate() (*ResetToken, error) {

	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	return &ResetToken{
		Value:     hex.EncodeToString(bytes),
		ExpiredAt: time.Now().Add(s.duration),
	}, nil
}

func (s *ResetTokenService) Validate(
	token ResetToken,
) error {

	if token.Value == "" {
		return ErrInvalidResetToken
	}

	if time.Now().After(token.ExpiredAt) {
		return ErrExpiredResetToken
	}

	return nil
}
