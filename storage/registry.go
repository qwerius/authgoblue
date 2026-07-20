package storage

import "sync"

type Registry struct {
	mu sync.RWMutex

	stores map[string]Store
}

func NewRegistry() *Registry {

	return &Registry{

		stores: make(
			map[string]Store,
		),
	}
}

func (r *Registry) Register(
	name string,
	store Store,
) {

	r.mu.Lock()

	defer r.mu.Unlock()

	r.stores[name] = store
}

func (r *Registry) Get(
	name string,
) (Store, bool) {

	r.mu.RLock()

	defer r.mu.RUnlock()

	store, ok := r.stores[name]

	return store, ok
}

func (r *Registry) Remove(
	name string,
) {

	r.mu.Lock()

	defer r.mu.Unlock()

	delete(
		r.stores,
		name,
	)
}

func (r *Registry) All() map[string]Store {

	r.mu.RLock()

	defer r.mu.RUnlock()

	result := make(
		map[string]Store,
		len(r.stores),
	)

	for name, store := range r.stores {

		result[name] = store
	}

	return result
}
