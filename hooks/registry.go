package hooks

import (
	"context"
	"sync"
)

type Registry struct {
	mu sync.RWMutex

	handlers map[Event][]Handler
}

func NewRegistry() *Registry {

	return &Registry{

		handlers: make(map[Event][]Handler),
	}
}

func (r *Registry) Register(
	event Event,
	handler Handler,
) {

	r.mu.Lock()
	defer r.mu.Unlock()

	r.handlers[event] =
		append(
			r.handlers[event],
			handler,
		)
}

func (r *Registry) Fire(
	ctx context.Context,
	event Event,
	payload Payload,
) error {

	r.mu.RLock()

	listeners :=
		append(
			[]Handler(nil),
			r.handlers[event]...,
		)

	r.mu.RUnlock()

	for _, handler := range listeners {

		if err :=
			handler(
				ctx,
				payload,
			); err != nil {

			return err
		}
	}

	return nil
}
