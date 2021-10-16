package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"planningpoker/internal/domain"
	"sync"
	"time"
)

type cardsDeckDTO struct {
	Name  string   `json:"name"`
	Cards []string `json:"cards"`
}

func newCardsDeckDTO(d domain.CardsDeck) cardsDeckDTO {
	dto := cardsDeckDTO{
		Name:  d.Name(),
		Cards: make([]string, len(d.Cards())),
	}

	for i, c := range d.Cards() {
		dto.Cards[i] = c.Type()
	}

	return dto
}

func (d cardsDeckDTO) ToDomain() (*domain.CardsDeck, error) {
	cards := make([]domain.Card, len(d.Cards))
	for i, v := range d.Cards {
		c, err := domain.NewCard(v)
		if err != nil {
			return nil, fmt.Errorf("card creation: %w", err)
		}
		cards[i] = *c
	}

	return domain.NewCardsDeck(d.Name, cards)
}

type playerDTO struct {
	VotedCard string    `json:"voted_card"`
	CanReveal bool      `json:"can_reveal"`
	LastPing  time.Time `json:"last_ping"`
}

func (d playerDTO) ToDomain() (*domain.Player, error) {
	var votedCard *domain.Card
	var err error

	if d.VotedCard != "" {
		votedCard, err = domain.NewCard(d.VotedCard)
		if err != nil {
			return nil, err
		}
	}

	return &domain.Player{
		VotedCard: votedCard,
		CanReveal: d.CanReveal,
		LastPing:  d.LastPing,
	}, nil
}

type gameDTO struct {
	ID                string               `json:"id"`
	Name              string               `json:"name"`
	TicketURL         string               `json:"ticket_url"`
	CardsDeck         cardsDeckDTO         `json:"cards_deck"`
	Players           map[string]playerDTO `json:"players"`
	State             string               `json:"state"`
	ChangeID          string               `json:"change_id"`
	EveryoneCanReveal bool                 `json:"everyone_can_reveal"`
}

func (d gameDTO) ToDomain() (*domain.Game, error) {
	game := &domain.Game{
		ID:                d.ID,
		Name:              d.Name,
		TicketURL:         d.TicketURL,
		Players:           make(map[string]*domain.Player),
		State:             d.State,
		ChangeID:          d.ChangeID,
		EveryoneCanReveal: d.EveryoneCanReveal,
	}

	deck, err := d.CardsDeck.ToDomain()
	if err != nil {
		return nil, err
	}
	game.CardsDeck = *deck

	for i, p := range d.Players {
		game.Players[i], err = p.ToDomain()
		if err != nil {
			return nil, err
		}
	}

	return game, err
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
		ID:                game.ID,
		Name:              game.Name,
		TicketURL:         game.TicketURL,
		CardsDeck:         newCardsDeckDTO(game.CardsDeck),
		Players:           make(map[string]playerDTO),
		State:             game.State,
		ChangeID:          game.ChangeID,
		EveryoneCanReveal: game.EveryoneCanReveal,
	}

	for id, p := range game.Players {
		votedCard := ""
		if p.VotedCard != nil {
			votedCard = p.VotedCard.Type()
		}
		dto.Players[id] = playerDTO{
			VotedCard: votedCard,
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

	return dto.ToDomain()
}

func (r *MemoryGameRepository) GetActiveGamesByPlayerID(playerID string) ([]domain.Game, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	list := make([]domain.Game, 0)
	for _, j := range r.games {
		dto := gameDTO{}
		err := json.Unmarshal(j, &dto)
		if err != nil {
			return nil, err
		}
		if dto.State != domain.GameStateStarted {
			continue
		}
		if _, ok := dto.Players[playerID]; !ok {
			continue
		}

		g, err := dto.ToDomain()
		if err != nil {
			return nil, err
		}

		list = append(list, *g)
	}

	return list, nil
}
