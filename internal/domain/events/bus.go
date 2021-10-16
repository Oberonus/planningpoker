package events

type Consumer func(DomainEvent)

type EventBus interface {
	Publish(event DomainEvent) error
	Subscribe(consumer Consumer, eventTypes ...string)
}
