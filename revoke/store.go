package revoke

import "time"

type Store interface {
	Revoke(
		tokenID string,
		expireAt time.Time,
	) error

	IsRevoked(
		tokenID string,
	) (bool, error)

	Consume(
		tokenID string,
		expireAt time.Time,
	) (bool, error)
}
