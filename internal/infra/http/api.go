// Package http contains http related infra logic.
package http

import (
	"errors"

	"github.com/gin-gonic/gin"
	"planningpoker/internal/domain/users"
)

// UserAuthenticator is a contract to authenticate users.
type userAuthenticator interface {
	// AuthenticateByToken returns a user ID or error if user is not authenticated
	AuthenticateByToken(token string) (string, error)
}

// UsersService is a contract to perform user related actions.
type UsersService interface {
	Register(cmd users.RegisterCommand) (*users.User, error)
	Update(cmd users.UpdateCommand) (*users.User, error)
	Get(userID string) (*users.User, error)
}

// API contains all HTTP API handlers.
type API struct {
	usersService  UsersService
	authenticator userAuthenticator
}

// NewAPI creates a new API instance.
func NewAPI(us UsersService, auth userAuthenticator) (*API, error) {
	if us == nil {
		return nil, errors.New("users service should be provided")
	}

	if auth == nil {
		return nil, errors.New("user authenticator should be provided")
	}

	return &API{
		usersService:  us,
		authenticator: auth,
	}, nil
}

// SetupRoutes creates HTTP API routes and binds them to handlers.
func (h *API) SetupRoutes(r gin.IRoutes) {
	r.GET("/alive", h.Alive)

	r.POST("/api/v1/register", h.register)

	r.GET("/api/v1/me", h.withUser(h.currentUser))
	r.PUT("/api/v1/me", h.withUser(h.changeUserData))
}

// Alive returns status 200 with empty body.
func (h *API) Alive(ctx *gin.Context) {
	success(ctx, nil)
}
