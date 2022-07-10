package repository

import (
	"fmt"
	"planningpoker/internal/domain/games"
)

type cardsDeckDTO struct {
	Name  string   `json:"name"`
	Cards []string `json:"cards"`
}

func newCardsDeckDTO(d games.CardsDeck) cardsDeckDTO {
	dto := cardsDeckDTO{
		Name:  d.Name(),
		Cards: make([]string, len(d.Cards())),
	}

	for i, c := range d.Cards() {
		dto.Cards[i] = c.Type()
	}

	return dto
}

func (d cardsDeckDTO) toDomain() (*games.CardsDeck, error) {
	cards := make([]games.Card, len(d.Cards))
	for i, v := range d.Cards {
		c, err := games.NewCard(v)
		if err != nil {
			return nil, fmt.Errorf("card creation: %w", err)
		}
		cards[i] = *c
	}

	return games.NewCardsDeck(d.Name, cards)
}

type playerDTO struct {
	VotedCard string `json:"voted_card"`
	CanReveal bool   `json:"can_reveal"`
	Active    bool   `json:"active"`
}

func (d playerDTO) toDomain() (*games.Player, error) {
	var votedCard *games.Card
	var err error

	if d.VotedCard != "" {
		votedCard, err = games.NewCard(d.VotedCard)
		if err != nil {
			return nil, err
		}
	}

	return &games.Player{
		VotedCard: votedCard,
		CanReveal: d.CanReveal,
		Active:    d.Active,
	}, nil
}

type gameDTO struct {
	ID                string               `json:"id"`
	Name              string               `json:"name"`
	TicketURL         string               `json:"ticket_url"`
	CardsDeck         cardsDeckDTO         `json:"cards_deck"`
	Players           map[string]playerDTO `json:"players"`
	State             string               `json:"state"`
	EveryoneCanReveal bool                 `json:"everyone_can_reveal"`
}

func (d gameDTO) toDomain() (*games.Game, error) {
	deck, err := d.CardsDeck.toDomain()
	if err != nil {
		return nil, err
	}

	players := make(map[string]*games.Player)
	for i, p := range d.Players {
		players[i], err = p.toDomain()
		if err != nil {
			return nil, err
		}
	}

	game := games.NewRaw(d.ID, d.Name, d.TicketURL, *deck, players, d.State, d.EveryoneCanReveal)

	return game, err
}
