// Package auth provides user authentication functionality.
package auth

import (
	"fmt"

	"planningpoker/internal/domain/users"
)

type usersService interface {
	AuthenticateByID(cmd users.AuthByIDCommand) (*users.User, error)
}

// UserAuthenticator is a service to authenticate a use.
type UserAuthenticator struct {
	usersService
}

// NewUserAuthenticator creates a new user authenticator instance.
func NewUserAuthenticator(us usersService) *UserAuthenticator {
	return &UserAuthenticator{
		usersService: us,
	}
}

// AuthenticateByToken authenticates user by bare token.
func (a *UserAuthenticator) AuthenticateByToken(token string) (string, error) {
	// for simplicity for now token is actual user ID
	cmd, err := users.NewAuthByIDCommand(token)
	if err != nil {
		return "", fmt.Errorf("unable to create an auth command: %w", err)
	}

	user, err := a.usersService.AuthenticateByID(*cmd)
	if err != nil {
		return "", fmt.Errorf("unable to authenticate by ID: %w", err)
	}

	if user == nil {
		return "", fmt.Errorf("user not found")
	}

	return user.ID(), nil
}
