package events

import "time"

type DomainEventBuilder struct {
	e DomainEvent
}

func NewDomainEventBuilder(eventType string) *DomainEventBuilder {
	return &DomainEventBuilder{
		e: DomainEvent{
			eventType:  eventType,
			occurredAt: time.Now(),
		},
	}
}

func (b *DomainEventBuilder) ForAggregate(id string) *DomainEventBuilder {
	b.e.aggregateID = id
	return b
}

func (b *DomainEventBuilder) Build() DomainEvent {
	return b.e
}
