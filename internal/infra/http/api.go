// Package http contains http related infra logic.
package http

import (
	"errors"
	"planningpoker/internal/domain/games"
	"planningpoker/internal/domain/state"
	"planningpoker/internal/domain/users"

	"github.com/gin-gonic/gin"
)

// GamesService is a contract in order to perform actions on games.
type GamesService interface {
	Create(cmd games.CreateGameCommand) (string, error)
	Update(cmd games.UpdateGameCommand) error
	Join(cmd games.JoinGameCommand) error
	Vote(cmd games.VoteCommand) error
	UnVote(cmd games.UnVoteCommand) error
	Reveal(cmd games.RevealCardsCommand) error
	Restart(cmd games.RestartGameCommand) error
	Ping(cmd games.PlayerPingCommand) error
}

// StateService is a contract to get a game state.
type StateService interface {
	GameState(cmd state.GameStateCommand) (*state.GameState, error)
}

// UsersService is a contract to perform user related actions.
type UsersService interface {
	Register(cmd users.RegisterCommand) (*users.User, error)
	Update(cmd users.UpdateCommand) (*users.User, error)
	AuthenticateByID(cmd users.AuthByIDCommand) (*users.User, error)
}

// API contains all HTTP API handlers.
type API struct {
	gamesService GamesService
	usersService UsersService
	stateService StateService
}

// NewAPI creates a new API instance.
func NewAPI(gs GamesService, us UsersService, ss StateService) (*API, error) {
	if gs == nil {
		return nil, errors.New("games service should be provided")
	}
	if us == nil {
		return nil, errors.New("users service should be provided")
	}
	if ss == nil {
		return nil, errors.New("state service should be provided")
	}

	return &API{
		gamesService: gs,
		usersService: us,
		stateService: ss,
	}, nil
}

// SetupRoutes creates HTTP API routes and binds them to handlers.
func (h *API) SetupRoutes(r gin.IRoutes) {
	r.GET("/alive", h.Alive)

	r.POST("/api/v1/register", h.register)

	r.GET("/api/v1/me", h.withUser(h.currentUser))
	r.PUT("/api/v1/me", h.withUser(h.changeUserData))

	r.POST("/api/v1/games", h.withUser(h.createGame))
	r.PUT("/api/v1/games/:gameID", h.withUser(h.updateGame))
	r.POST("/api/v1/games/:gameID/join", h.withUser(h.join))
	r.POST("/api/v1/games/:gameID/ping", h.withUser(h.ping))
	r.POST("/api/v1/games/:gameID/votes/:vote", h.withUser(h.vote))
	r.POST("/api/v1/games/:gameID/unvote", h.withUser(h.unVote))
	r.POST("/api/v1/games/:gameID/restart", h.withUser(h.restartGame))
	r.POST("/api/v1/games/:gameID/reveal", h.withUser(h.revealCards))

	r.GET("/api/v1/games/:gameID", h.withUser(h.gameState))
}

// Alive returns status 200 with empty body.
func (h *API) Alive(ctx *gin.Context) {
	success(ctx, nil)
}
