package users

import (
	"errors"
	"fmt"
	"planningpoker/internal/domain/events"
)

// Service is a user related application service.
type Service struct {
	usersRepo Repository
	eventBus  events.EventBus
}

// NewService creates a new users service instance.
func NewService(ur Repository, eb events.EventBus) (*Service, error) {
	if ur == nil {
		return nil, errors.New("users repository should be provided")
	}
	if eb == nil {
		return nil, errors.New("event bus should be provided")
	}

	return &Service{
		usersRepo: ur,
		eventBus:  eb,
	}, nil
}

// Register is a first time registration, without userID known.
func (s *Service) Register(cmd RegisterCommand) (*User, error) {
	u, err := NewUser(cmd.Name)
	if err != nil {
		return nil, fmt.Errorf("user creation: %w", err)
	}

	if err := s.usersRepo.Save(*u); err != nil {
		return nil, err
	}

	return u, nil
}

// Update updates some user details.
func (s *Service) Update(cmd UpdateCommand) (*User, error) {
	u, err := s.usersRepo.Get(cmd.ID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("user not found")
	}

	if err := u.NameAs(cmd.Name); err != nil {
		return nil, err
	}

	if err := s.usersRepo.Save(*u); err != nil {
		return nil, err
	}

	// current implementation provides at-most-once delivery guarantees.
	err = s.eventBus.Publish(events.NewDomainEventBuilder(events.EventTypeUserUpdated).ForAggregate(u.ID()).Build())
	if err != nil {
		fmt.Printf("failed to publish 'user updated' domain event for user id=%s", u.ID())
	}

	return u, nil
}

// AuthenticateByID checks that the user with provided ID exists.
func (s *Service) AuthenticateByID(cmd AuthByIDCommand) (*User, error) {
	u, err := s.usersRepo.Get(cmd.ID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("user not found")
	}

	return u, nil
}
