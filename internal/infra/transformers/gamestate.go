// Package transformers contains all transformers from/to domain models.
package transformers

import (
	"planningpoker/internal/domain/games"
	"planningpoker/internal/domain/state"
)

// PlayerStateResponse is a response payload for a player.
type PlayerStateResponse struct {
	Name      string `json:"name"`
	VotedCard string `json:"voted_card"`
}

func newPlayerStateResponse(gState state.GameState, pState state.PlayerState) PlayerStateResponse {
	resp := PlayerStateResponse{
		Name: pState.Name,
	}
	if pState.VotedCard != nil {
		if gState.State != games.GameStateFinished {
			resp.VotedCard = games.NewUnrevealedCard().Type()
		} else {
			resp.VotedCard = pState.VotedCard.Type()
		}
	}
	return resp
}

type cardsDeckResponse struct {
	Name  string   `json:"name"`
	Cards []string `json:"cards"`
}

func newCardsDeckResponse(cd games.CardsDeck) cardsDeckResponse {
	resp := cardsDeckResponse{Name: cd.Name()}
	for _, c := range cd.Cards() {
		resp.Cards = append(resp.Cards, c.Type())
	}
	return resp
}

// GameStateResponse is a response payload with game state.
type GameStateResponse struct {
	Name      string                `json:"name"`
	TicketURL string                `json:"ticket_url"`
	CardsDeck cardsDeckResponse     `json:"cards_deck"`
	Players   []PlayerStateResponse `json:"players"`
	State     string                `json:"state"`
	VotedCard string                `json:"voted_card"`
	CanReveal bool                  `json:"can_reveal"`
}

// NewGameStateResponse creates a new game state response.
func NewGameStateResponse(player state.PlayerState, state state.GameState) GameStateResponse {
	resp := GameStateResponse{
		Name:      state.Name,
		TicketURL: state.TicketURL,
		CardsDeck: newCardsDeckResponse(state.CardsDeck),
		State:     state.State,
		CanReveal: player.CanReveal,
	}
	if player.VotedCard != nil {
		resp.VotedCard = player.VotedCard.Type()
	}
	for _, p := range state.Players {
		resp.Players = append(resp.Players, newPlayerStateResponse(state, p))
	}
	return resp
}
