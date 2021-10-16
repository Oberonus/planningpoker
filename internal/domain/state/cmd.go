package state

import "time"

type GameStateCommand struct {
	GameID       string
	UserID       string
	WaitFor      time.Duration
	LastChangeID string
}

func NewGameStateCommand(gameID, userID string, waitFor time.Duration, lastChangeID string) (*GameStateCommand, error) {
	return &GameStateCommand{
		GameID:       gameID,
		UserID:       userID,
		WaitFor:      waitFor,
		LastChangeID: lastChangeID,
	}, nil
}
