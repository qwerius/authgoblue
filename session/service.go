package session

import (
	"time"

	"github.com/google/uuid"
)

const DefaultSessionTTL = 7 * 24 * time.Hour

type Service struct {
	store Store
}

func NewService(
	store Store,
) *Service {

	return &Service{
		store: store,
	}
}

// Store mengembalikan store yang digunakan oleh Session Service.
func (s *Service) Store() Store {
	return s.store
}

// Create membuat session baru tanpa informasi device.
func (s *Service) Create(
	userID string,
) (Session, error) {

	now := time.Now()

	session := Session{
		ID: uuid.NewString(),

		UserID: userID,

		Revoked: false,

		LastSeenAt: now,

		ExpiresAt: now.Add(
			DefaultSessionTTL,
		),

		CreatedAt: now,
	}

	err := s.store.Create(
		session,
	)

	if err != nil {
		return session, err
	}

	return session, nil
}

// CreateWithDevice membuat session baru dengan informasi device.
func (s *Service) CreateWithDevice(
	userID string,
	deviceID string,
	deviceName string,
	platform string,
	ipAddress string,
	userAgent string,
) (Session, error) {

	now := time.Now()

	session := Session{
		ID: uuid.NewString(),

		UserID: userID,

		DeviceID: deviceID,

		DeviceName: deviceName,

		Platform: platform,

		IPAddress: ipAddress,

		UserAgent: userAgent,

		Revoked: false,

		LastSeenAt: now,

		ExpiresAt: now.Add(
			DefaultSessionTTL,
		),

		CreatedAt: now,
	}

	err := s.store.Create(
		session,
	)

	if err != nil {
		return session, err
	}

	return session, nil
}

// Get mengambil session berdasarkan ID.
func (s *Service) Get(
	id string,
) (Session, error) {

	return s.store.Get(
		id,
	)
}

// GetByUserID mengambil semua session milik user.
func (s *Service) GetByUserID(
	userID string,
) ([]Session, error) {

	return s.store.GetByUserID(
		userID,
	)
}

// Revoke mencabut satu session.
//
// Session tetap ada, tetapi tidak dapat digunakan
// untuk autentikasi berikutnya.
func (s *Service) Revoke(
	id string,
) error {

	return s.store.Revoke(
		id,
	)
}

// RevokeAll mencabut seluruh session milik user.
func (s *Service) RevokeAll(
	userID string,
) error {

	return s.store.RevokeAll(
		userID,
	)
}

// CheckSession memastikan session masih valid.
//
// Digunakan oleh flow yang membutuhkan session aktif,
// seperti refresh token rotation.
func (s *Service) CheckSession(
	sessionID string,
) error {

	sess, err :=
		s.Get(
			sessionID,
		)

	if err != nil {
		return err
	}

	if sess.Revoked {
		return ErrSessionRevoked
	}

	if time.Now().After(
		sess.ExpiresAt,
	) {
		return ErrSessionExpired
	}

	return nil
}

// Touch memperbarui waktu aktivitas terakhir session.
//
// Store.Create digunakan sebagai operasi save/upsert
// sesuai kontrak Store.
func (s *Service) Touch(
	id string,
) error {

	session, err :=
		s.store.Get(
			id,
		)

	if err != nil {
		return err
	}

	session.LastSeenAt =
		time.Now()

	return s.store.Create(
		session,
	)
}

// DeleteExpired menghapus session yang sudah melewati masa aktif.
func (s *Service) DeleteExpired(
	now time.Time,
) error {

	return s.store.DeleteExpired(
		now,
	)
}

// EnforceLimit membatasi jumlah session aktif user.
//
// Jika jumlah session aktif melebihi batas,
// session tertua akan dicabut.
func (s *Service) EnforceLimit(
	userID string,
	maxSessions int,
) error {

	if maxSessions <= 0 {
		return nil
	}

	sessions, err :=
		s.store.GetByUserID(
			userID,
		)

	if err != nil {
		return err
	}

	now := time.Now()

	active := make(
		[]Session,
		0,
	)

	for _, sess := range sessions {

		if !sess.Revoked &&
			now.Before(sess.ExpiresAt) {

			active = append(
				active,
				sess,
			)
		}
	}

	if len(active) < maxSessions {
		return nil
	}

	oldest := active[0]

	for _, sess := range active {

		if sess.CreatedAt.Before(
			oldest.CreatedAt,
		) {

			oldest = sess
		}
	}

	return s.store.Revoke(
		oldest.ID,
	)
}
