package domain

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type User struct {
	ID   string
	Name string
}

type UsersService struct {
	usersRepo UsersRepository
}

func NewUsersService(ur UsersRepository) (*UsersService, error) {
	if ur == nil {
		return nil, errors.New("users repository should be provided")
	}

	return &UsersService{
		usersRepo: ur,
	}, nil
}

// Register is a first time registration, without userID known
func (s *UsersService) Register(name string) (*User, error) {
	if name == "" {
		return nil, errors.New("name can not be empty")
	}

	u := &User{
		ID:   strings.Replace(uuid.New().String(), "-", "", -1),
		Name: name,
	}

	if err := s.usersRepo.Save(u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *UsersService) Update(id, name string) (*User, error) {
	if name == "" {
		return nil, errors.New("name can not be empty")
	}

	u, err := s.usersRepo.Get(id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("user not found")
	}

	u.Name = name

	if err := s.usersRepo.Save(u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *UsersService) AuthenticateByID(id string) (*User, error) {
	u, err := s.usersRepo.Get(id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("user not found")
	}

	return u, nil
}
