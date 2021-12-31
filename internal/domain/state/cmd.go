package state

import "time"

// GameStateCommand is a command to get a game state.
type GameStateCommand struct {
	GameID       string
	UserID       string
	WaitFor      time.Duration
	LastChangeID string
}

// NewGameStateCommand creates a new command instance.
func NewGameStateCommand(gameID, userID string, waitFor time.Duration, lastChangeID string) (*GameStateCommand, error) {
	return &GameStateCommand{
		GameID:       gameID,
		UserID:       userID,
		WaitFor:      waitFor,
		LastChangeID: lastChangeID,
	}, nil
}
