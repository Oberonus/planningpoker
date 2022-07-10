package async

import (
	"planningpoker/internal/infra/transformers"

	"github.com/sirupsen/logrus"
	"planningpoker/internal/domain/events"
	"planningpoker/internal/domain/state"
)

// Publisher sends a message to a player.
type Publisher interface {
	SendToPlayer(gameID, userID string, state transformers.GameStateResponse) error
}

// StateProvider retrieves a game state.
type StateProvider interface {
	GameState(gameID string) (*state.GameState, error)
}

// Broadcaster sends a game state to all players.
type Broadcaster struct {
	stateProvider StateProvider
	publisher     Publisher
}

// NewBroadcaster creates a new broadcaster instance.
func NewBroadcaster(stateProvider StateProvider, pub Publisher, eventBus events.EventBus) (*Broadcaster, error) {
	b := &Broadcaster{
		publisher:     pub,
		stateProvider: stateProvider,
	}
	eventBus.Subscribe(b.processGameUpdated, events.EventTypeGameUpdated)
	return b, nil
}

func (b *Broadcaster) processGameUpdated(e events.DomainEvent) {
	gameState, err := b.stateProvider.GameState(e.AggregateID())
	if err != nil {
		logrus.Errorf("failed to fetch game state %v", err)
	}

	for _, playerState := range gameState.Players {
		if err := b.publisher.SendToPlayer(gameState.GameID, playerState.UserID, transformers.NewGameStateResponse(playerState, *gameState)); err != nil {
			logrus.Errorf("failed to send state to the player with ID=%s, %+v, %v", playerState.UserID, gameState, err)
		}
	}
}
