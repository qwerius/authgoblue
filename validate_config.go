package authgoblue

func ValidateConfig(config Config) error {

	if config.Secret == "" {
		return ErrSecretRequired
	}

	if config.Issuer == "" {
		return ErrIssuerRequired
	}

	if config.Provider == nil {
		return ErrProviderRequired
	}

	if config.SessionStore == nil {
		return ErrSessionStoreRequired
	}

	if config.RevokeStore == nil {
		return ErrRevokeStoreRequired
	}

	if config.AccessTokenTTL <= 0 {
		return ErrInvalidAccessTokenTTL
	}

	if config.RefreshTokenTTL <= 0 {
		return ErrInvalidRefreshTokenTTL
	}

	if config.MaxSessions <= 0 {
		return ErrInvalidMaxSessions
	}

	return nil
}
