package http

import (
	"planningpoker/internal/domain/games"
	"planningpoker/internal/domain/users"

	"github.com/gin-gonic/gin"
)

type gamePayload struct {
	Name      string `json:"name"`
	TicketURL string `json:"url"`
	CardsDeck struct {
		Name  string   `json:"name"`
		Types []string `json:"types"`
	} `json:"cards_deck"`
	EveryoneCanReveal bool `json:"everyone_can_reveal"`
}

func (h *API) CreateGame(c *gin.Context, user *users.User) {
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

func (h *API) UpdateGame(c *gin.Context, user *users.User) {
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

func (h *API) Join(c *gin.Context, user *users.User) {
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

func (h *API) Ping(c *gin.Context, user *users.User) {
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

func (h *API) Vote(c *gin.Context, user *users.User) {
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

func (h *API) UnVote(c *gin.Context, user *users.User) {
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

func (h *API) RevealCards(c *gin.Context, user *users.User) {
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

func (h *API) RestartGame(c *gin.Context, user *users.User) {
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
