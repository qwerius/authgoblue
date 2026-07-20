package session

import "time"

// Store adalah abstraksi penyimpanan session.
//
// Create digunakan untuk membuat session baru
// sekaligus menyimpan perubahan session yang sudah ada
// (upsert).
type Store interface {

	// Create menyimpan session.
	// Jika ID sudah ada, data lama akan diganti.
	Create(
		session Session,
	) error

	// Get mengambil session berdasarkan ID.
	Get(
		id string,
	) (Session, error)

	// GetByUserID mengambil seluruh session milik user.
	GetByUserID(
		userID string,
	) ([]Session, error)

	// Revoke mencabut satu session.
	Revoke(
		id string,
	) error

	// RevokeAll mencabut seluruh session milik user.
	RevokeAll(
		userID string,
	) error

	// DeleteExpired menghapus session yang sudah expired.
	DeleteExpired(
		now time.Time,
	) error
}
