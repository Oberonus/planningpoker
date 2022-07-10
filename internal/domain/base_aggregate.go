package domain

import "planningpoker/internal/domain/events"

// BaseAggregate is a base struct for all aggregates.
type BaseAggregate struct {
	events []events.DomainEvent
}

// GetEvents returns all created domain events.
func (a *BaseAggregate) GetEvents() []events.DomainEvent {
	return a.events
}

// AddEvent adds one domain event to the list.
func (a *BaseAggregate) AddEvent(e events.DomainEvent) {
	a.events = append(a.events, e)
}

// ClearEvents deletes all domain events.
func (a *BaseAggregate) ClearEvents() {
	a.events = nil
}
