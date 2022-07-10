package events

import "time"

const (
	// EventTypeUserUpdated is a domain event that user information was updated.
	EventTypeUserUpdated = "user:updated"

	// EventTypeGameUpdated is a domain event that game state has changed.
	EventTypeGameUpdated = "game:updated"
)

// DomainEvent is a generic domain event.
type DomainEvent struct {
	eventType   string
	aggregateID string
	occurredAt  time.Time
}

// EventType returns the domain event type.
func (e DomainEvent) EventType() string {
	return e.eventType
}

// AggregateID returns the aggregate ID for the event.
func (e DomainEvent) AggregateID() string {
	return e.aggregateID
}
