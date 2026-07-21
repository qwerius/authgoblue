package authgoblue

func applyDefaults(config Config) Config {

	defaults := DefaultConfig()

	if config.AccessTokenTTL == 0 {
		config.AccessTokenTTL =
			defaults.AccessTokenTTL
	}

	if config.RefreshTokenTTL == 0 {
		config.RefreshTokenTTL =
			defaults.RefreshTokenTTL
	}

	if config.Header == "" {
		config.Header =
			defaults.Header
	}

	if config.Prefix == "" {
		config.Prefix =
			defaults.Prefix
	}

	if config.AccessCookieName == "" {
		config.AccessCookieName =
			defaults.AccessCookieName
	}

	if config.RefreshCookieName == "" {
		config.RefreshCookieName =
			defaults.RefreshCookieName
	}

	if config.MaxSessions == 0 {
		config.MaxSessions =
			defaults.MaxSessions
	}

	return config
}
