package eventbus

import (
	"planningpoker/internal/domain/events"
	"sync"
)

type InternalBus struct {
	m           sync.RWMutex
	subscribers map[string][]events.Consumer
}

func NewInternalBus() *InternalBus {
	return &InternalBus{
		subscribers: make(map[string][]events.Consumer),
	}
}

func (b *InternalBus) Publish(event events.DomainEvent) error {
	b.m.RLock()
	defer b.m.RUnlock()

	for _, c := range b.subscribers[event.EventType()] {
		go c(event)
	}

	return nil
}

func (b *InternalBus) Subscribe(consumer events.Consumer, eventTypes ...string) {
	b.m.Lock()
	defer b.m.Unlock()

	for _, typ := range eventTypes {
		b.subscribers[typ] = append(b.subscribers[typ], consumer)
	}
}
