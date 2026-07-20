package providers

import (
	"errors"
	"sync"
)

var ErrProviderNotFound = errors.New(
	"authgoblue: provider not found",
)

type Registry struct {
	mu sync.RWMutex

	providers map[string]Provider
}

func NewRegistry() *Registry {

	return &Registry{

		providers: make(map[string]Provider),
	}
}

func (r *Registry) Register(
	name string,
	provider Provider,
) {

	r.mu.Lock()
	defer r.mu.Unlock()

	r.providers[name] = provider
}

func (r *Registry) Get(
	name string,
) (Provider, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, ok :=
		r.providers[name]

	if !ok {

		return nil, ErrProviderNotFound
	}

	return provider, nil
}
