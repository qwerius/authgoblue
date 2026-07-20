package claims

import "time"

func (c Claims) Expired() bool {

	if c.ExpiresAt == 0 {
		return false
	}

	return time.Now().Unix() > c.ExpiresAt
}

func (c Claims) IssuedTime() time.Time {

	if c.IssuedAt == 0 {
		return time.Time{}
	}

	return time.Unix(c.IssuedAt, 0)
}

func (c Claims) ExpirationTime() time.Time {

	if c.ExpiresAt == 0 {
		return time.Time{}
	}

	return time.Unix(c.ExpiresAt, 0)
}
