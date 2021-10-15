package users

import (
	"errors"
	"fmt"
)

type Service struct {
	usersRepo Repository
}

func NewService(ur Repository) (*Service, error) {
	if ur == nil {
		return nil, errors.New("users repository should be provided")
	}

	return &Service{
		usersRepo: ur,
	}, nil
}

// Register is a first time registration, without userID known
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

	return u, nil
}

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
