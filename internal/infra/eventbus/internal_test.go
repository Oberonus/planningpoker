package eventbus_test

import (
	"planningpoker/internal/domain/events"
	"planningpoker/internal/infra/eventbus"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInternalBus(t *testing.T) {
	bus := eventbus.NewInternalBus()

	calledCh := make(chan struct{})
	eventType := "event1"

	consumer := func(event events.DomainEvent) {
		assert.Equal(t, eventType, event.EventType())
		calledCh <- struct{}{}
	}

	bus.Subscribe(consumer, "event1", "event2")

	err := bus.Publish(events.NewDomainEventBuilder(eventType).Build())
	require.NoError(t, err)

	select {
	case <-calledCh:
		break
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("consumer func was not called")
	}
}
