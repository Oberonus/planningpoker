package eventbus_test

import (
	"planningpoker/internal/domain/events"
	"planningpoker/internal/infra/eventbus"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInternalBus_Subscribe(t *testing.T) {
	consumer := func(event events.DomainEvent) {
		t.Logf("BLA!!!!!!!!!!!!!!!!!!!!!!!!")
	}

	bus := eventbus.NewInternalBus()
	bus.Subscribe(consumer, "event1", "event2")

	err := bus.Publish(events.NewDomainEventBuilder("event1").Build())
	require.NoError(t, err)

	time.Sleep(50 * time.Millisecond)
}
