package state

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"planningpoker/internal/domain/events"
)

// Publisher sends a message to a player.
type Publisher interface {
	SendToPlayer(gameState GameState, userID string) error
}

// Service is a game state service.
type Service struct {
	gamesRepo GameRepository
	usersRepo UsersRepository
	publisher Publisher
}

// NewService creates a new game state service instance.
func NewService(gr GameRepository, ur UsersRepository, pub Publisher, eventBus events.EventBus) (*Service, error) {
	if gr == nil {
		return nil, errors.New("games repository should be provided")
	}
	if ur == nil {
		return nil, errors.New("users repository should be provided")
	}
	if eventBus == nil {
		return nil, errors.New("event bus should be provided")
	}
	if pub == nil {
		return nil, errors.New("publisher should be provided")
	}

	srv := &Service{
		gamesRepo: gr,
		usersRepo: ur,
		publisher: pub,
	}

	eventBus.Subscribe(srv.processGameUpdated, events.EventTypeGameUpdated)

	return srv, nil
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

func (s *Service) processGameUpdated(e events.DomainEvent) {
	gameState, err := s.GameState(e.AggregateID())
	if err != nil {
		logrus.Errorf("failed to fetch game state %v", err)
	}

	for _, playerState := range gameState.Players {
		if err := s.publisher.SendToPlayer(*gameState, playerState.UserID); err != nil {
			logrus.Errorf("failed to send state to the player with ID=%s, %+v, %v", playerState.UserID, gameState, err)
		}
	}
}
