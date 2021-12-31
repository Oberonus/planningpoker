package events

// Consumer is a generic domain event consumer function.
type Consumer func(DomainEvent)

// EventBus represents a domain events bus generic contract.
type EventBus interface {
	Publish(event DomainEvent) error
	Subscribe(consumer Consumer, eventTypes ...string)
}
