// Package events contains domain events logic.
package events

import "time"

// DomainEventBuilder a builder to create domain events.
type DomainEventBuilder struct {
	e DomainEvent
}

// NewDomainEventBuilder creates a new domain event builder instance.
func NewDomainEventBuilder(eventType string) *DomainEventBuilder {
	return &DomainEventBuilder{
		e: DomainEvent{
			eventType:  eventType,
			occurredAt: time.Now(),
		},
	}
}

// ForAggregate sets the aggregate ID for event.
func (b *DomainEventBuilder) ForAggregate(id string) *DomainEventBuilder {
	b.e.aggregateID = id
	return b
}

// Build creates a domain event.
func (b *DomainEventBuilder) Build() DomainEvent {
	return b.e
}
