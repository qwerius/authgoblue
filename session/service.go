package session

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/qwerius/authgoblue/hooks"
)

const DefaultSessionTTL = 7 * 24 * time.Hour

type Service struct {
	store Store

	hooks *hooks.Registry
}

func NewService(
	store Store,
	hookRegistry *hooks.Registry,
) *Service {

	return &Service{
		store: store,

		hooks: hookRegistry,
	}
}

// Store mengembalikan store yang digunakan.
func (s *Service) Store() Store {

	return s.store
}

func (s *Service) fire(
	ctx context.Context,
	event hooks.Event,
	payload hooks.Payload,
) error {

	if s.hooks == nil {
		return nil
	}

	return s.hooks.Fire(
		ctx,
		event,
		payload,
	)
}

// Create membuat session baru tanpa device.
func (s *Service) Create(
	userID string,
) (Session, error) {

	now := time.Now()

	sess := Session{

		ID: uuid.NewString(),

		UserID: userID,

		Revoked: false,

		LastSeenAt: now,

		ExpiresAt: now.Add(
			DefaultSessionTTL,
		),

		CreatedAt: now,
	}

	if err := s.store.Create(
		sess,
	); err != nil {

		return sess, err
	}

	if err := s.fire(
		context.Background(),
		hooks.EventSessionCreated,
		hooks.Payload{

			UserID: sess.UserID,

			SessionID: sess.ID,

			Metadata: map[string]any{
				"type": "session",
			},
		},
	); err != nil {

		return sess, err
	}

	return sess, nil
}

// CreateWithDevice membuat session dengan informasi device.
func (s *Service) CreateWithDevice(
	userID string,
	deviceID string,
	deviceName string,
	platform string,
	ipAddress string,
	userAgent string,
) (Session, error) {

	now := time.Now()

	sess := Session{

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

	if err := s.store.Create(
		sess,
	); err != nil {

		return sess, err
	}

	if err := s.fire(
		context.Background(),
		hooks.EventSessionCreated,
		hooks.Payload{

			UserID: sess.UserID,

			SessionID: sess.ID,

			Metadata: map[string]any{

				"device_id": sess.DeviceID,

				"device": sess.DeviceName,

				"platform": sess.Platform,
			},
		},
	); err != nil {

		return sess, err
	}

	return sess, nil
}

func (s *Service) Get(
	id string,
) (Session, error) {

	return s.store.Get(
		id,
	)
}

func (s *Service) GetByUserID(
	userID string,
) ([]Session, error) {

	return s.store.GetByUserID(
		userID,
	)
}

// Revoke mencabut satu session.
func (s *Service) Revoke(
	id string,
) error {

	if err := s.store.Revoke(
		id,
	); err != nil {

		return err
	}

	return s.fire(
		context.Background(),
		hooks.EventSessionRevoked,
		hooks.Payload{

			SessionID: id,

			Metadata: map[string]any{
				"type": "single",
			},
		},
	)
}

// RevokeAll mencabut seluruh session user.
func (s *Service) RevokeAll(
	userID string,
) error {

	sessions, err :=
		s.store.GetByUserID(
			userID,
		)

	if err != nil {
		return err
	}

	for _, sess := range sessions {

		if sess.Revoked {
			continue
		}

		if err :=
			s.Revoke(
				sess.ID,
			); err != nil {

			return err
		}
	}

	return nil
}

// CheckSession memastikan session masih aktif.
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

// Touch memperbarui aktivitas session.
func (s *Service) Touch(
	id string,
) error {

	sess, err :=
		s.store.Get(
			id,
		)

	if err != nil {
		return err
	}

	sess.LastSeenAt =
		time.Now()

	return s.store.Create(
		sess,
	)
}

// DeleteExpired menghapus session expired.
func (s *Service) DeleteExpired(
	now time.Time,
) error {

	return s.store.DeleteExpired(
		now,
	)
}

// EnforceLimit membatasi jumlah session aktif.
func (s *Service) EnforceLimit(
	userID string,
	maxSessions int,
) error {

	if maxSessions <= 0 {
		return nil
	}

	list, err :=
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

	for _, sess := range list {

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

	oldest :=
		active[0]

	for _, sess := range active {

		if sess.CreatedAt.Before(
			oldest.CreatedAt,
		) {

			oldest = sess
		}
	}

	return s.Revoke(
		oldest.ID,
	)
}
