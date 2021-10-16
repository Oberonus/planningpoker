package state

import (
	"context"
	"errors"
	"planningpoker/internal/domain/games"
	"time"
)

type Service struct {
	gamesRepo GameRepository
	usersRepo UsersRepository
}

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

func (s *Service) GameState(cmd GameStateCommand) (*GameState, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cmd.WaitFor)
	defer cancel()

	var game *games.Game
	err := s.gamesRepo.ModifyExclusively(cmd.GameID, func(g *games.Game) error {
		game = g
		return g.Ping(cmd.UserID)
	})
	if err != nil {
		return nil, err
	}

	game, err = s.getUpdatedState(ctx, cmd.GameID, cmd.LastChangeID)
	if err != nil {
		return nil, err
	}

	userIDs := make([]string, 0)
	for id := range game.Players {
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

		if game.ChangeID != lastKnownStateID {
			return game, nil
		}

		select {
		case <-ctx.Done():
			return game, nil
		case <-time.After(100 * time.Millisecond):
		}
	}
}
