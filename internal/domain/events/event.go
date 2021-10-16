package events

import "time"

const (
	EventTypeUserUpdated = "user:updated"
)

type DomainEvent struct {
	eventType   string
	aggregateID string
	occurredAt  time.Time
}

func (e DomainEvent) EventType() string {
	return e.eventType
}

func (e DomainEvent) AggregateID() string {
	return e.aggregateID
}
