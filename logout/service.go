package logout

import (
	"authgoblue/session"
	"authgoblue/token"
)

type Service struct {
	token   *token.Service
	session *session.Service
}

func NewService(
	tokenService *token.Service,
	sessionService *session.Service,
) *Service {

	return &Service{
		token:   tokenService,
		session: sessionService,
	}
}

func (s *Service) Logout(
	accessToken string,
) error {

	if accessToken == "" {
		return ErrMissingToken
	}

	claims, err := s.token.ParseAccessToken(
		accessToken,
	)

	if err != nil {
		return ErrInvalidToken
	}

	if err := s.token.ValidateAccessToken(
		claims,
	); err != nil {
		return err
	}

	if claims.SessionID == "" {
		return ErrMissingSessionID
	}

	return s.session.Revoke(
		claims.SessionID,
	)
}
