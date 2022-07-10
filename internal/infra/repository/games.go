package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"planningpoker/internal/domain/events"
	"planningpoker/internal/domain/games"
)

// MemoryGameRepository is a simple in-memory linear games repository.
type MemoryGameRepository struct {
	gm       sync.RWMutex
	m        sync.RWMutex
	games    map[string][]byte
	eventBus events.EventBus
}

// NewMemoryGameRepository creates a new in-memory repository instance.
func NewMemoryGameRepository(bus events.EventBus) *MemoryGameRepository {
	return &MemoryGameRepository{
		games:    make(map[string][]byte),
		eventBus: bus,
	}
}

// ModifyExclusively does exclusive blocking modification, so no other goroutines can modify the database exclusively
// quick and dirty implementation, should evolve in something blocking on an external database level...
func (r *MemoryGameRepository) ModifyExclusively(id string, cb func(*games.Game) error) error {
	r.gm.Lock()
	defer r.gm.Unlock()

	game, err := r.Get(id)
	if err != nil {
		return fmt.Errorf("game fetching: %w", err)
	}
	if game == nil {
		return errors.New("game not found")
	}

	if err := cb(game); err != nil {
		return err
	}

	if err := r.Save(game); err != nil {
		return fmt.Errorf("game save: %w", err)
	}

	return nil
}

// Save persists the game.
func (r *MemoryGameRepository) Save(game *games.Game) error {
	dto := gameDTO{
		ID:                game.ID(),
		Name:              game.Name(),
		TicketURL:         game.TicketURL(),
		CardsDeck:         newCardsDeckDTO(game.CardsDeck()),
		Players:           make(map[string]playerDTO),
		State:             game.State(),
		EveryoneCanReveal: game.EveryoneCanReveal(),
	}

	for id, p := range game.Players() {
		votedCard := ""
		if p.VotedCard != nil {
			votedCard = p.VotedCard.Type()
		}
		dto.Players[id] = playerDTO{
			VotedCard: votedCard,
			CanReveal: p.CanReveal,
			Active:    p.Active,
		}
	}

	raw, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	r.m.Lock()
	defer r.m.Unlock()
	r.games[game.ID()] = raw

	for _, e := range game.GetEvents() {
		if err := r.eventBus.Publish(e); err != nil {
			// this is a simple handler, which provides at most once delivery
			// in case if a bus is broken, it will log the error and continue
			// TODO: implement outbox pattern in order to mitigate potential distributed transactions
			logrus.Errorf("failed to publish a game event: %v", err)
		}
	}

	return nil
}

// Get retrieves the game.
func (r *MemoryGameRepository) Get(id string) (*games.Game, error) {
	r.m.RLock()
	raw, ok := r.games[id]
	r.m.RUnlock()

	if !ok {
		return nil, nil
	}

	dto := gameDTO{}
	err := json.Unmarshal(raw, &dto)
	if err != nil {
		return nil, err
	}

	return dto.toDomain()
}

// GetActiveGamesByPlayerID returns all games where specific user is an active participant.
func (r *MemoryGameRepository) GetActiveGamesByPlayerID(playerID string) ([]games.Game, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	list := make([]games.Game, 0)
	for _, j := range r.games {
		dto := gameDTO{}
		err := json.Unmarshal(j, &dto)
		if err != nil {
			return nil, err
		}
		if dto.State != games.GameStateStarted {
			continue
		}
		if _, ok := dto.Players[playerID]; !ok {
			continue
		}

		g, err := dto.toDomain()
		if err != nil {
			return nil, err
		}

		list = append(list, *g)
	}

	return list, nil
}
