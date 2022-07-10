package state

import (
	"errors"
	"fmt"
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
func (s *Service) GameState(gameID string) (*GameState, error) {
	game, err := s.gamesRepo.Get(gameID)
	if err != nil || game == nil {
		return nil, fmt.Errorf("get game: %w", err)
	}

	userIDs := make([]string, 0)
	for id := range game.Players() {
		userIDs = append(userIDs, id)
	}

	users, err := s.usersRepo.GetMany(userIDs)
	if err != nil {
		return nil, err
	}

	state := NewStateForGame(*game, users)

	return &state, nil
}
