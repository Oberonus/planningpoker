package games

import (
	"errors"
	"fmt"

	"planningpoker/internal/domain/events"
)

// Service is the game related application service.
type Service struct {
	gamesRepo GameRepository
}

// NewService creates a new game domain service instance.
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
	eb.Subscribe(gs.processUserUpdated, events.EventTypeUserUpdated)

	return gs, nil
}

// Create creates a game.
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

// Update updates a game.
func (s *Service) Update(cmd UpdateGameCommand) error {
	return s.gamesRepo.ModifyExclusively(cmd.GameID, func(game *Game) error {
		return game.Update(cmd)
	})
}

// Join adds a player to the game.
func (s *Service) Join(cmd JoinGameCommand) error {
	return s.gamesRepo.ModifyExclusively(cmd.GameID, func(game *Game) error {
		return game.Join(cmd)
	})
}

// Restart restarts the game.
func (s *Service) Restart(cmd RestartGameCommand) error {
	return s.gamesRepo.ModifyExclusively(cmd.GameID, func(game *Game) error {
		return game.Restart(cmd)
	})
}

// Vote performs player voting.
func (s *Service) Vote(cmd VoteCommand) error {
	return s.gamesRepo.ModifyExclusively(cmd.GameID, func(game *Game) error {
		return game.Vote(cmd)
	})
}

// UnVote removes a player vote.
func (s *Service) UnVote(cmd UnVoteCommand) error {
	return s.gamesRepo.ModifyExclusively(cmd.GameID, func(game *Game) error {
		return game.UnVote(cmd)
	})
}

// Reveal opens all cards and stops the game.
func (s *Service) Reveal(cmd RevealCardsCommand) error {
	return s.gamesRepo.ModifyExclusively(cmd.GameID, func(game *Game) error {
		return game.Reveal(cmd)
	})
}

// Ping updates player state (assure that the player is still active).
func (s *Service) Ping(cmd PlayerPingCommand) error {
	return s.gamesRepo.ModifyExclusively(cmd.GameID, func(g *Game) error {
		return g.Ping(cmd.UserID)
	})
}

func (s *Service) processUserUpdated(e events.DomainEvent) {
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
			fmt.Printf("no way to update the game, will be updated eventually")
		}
	}
}
