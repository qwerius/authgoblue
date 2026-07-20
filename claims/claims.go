package claims

type Claims struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id,omitempty"`
	Username  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`

	Role string `json:"role,omitempty"`

	Permissions []string `json:"permissions,omitempty"`

	TokenType TokenType `json:"token_type,omitempty"`
	TokenID   string    `json:"jti,omitempty"`

	Issuer string `json:"iss,omitempty"`

	IssuedAt  int64 `json:"iat,omitempty"`
	ExpiresAt int64 `json:"exp,omitempty"`
}
