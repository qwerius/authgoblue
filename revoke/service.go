package revoke

import "time"

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

func (s *Service) Revoke(
	tokenID string,
	expireAt time.Time,
) error {

	return s.store.Revoke(
		tokenID,
		expireAt,
	)
}

func (s *Service) IsRevoked(
	tokenID string,
) (bool, error) {

	return s.store.IsRevoked(
		tokenID,
	)
}
