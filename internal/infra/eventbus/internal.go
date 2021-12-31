// Package eventbus contains internal domain events bus implementation.
package eventbus

import (
	"planningpoker/internal/domain/events"
	"sync"
)

// InternalBus represents a simple in-service pub/sub service for domain events.
type InternalBus struct {
	m           sync.RWMutex
	subscribers map[string][]events.Consumer
}

// NewInternalBus creates a new internal bus instance.
func NewInternalBus() *InternalBus {
	return &InternalBus{
		subscribers: make(map[string][]events.Consumer),
	}
}

// Publish publishes the provided event.
func (b *InternalBus) Publish(event events.DomainEvent) error {
	b.m.RLock()
	defer b.m.RUnlock()

	for _, c := range b.subscribers[event.EventType()] {
		go c(event)
	}

	return nil
}

// Subscribe is a way for services to subscribe to some events.
func (b *InternalBus) Subscribe(consumer events.Consumer, eventTypes ...string) {
	b.m.Lock()
	defer b.m.Unlock()

	for _, typ := range eventTypes {
		b.subscribers[typ] = append(b.subscribers[typ], consumer)
	}
}
