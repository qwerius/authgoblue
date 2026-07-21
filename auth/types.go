package auth

import "github.com/qwerius/authgoblue/claims"

type User struct {
	ID          string
	Username    string
	Email       string
	Role        string
	Permissions []string

	Data any
}

type AuthResult struct {
	User *User

	Claims claims.Claims

	AccessToken string

	RefreshToken string
}
