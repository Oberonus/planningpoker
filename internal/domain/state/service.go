package state

import (
	"context"
	"errors"
	"planningpoker/internal/domain/games"
	"time"
)

// Service is a game state service.
type Service struct {
	gamesRepo GameRepository
	usersRepo UsersRepository
}

// NewService creates a new game state service instance.
func NewService(gr GameRepository, ur UsersRepository) (*Service, error) {
	if gr == nil {
		return nil, errors.New("games repository should be provided")
	}
	if ur == nil {
		return nil, errors.New("users repository should be provided")
	}

	return &Service{
		gamesRepo: gr,
		usersRepo: ur,
	}, nil
}

// GameState returns a current state of a game.
func (s *Service) GameState(cmd GameStateCommand) (*GameState, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cmd.WaitFor)
	defer cancel()

	game, err := s.getUpdatedState(ctx, cmd.GameID, cmd.LastChangeID)
	if err != nil {
		return nil, err
	}

	if !game.IsPlayer(cmd.UserID) {
		return nil, errors.New("user is not a player")
	}

	userIDs := make([]string, 0)
	for id := range game.Players() {
		userIDs = append(userIDs, id)
	}

	users, err := s.usersRepo.GetMany(userIDs)
	if err != nil {
		return nil, err
	}

	state := NewStateForGame(cmd.UserID, *game, users)

	return &state, nil
}

func (s *Service) getUpdatedState(ctx context.Context, gameID, lastKnownStateID string) (*games.Game, error) {
	for {
		game, err := s.gamesRepo.Get(gameID)
		if err != nil {
			return nil, err
		}

		if game.ChangeID() != lastKnownStateID {
			return game, nil
		}

		select {
		case <-ctx.Done():
			return game, nil
		case <-time.After(100 * time.Millisecond):
		}
	}
}