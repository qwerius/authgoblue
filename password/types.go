package password

import "time"

type ResetToken struct {
	Value     string
	ExpiredAt time.Time
}
