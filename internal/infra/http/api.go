package http

import (
	"errors"
	"planningpoker/internal/domain/games"
	"planningpoker/internal/domain/state"
	"planningpoker/internal/domain/users"

	"github.com/gin-gonic/gin"
)

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

type StateService interface {
	GameState(cmd state.GameStateCommand) (*state.GameState, error)
}

type UsersService interface {
	Register(cmd users.RegisterCommand) (*users.User, error)
	Update(cmd users.UpdateCommand) (*users.User, error)
	AuthenticateByID(cmd users.AuthByIDCommand) (*users.User, error)
}

type API struct {
	gamesService GamesService
	usersService UsersService
	stateService StateService
}

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

func (h *API) SetupRoutes(r gin.IRoutes) {
	r.POST("/api/v1/register", h.Register)

	r.GET("/api/v1/me", h.withUser(h.CurrentUser))
	r.PUT("/api/v1/me", h.withUser(h.ChangeUserData))

	r.POST("/api/v1/games", h.withUser(h.CreateGame))
	r.PUT("/api/v1/games/:gameID", h.withUser(h.UpdateGame))
	r.POST("/api/v1/games/:gameID/join", h.withUser(h.Join))
	r.POST("/api/v1/games/:gameID/ping", h.withUser(h.Ping))
	r.POST("/api/v1/games/:gameID/votes/:vote", h.withUser(h.Vote))
	r.POST("/api/v1/games/:gameID/unvote", h.withUser(h.UnVote))
	r.POST("/api/v1/games/:gameID/restart", h.withUser(h.RestartGame))
	r.POST("/api/v1/games/:gameID/reveal", h.withUser(h.RevealCards))

	r.GET("/api/v1/games/:gameID", h.withUser(h.GameState))
}
