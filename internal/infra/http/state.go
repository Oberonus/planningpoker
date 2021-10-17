package http

import (
	"planningpoker/internal/domain/games"
	"planningpoker/internal/domain/state"
	"planningpoker/internal/domain/users"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *API) GameState(c *gin.Context, user *users.User) {
	gameID := c.Param("gameID")
	lastChangeID := c.Query("lastChangeID")

	cmd, err := state.NewGameStateCommand(gameID, user.ID(), 5*time.Second, lastChangeID)
	if err != nil {
		badRequestError(c, err)
		return
	}

	st, err := h.stateService.GameState(*cmd)
	if err != nil {
		internalError(c, err)
		return
	}

	success(c, newGameStateResponse(*st))
}

type playerStateResponse struct {
	Name      string `json:"name"`
	VotedCard string `json:"voted_card"`
}

func newPlayerStateResponse(state state.PlayerState) playerStateResponse {
	resp := playerStateResponse{
		Name: state.Name,
	}
	if state.VotedCard != nil {
		resp.VotedCard = state.VotedCard.Type()
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

type gameStateResponse struct {
	Name      string                `json:"name"`
	TicketURL string                `json:"ticket_url"`
	ChangeID  string                `json:"change_id"`
	CardsDeck cardsDeckResponse     `json:"cards_deck"`
	Players   []playerStateResponse `json:"players"`
	State     string                `json:"state"`
	VotedCard string                `json:"voted_card"`
	CanReveal bool                  `json:"can_reveal"`
}

func newGameStateResponse(state state.GameState) gameStateResponse {
	resp := gameStateResponse{
		Name:      state.Name,
		TicketURL: state.TicketURL,
		ChangeID:  state.ChangeID,
		CardsDeck: newCardsDeckResponse(state.CardsDeck),
		State:     state.State,
		CanReveal: state.CanReveal,
	}
	if state.VotedCard != nil {
		resp.VotedCard = state.VotedCard.Type()
	}
	for _, p := range state.Players {
		resp.Players = append(resp.Players, newPlayerStateResponse(p))
	}
	return resp
}
