// Package http contains http related infra logic.
package http

import (
	"context"
	"errors"
	"net/http"
	"time"

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
	usersService   UsersService
	authenticator  userAuthenticator
	readinessCheck func(context.Context) error
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

// SetReadinessCheck adds a dependency check used by the deployment proxy.
func (h *API) SetReadinessCheck(check func(context.Context) error) {
	h.readinessCheck = check
}

// Alive returns status 200 when storage is reachable, or 503 when it is unavailable.
func (h *API) Alive(ctx *gin.Context) {
	if h.readinessCheck != nil {
		checkCtx, cancel := context.WithTimeout(ctx.Request.Context(), 2*time.Second)
		defer cancel()

		if err := h.readinessCheck(checkCtx); err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "storage unavailable"})
			return
		}
	}

	success(ctx, nil)
}
