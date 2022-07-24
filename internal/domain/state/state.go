// Package state contains domain level game state logic.
package state

import (
	"fmt"

	"planningpoker/internal/domain/games"
	"planningpoker/internal/domain/users"
)

// PlayerState represents a player state.
type PlayerState struct {
	UserID     string
	Name       string
	VotedCard  *games.Card
	Confidence string
	CanReveal  bool
}

// GameState represents a game state.
type GameState struct {
	GameID    string
	Name      string
	TicketURL string
	CardsDeck games.CardsDeck
	Players   []PlayerState
	State     string
}

// NewStateForGame creates a new game state.
func NewStateForGame(game games.Game, gamers []users.User) GameState {
	state := GameState{
		GameID:    game.ID(),
		CardsDeck: game.CardsDeck(),
		Name:      game.Name(),
		TicketURL: game.TicketURL(),
		Players:   make([]PlayerState, 0, len(game.Players())),
		State:     game.State(),
	}

	for uid, p := range game.Players() {
		userName := "Unknown"
		if u := findUserInListByID(uid, gamers); u != nil {
			userName = u.Name()
		}

		state.Players = append(state.Players, PlayerState{
			UserID:     uid,
			Name:       userName,
			VotedCard:  p.VotedCard,
			Confidence: p.Confidence,
			CanReveal:  p.CanReveal,
		})
	}

	return state
}

func findUserInListByID(id string, users []users.User) *users.User {
	for _, u := range users {
		if u.ID() == id {
			return &u
		}
	}
	return nil
}

// PlayerByID returns a player state by user ID.
func (s GameState) PlayerByID(userID string) (*PlayerState, error) {
	for _, player := range s.Players {
		if player.UserID == userID {
			return &player, nil
		}
	}
	return nil, fmt.Errorf("player with id=%v not found", userID)
}
