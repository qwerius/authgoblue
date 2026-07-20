package storage

import (
	"context"
)

type Store interface {

	// Save menyimpan data berdasarkan key.
	Set(
		ctx context.Context,
		key string,
		value []byte,
	) error

	// Get mengambil data berdasarkan key.
	Get(
		ctx context.Context,
		key string,
	) ([]byte, error)

	// Delete menghapus data berdasarkan key.
	Delete(
		ctx context.Context,
		key string,
	) error
}
