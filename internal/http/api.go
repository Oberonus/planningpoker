package http

import (
	"errors"
	"strings"
	"time"

	"planningpoker/internal/domain"

	"github.com/gin-gonic/gin"
)

type API struct {
	gamesService *domain.GamesService
	usersService *domain.UsersService
}

func NewAPI(gs *domain.GamesService, us *domain.UsersService) (*API, error) {
	if gs == nil {
		return nil, errors.New("games service should be provided")
	}
	if us == nil {
		return nil, errors.New("users service should be provided")
	}

	return &API{
		gamesService: gs,
		usersService: us,
	}, nil
}

func (h *API) SetupRoutes(r gin.IRoutes) {
	r.GET("/api/v1/me", h.CurrentUser)
	r.PUT("/api/v1/me", h.ChangeUserData)
	r.POST("/api/v1/register", h.Register)
	r.POST("/api/v1/games", h.CreateGame)
	r.PUT("/api/v1/games/:gameID", h.UpdateGame)
	r.POST("/api/v1/games/:gameID/join", h.Join)
	r.POST("/api/v1/games/:gameID/votes/:vote", h.Vote)
	r.POST("/api/v1/games/:gameID/unvote", h.UnVote)
	r.POST("/api/v1/games/:gameID/restart", h.RestartGame)
	r.POST("/api/v1/games/:gameID/reveal", h.RevealCards)
	r.GET("/api/v1/games/:gameID", h.GameState)
}

func (h *API) Register(c *gin.Context) {
	pl := struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	}{}
	if err := c.BindJSON(&pl); err != nil {
		badRequestError(c, err)
		return
	}

	user, err := h.usersService.Register(pl.Name)
	if err != nil {
		badRequestError(c, err)
		return
	}

	success(c, gin.H{
		"user_id": user.ID,
	})
}

func (h *API) CurrentUser(c *gin.Context) {
	user := h.authenticate(c)
	if user == nil {
		return
	}

	success(c, gin.H{
		"name": user.Name,
	})
}

func (h *API) ChangeUserData(c *gin.Context) {
	user := h.authenticate(c)
	if user == nil {
		return
	}

	pl := struct {
		Name string `json:"name"`
	}{}
	if err := c.BindJSON(&pl); err != nil {
		badRequestError(c, err)
		return
	}

	_, err := h.usersService.Update(user.ID, pl.Name)
	if err != nil {
		internalError(c, err)
		return
	}

	success(c, gin.H{})
}

type gamePayload struct {
	Name      string `json:"name"`
	TicketURL string `json:"url"`
}

func (h *API) CreateGame(c *gin.Context) {
	user := h.authenticate(c)
	if user == nil {
		return
	}

	pl := gamePayload{}
	if err := c.BindJSON(&pl); err != nil {
		badRequestError(c, err)
		return
	}

	cmd, err := domain.NewCreateGameCommand(pl.Name, pl.TicketURL, user.ID)
	if err != nil {
		badRequestError(c, err)
		return
	}

	gameID, err := h.gamesService.Create(*cmd)
	if err != nil {
		internalError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"game_id": gameID,
	})
}

func (h *API) UpdateGame(c *gin.Context) {
	user := h.authenticate(c)
	if user == nil {
		return
	}
	gameID := c.Param("gameID")

	pl := gamePayload{}
	if err := c.BindJSON(&pl); err != nil {
		badRequestError(c, err)
		return
	}

	cmd, err := domain.NewUpdateGameCommand(gameID, pl.Name, pl.TicketURL, user.ID)
	if err != nil {
		badRequestError(c, err)
		return
	}

	err = h.gamesService.Update(*cmd)
	if err != nil {
		internalError(c, err)
		return
	}

	c.JSON(200, gin.H{})
}

func (h *API) Join(c *gin.Context) {
	user := h.authenticate(c)
	if user == nil {
		return
	}
	gameID := c.Param("gameID")

	err := h.gamesService.Join(gameID, user.ID)
	if err != nil {
		internalError(c, err)
		return
	}

	success(c, nil)
}

func (h *API) Vote(c *gin.Context) {
	user := h.authenticate(c)
	if user == nil {
		return
	}
	gameID := c.Param("gameID")
	vote := c.Param("vote")

	err := h.gamesService.Vote(gameID, user.ID, vote)
	if err != nil {
		internalError(c, err)
		return
	}

	success(c, nil)
}

func (h *API) UnVote(c *gin.Context) {
	user := h.authenticate(c)
	if user == nil {
		return
	}
	gameID := c.Param("gameID")

	err := h.gamesService.UnVote(gameID, user.ID)
	if err != nil {
		internalError(c, err)
		return
	}

	success(c, nil)
}

func (h *API) GameState(c *gin.Context) {
	user := h.authenticate(c)
	if user == nil {
		return
	}
	gameID := c.Param("gameID")
	lastChangeID := c.Query("lastChangeID")

	state, err := h.gamesService.GameState(gameID, user.ID, 5*time.Second, lastChangeID)
	if err != nil {
		internalError(c, err)
		return
	}

	success(c, state)
}

func (h *API) RestartGame(c *gin.Context) {
	user := h.authenticate(c)
	if user == nil {
		return
	}
	gameID := c.Param("gameID")

	err := h.gamesService.Restart(gameID, user.ID)
	if err != nil {
		internalError(c, err)
		return
	}

	success(c, nil)
}

func (h *API) RevealCards(c *gin.Context) {
	user := h.authenticate(c)
	if user == nil {
		return
	}
	gameID := c.Param("gameID")

	err := h.gamesService.Reveal(gameID, user.ID)
	if err != nil {
		internalError(c, err)
		return
	}

	success(c, nil)
}

func (h *API) authenticate(c *gin.Context) *domain.User {
	auth := c.GetHeader("Authorization")
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		unauthorizedError(c, errors.New("unauthorized"))
		return nil
	}

	user, err := h.usersService.AuthenticateByID(parts[1])
	if err != nil {
		unauthorizedError(c, err)
		return nil
	}

	if user == nil {
		notFoundError(c, errors.New("user not found"))
		return nil
	}

	return user
}
