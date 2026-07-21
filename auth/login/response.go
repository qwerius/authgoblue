package login

import (
	"github.com/google/uuid"
	"github.com/qwerius/authgoblue/auth"
)

type Response struct {
	Result *auth.AuthResult
}

type LoginResult struct {
	AccessToken string

	RefreshToken string

	UserID uuid.UUID

	Email string

	Role string
}

func (r *Response) Output() (*LoginResult, error) {

	userID, err := uuid.Parse(
		r.Result.Claims.UserID,
	)

	if err != nil {
		return nil, err
	}

	return &LoginResult{

		AccessToken: r.Result.AccessToken,

		RefreshToken: r.Result.RefreshToken,

		UserID: userID,

		Email: r.Result.Claims.Email,

		Role: r.Result.Claims.Role,
	}, nil
}
