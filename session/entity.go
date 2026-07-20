package session

import "time"

type Session struct {
	ID string

	UserID string

	// device information
	DeviceID string

	DeviceName string

	Platform string

	// request metadata
	IPAddress string

	UserAgent string

	Revoked bool

	LastSeenAt time.Time

	ExpiresAt time.Time

	CreatedAt time.Time
}
