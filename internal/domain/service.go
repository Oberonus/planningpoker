package domain

import (
	"context"
	"errors"
	"time"
)

type GamesService struct {
	gamesRepo GameRepository
	usersRepo UsersRepository
}

func NewGamesService(gr GameRepository, ur UsersRepository) (*GamesService, error) {
	if gr == nil {
		return nil, errors.New("games repository should be provided")
	}
	if ur == nil {
		return nil, errors.New("users repository should be provided")
	}

	return &GamesService{
		gamesRepo: gr,
		usersRepo: ur,
	}, nil
}

func (s *GamesService) Create(cmd CreateGameCommand) (string, error) {
	game := NewGame(cmd)

	joinCmd, err := NewJoinGameCommand(game.ID, cmd.UserID)
	if err != nil {
		return "", err
	}

	if err := game.Join(*joinCmd); err != nil {
		return "", err
	}
	if err := s.gamesRepo.Save(game); err != nil {
		return "", err
	}

	return game.ID, nil
}

func (s *GamesService) Update(cmd UpdateGameCommand) error {
	return s.gamesRepo.ModifyExclusively(cmd.GameID, func(game *Game) error {
		return game.Update(cmd)
	})
}

func (s *GamesService) Join(cmd JoinGameCommand) error {
	return s.gamesRepo.ModifyExclusively(cmd.GameID, func(game *Game) error {
		return game.Join(cmd)
	})
}

func (s *GamesService) Restart(cmd RestartGameCommand) error {
	return s.gamesRepo.ModifyExclusively(cmd.GameID, func(game *Game) error {
		return game.Restart(cmd)
	})
}

func (s *GamesService) Vote(cmd VoteCommand) error {
	return s.gamesRepo.ModifyExclusively(cmd.GameID, func(game *Game) error {
		return game.Vote(cmd)
	})
}

func (s *GamesService) UnVote(cmd UnVoteCommand) error {
	return s.gamesRepo.ModifyExclusively(cmd.GameID, func(game *Game) error {
		return game.UnVote(cmd)
	})
}

func (s *GamesService) Reveal(cmd RevealCardsCommand) error {
	return s.gamesRepo.ModifyExclusively(cmd.GameID, func(game *Game) error {
		return game.Reveal(cmd)
	})
}

func (s *GamesService) GameState(cmd GameStateCommand) (*GameState, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cmd.WaitFor)
	defer cancel()

	var game *Game
	err := s.gamesRepo.ModifyExclusively(cmd.GameID, func(g *Game) error {
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

func (s *GamesService) getUpdatedState(ctx context.Context, gameID, lastKnownStateID string) (*Game, error) {
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
