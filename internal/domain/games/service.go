package games

import (
	"errors"
	"fmt"
	"planningpoker/internal/domain/events"
)

type Service struct {
	gamesRepo GameRepository
}

func NewService(gr GameRepository, eb events.EventBus) (*Service, error) {
	if gr == nil {
		return nil, errors.New("games repository should be provided")
	}
	if eb == nil {
		return nil, errors.New("event bus should be provided")
	}

	gs := &Service{
		gamesRepo: gr,
	}
	eb.Subscribe(gs.ProcessUserUpdated, events.EventTypeUserUpdated)

	return gs, nil
}

func (s *Service) Create(cmd CreateGameCommand) (string, error) {
	game := NewGame(cmd)

	joinCmd, err := NewJoinGameCommand(game.id, cmd.UserID)
	if err != nil {
		return "", err
	}

	if err := game.Join(*joinCmd); err != nil {
		return "", err
	}
	if err := s.gamesRepo.Save(game); err != nil {
		return "", err
	}

	return game.id, nil
}

func (s *Service) Update(cmd UpdateGameCommand) error {
	return s.gamesRepo.ModifyExclusively(cmd.GameID, func(game *Game) error {
		return game.Update(cmd)
	})
}

func (s *Service) Join(cmd JoinGameCommand) error {
	return s.gamesRepo.ModifyExclusively(cmd.GameID, func(game *Game) error {
		return game.Join(cmd)
	})
}

func (s *Service) Restart(cmd RestartGameCommand) error {
	return s.gamesRepo.ModifyExclusively(cmd.GameID, func(game *Game) error {
		return game.Restart(cmd)
	})
}

func (s *Service) Vote(cmd VoteCommand) error {
	return s.gamesRepo.ModifyExclusively(cmd.GameID, func(game *Game) error {
		return game.Vote(cmd)
	})
}

func (s *Service) UnVote(cmd UnVoteCommand) error {
	return s.gamesRepo.ModifyExclusively(cmd.GameID, func(game *Game) error {
		return game.UnVote(cmd)
	})
}

func (s *Service) Reveal(cmd RevealCardsCommand) error {
	return s.gamesRepo.ModifyExclusively(cmd.GameID, func(game *Game) error {
		return game.Reveal(cmd)
	})
}

func (s *Service) Ping(cmd PlayerPingCommand) error {
	return s.gamesRepo.ModifyExclusively(cmd.GameID, func(g *Game) error {
		return g.Ping(cmd.UserID)
	})
}

func (s *Service) ProcessUserUpdated(e events.DomainEvent) {
	list, err := s.gamesRepo.GetActiveGamesByPlayerID(e.AggregateID())
	if err != nil {
		fmt.Printf("error fetching games for player id=%s: %v", e.AggregateID(), err)
	}

	for _, g := range list {
		err := s.gamesRepo.ModifyExclusively(g.id, func(g *Game) error {
			g.ForceChanged()
			return nil
		})
		if err != nil {
			fmt.Printf("no way to update game, will be updated eventually")
		}
	}
}
