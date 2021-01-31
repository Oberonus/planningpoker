package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"planningpoker/internal/domain"
	"sync"
	"time"
)

type playerDTO struct {
	VotedCard string    `json:"voted_card"`
	CanReveal bool      `json:"can_reveal"`
	LastPing  time.Time `json:"last_ping"`
}

type gameDTO struct {
	ID      string
	Cards   []string             `json:"cards"`
	Players map[string]playerDTO `json:"players"`
	State   string               `json:"state"`
}

type MemoryGameRepository struct {
	gm    sync.RWMutex
	m     sync.RWMutex
	games map[string][]byte
}

func NewMemoryGameRepository() *MemoryGameRepository {
	return &MemoryGameRepository{
		games: make(map[string][]byte),
	}
}

// ModifyExclusively does exclusive blocking modification, so no other goroutines can modify the database exclusively
// quick and dirty implementation, should evolve in something blocking on an external database level...
func (r *MemoryGameRepository) ModifyExclusively(id string, cb func(*domain.Game) error) error {
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

func (r *MemoryGameRepository) Save(game *domain.Game) error {
	dto := gameDTO{
		ID:      game.ID,
		Cards:   game.Cards,
		Players: make(map[string]playerDTO),
		State:   game.State,
	}

	for id, p := range game.Players {
		dto.Players[id] = playerDTO{
			VotedCard: p.VotedCard,
			CanReveal: p.CanReveal,
			LastPing:  p.LastPing,
		}
	}

	raw, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	r.m.Lock()
	defer r.m.Unlock()
	r.games[game.ID] = raw

	return nil
}

func (r *MemoryGameRepository) Get(id string) (*domain.Game, error) {
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

	game := domain.Game{
		ID:      dto.ID,
		Cards:   dto.Cards,
		Players: make(map[string]*domain.Player),
		State:   dto.State,
	}

	for i, p := range dto.Players {
		game.Players[i] = &domain.Player{
			VotedCard: p.VotedCard,
			CanReveal: p.CanReveal,
			LastPing:  p.LastPing,
		}
	}

	return &game, nil
}
