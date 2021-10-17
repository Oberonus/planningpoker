package http

import (
	"errors"
	"planningpoker/internal/domain/games"
	"planningpoker/internal/domain/state"
	"planningpoker/internal/domain/users"
	"strings"
	"time"

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
	r.GET("/api/v1/me", h.CurrentUser)
	r.PUT("/api/v1/me", h.ChangeUserData)
	r.POST("/api/v1/register", h.Register)
	r.POST("/api/v1/games", h.CreateGame)
	r.PUT("/api/v1/games/:gameID", h.UpdateGame)
	r.POST("/api/v1/games/:gameID/join", h.Join)
	r.POST("/api/v1/games/:gameID/ping", h.Ping)
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

	cmd, err := users.NewRegisterCommand(pl.Name)
	if err != nil {
		badRequestError(c, err)
		return
	}

	user, err := h.usersService.Register(*cmd)
	if err != nil {
		badRequestError(c, err)
		return
	}

	success(c, gin.H{
		"user_id": user.ID(),
	})
}

func (h *API) CurrentUser(c *gin.Context) {
	user := h.authenticate(c)
	if user == nil {
		return
	}

	success(c, gin.H{
		"name": user.Name(),
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

	cmd, err := users.NewUpdateCommand(user.ID(), pl.Name)
	if err != nil {
		badRequestError(c, err)
		return
	}

	_, err = h.usersService.Update(*cmd)
	if err != nil {
		internalError(c, err)
		return
	}

	success(c, gin.H{})
}

type gamePayload struct {
	Name      string `json:"name"`
	TicketURL string `json:"url"`
	CardsDeck struct {
		Name  string   `json:"name"`
		Types []string `json:"types"`
	} `json:"cards_deck"`
	EveryoneCanReveal bool `json:"everyone_can_reveal"`
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

	cards := make([]games.Card, len(pl.CardsDeck.Types))
	for i, v := range pl.CardsDeck.Types {
		card, err := games.NewCard(v)
		if err != nil {
			badRequestError(c, err)
			return
		}
		cards[i] = *card
	}

	deck, err := games.NewCardsDeck(pl.CardsDeck.Name, cards)
	if err != nil {
		badRequestError(c, err)
		return
	}

	cmd, err := games.NewCreateGameCommand(pl.Name, pl.TicketURL, user.ID(), *deck, pl.EveryoneCanReveal)
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

	cmd, err := games.NewUpdateGameCommand(gameID, pl.Name, pl.TicketURL, user.ID())
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

	cmd, err := games.NewJoinGameCommand(gameID, user.ID())
	if err != nil {
		badRequestError(c, err)
		return
	}

	err = h.gamesService.Join(*cmd)
	if err != nil {
		internalError(c, err)
		return
	}

	success(c, nil)
}

func (h *API) Ping(c *gin.Context) {
	user := h.authenticate(c)
	if user == nil {
		return
	}
	gameID := c.Param("gameID")

	cmd, err := games.NewPlayerPingCommand(gameID, user.ID())
	if err != nil {
		badRequestError(c, err)
		return
	}

	err = h.gamesService.Ping(*cmd)
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

	card, err := games.NewCard(vote)
	if err != nil {
		badRequestError(c, err)
		return
	}

	cmd, err := games.NewVoteCommand(gameID, user.ID(), *card)
	if err != nil {
		badRequestError(c, err)
		return
	}

	err = h.gamesService.Vote(*cmd)
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

	cmd, err := games.NewUnVoteCommand(gameID, user.ID())
	if err != nil {
		badRequestError(c, err)
		return
	}

	err = h.gamesService.UnVote(*cmd)
	if err != nil {
		internalError(c, err)
		return
	}

	success(c, nil)
}

type playerStateResponse struct {
	Name      string `json:"name"`
	VotedCard string `json:"voted_card"`
}

func newPlayerStateResponse(state state.PlayerState) playerStateResponse {
	resp := playerStateResponse{
		Name: state.Name,
	}
	if state.VotedCard != nil {
		resp.VotedCard = state.VotedCard.Type()
	}
	return resp
}

type cardsDeckResponse struct {
	Name  string   `json:"name"`
	Cards []string `json:"cards"`
}

func newCardsDeckResponse(cd games.CardsDeck) cardsDeckResponse {
	resp := cardsDeckResponse{Name: cd.Name()}
	for _, c := range cd.Cards() {
		resp.Cards = append(resp.Cards, c.Type())
	}
	return resp
}

type gameStateResponse struct {
	Name      string                `json:"name"`
	TicketURL string                `json:"ticket_url"`
	ChangeID  string                `json:"change_id"`
	CardsDeck cardsDeckResponse     `json:"cards_deck"`
	Players   []playerStateResponse `json:"players"`
	State     string                `json:"state"`
	VotedCard string                `json:"voted_card"`
	CanReveal bool                  `json:"can_reveal"`
}

func newGameStateResponse(state state.GameState) gameStateResponse {
	resp := gameStateResponse{
		Name:      state.Name,
		TicketURL: state.TicketURL,
		ChangeID:  state.ChangeID,
		CardsDeck: newCardsDeckResponse(state.CardsDeck),
		State:     state.State,
		CanReveal: state.CanReveal,
	}
	if state.VotedCard != nil {
		resp.VotedCard = state.VotedCard.Type()
	}
	for _, p := range state.Players {
		resp.Players = append(resp.Players, newPlayerStateResponse(p))
	}
	return resp
}

func (h *API) GameState(c *gin.Context) {
	user := h.authenticate(c)
	if user == nil {
		return
	}
	gameID := c.Param("gameID")
	lastChangeID := c.Query("lastChangeID")

	cmd, err := state.NewGameStateCommand(gameID, user.ID(), 5*time.Second, lastChangeID)
	if err != nil {
		badRequestError(c, err)
		return
	}

	st, err := h.stateService.GameState(*cmd)
	if err != nil {
		internalError(c, err)
		return
	}

	success(c, newGameStateResponse(*st))
}

func (h *API) RestartGame(c *gin.Context) {
	user := h.authenticate(c)
	if user == nil {
		return
	}
	gameID := c.Param("gameID")

	cmd, err := games.NewRestartGameCommand(gameID, user.ID())
	if err != nil {
		badRequestError(c, err)
		return
	}

	err = h.gamesService.Restart(*cmd)
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

	cmd, err := games.NewRevealCardsCommand(gameID, user.ID())
	if err != nil {
		badRequestError(c, err)
		return
	}

	err = h.gamesService.Reveal(*cmd)
	if err != nil {
		internalError(c, err)
		return
	}

	success(c, nil)
}

func (h *API) authenticate(c *gin.Context) *users.User {
	auth := c.GetHeader("Authorization")
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		unauthorizedError(c, errors.New("unauthorized"))
		return nil
	}

	cmd, err := users.NewAuthByIDCommand(parts[1])
	if err != nil {
		badRequestError(c, err)
		return nil
	}

	user, err := h.usersService.AuthenticateByID(*cmd)
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
