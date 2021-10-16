package state

import (
	"planningpoker/internal/domain/games"
	"planningpoker/internal/domain/users"
)

type PlayerState struct {
	Name      string
	VotedCard *games.Card
}

type GameState struct {
	Name      string
	TicketURL string
	ChangeID  string
	CardsDeck games.CardsDeck
	Players   []PlayerState
	State     string
	VotedCard *games.Card
	CanReveal bool
}

func NewStateForGame(userID string, game games.Game, gamers []users.User) GameState {
	state := GameState{
		CardsDeck: game.CardsDeck,
		Name:      game.Name,
		TicketURL: game.TicketURL,
		Players:   make([]PlayerState, 0, len(game.Players)),
		State:     game.State,
		VotedCard: game.Players[userID].VotedCard,
		CanReveal: game.Players[userID].CanReveal,
		ChangeID:  game.ChangeID,
	}

	for uid, p := range game.Players {
		userName := "Unknown"
		if u := findUserInListByID(uid, gamers); u != nil {
			userName = u.Name()
		}

		votedCard := p.VotedCard
		// mask real votes if game is running
		if game.State == games.GameStateStarted && votedCard != nil {
			unrevealedCard := games.NewUnrevealedCard()
			votedCard = &unrevealedCard
		}

		state.Players = append(state.Players, PlayerState{
			Name:      userName,
			VotedCard: votedCard,
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
