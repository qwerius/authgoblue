package login

import "github.com/google/uuid"

type Response struct {
	Result *Result
}

type Result struct {
	AccessToken string

	RefreshToken string

	UserID uuid.UUID

	Email string

	Role string

	SessionID string
}
