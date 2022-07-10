// Package async contains all socketio related functionality.
package async

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"planningpoker/internal/infra/transformers"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"

	"planningpoker/internal/domain/games"
)

const rootNameSpace = "/"
const genericErrorMessage = "error"

// API is a socket.io API implementation.
type API struct {
	server       *socketio.Server
	usersAuth    userAuthenticator
	gamesService gameService
}

// GameRepository is a contract to fetch games data.
type gameService interface {
	Create(cmd games.CreateGameCommand) (string, error)
	Join(cmd games.JoinGameCommand) error
	Leave(cmd games.LeaveGameCommand) error
	Update(cmd games.UpdateGameCommand) error
	Vote(cmd games.VoteCommand) error
	UnVote(cmd games.UnVoteCommand) error
	Reveal(cmd games.RevealCardsCommand) error
	Restart(cmd games.RestartGameCommand) error
}

type conContext struct {
	userID string
	gameID string
}

// UserAuthenticator is a contract to authenticate users.
type userAuthenticator interface {
	// AuthenticateByToken returns a user ID or error is user is not authenticated
	AuthenticateByToken(token string) (string, error)
}

// NewAPI creates a new socket.io related api.
func NewAPI(repository gameService, authenticator userAuthenticator) *API {
	p := &API{
		gamesService: repository,
		usersAuth:    authenticator,
		server:       socketio.NewServer(nil),
	}

	go func() {
		if err := p.server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	// defer server.Close()

	p.server.OnConnect(rootNameSpace, p.onConnect)
	p.server.OnDisconnect(rootNameSpace, p.onDisconnect)
	p.server.OnEvent(rootNameSpace, "create", p.create)
	p.server.OnEvent(rootNameSpace, "join", p.join)
	p.server.OnEvent(rootNameSpace, "leave", p.leave)
	p.server.OnEvent(rootNameSpace, "vote", p.vote)
	p.server.OnEvent(rootNameSpace, "update", p.update)
	p.server.OnEvent(rootNameSpace, "unvote", p.unVote)
	p.server.OnEvent(rootNameSpace, "reveal", p.reveal)
	p.server.OnEvent(rootNameSpace, "restart", p.restart)
	p.server.OnError(rootNameSpace, func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})
	return p
}

// SetupRoutes sets up socket.io related routes.
func (p *API) SetupRoutes(r gin.IRoutes) {
	r.GET("/socket.io/*any", gin.WrapH(p.server))
	r.POST("/socket.io/*any", gin.WrapH(p.server))
}

// SendToPlayer sends a message to a player of specific game.
func (p *API) SendToPlayer(gameID, userID string, state transformers.GameStateResponse) error {
	ok := p.server.BroadcastToRoom(rootNameSpace, gameID+userID, "gameState", state)
	if !ok {
		return fmt.Errorf("broadcast to gameID=%s failed", gameID)
	}
	return nil
}

func (p *API) onConnect(conn socketio.Conn) error {
	urlVal := conn.URL()
	token := urlVal.Query().Get("token")

	uid, err := p.usersAuth.AuthenticateByToken(token)
	if err != nil {
		logrus.Errorf("socket authentication error: %v", err)
		return errors.New("unauthorized")
	}

	conn.SetContext(conContext{userID: uid})

	return nil
}

func (p *API) onDisconnect(conn socketio.Conn, reason string) {
	uid, ok := conn.Context().(string)
	if !ok {
		logrus.Infof("client unknown id disconnected with reason: %s", reason)
	}

	logrus.Infof("client with id=%s disconnected with reason: %s", uid, reason)
}

type createPayload struct {
	Name      string `json:"name"`
	TicketURL string `json:"url"`
	CardsDeck struct {
		Name  string   `json:"name"`
		Types []string `json:"types"`
	} `json:"cards_deck"`
	EveryoneCanReveal bool `json:"everyone_can_reveal"`
}

func (p *API) create(conn socketio.Conn, pl createPayload) interface{} {
	cc, ok := conn.Context().(conContext)
	if !ok {
		logrus.Errorf("socket listen for updates: unable to get the context")
		return genericErrorMessage
	}

	cards := make([]games.Card, len(pl.CardsDeck.Types))
	for i, v := range pl.CardsDeck.Types {
		card, err := games.NewCard(v)
		if err != nil {
			return genericErrorMessage
		}
		cards[i] = *card
	}

	deck, err := games.NewCardsDeck(pl.CardsDeck.Name, cards)
	if err != nil {
		return genericErrorMessage
	}

	cmd, err := games.NewCreateGameCommand(pl.Name, pl.TicketURL, cc.userID, *deck, pl.EveryoneCanReveal)
	if err != nil {
		return genericErrorMessage
	}

	gameID, err := p.gamesService.Create(*cmd)
	if err != nil {
		return genericErrorMessage
	}

	return gin.H{
		"game_id": gameID,
	}
}

func (p *API) join(conn socketio.Conn, gameID string) string {
	cc, ok := conn.Context().(conContext)
	if !ok {
		logrus.Errorf("socket listen for updates: unable to get the context")
		return genericErrorMessage
	}

	cmd, err := games.NewJoinGameCommand(gameID, cc.userID)
	if err != nil {
		logrus.Errorf("socket listen for updates: unable to create player join command")
		return genericErrorMessage
	}

	// join the room in order to receive personalized game updates
	conn.Join(gameID + cc.userID)

	if err := p.gamesService.Join(*cmd); err != nil {
		logrus.Errorf("socket listen for updates: unable to find the game id=%s: %v", gameID, err)
		return genericErrorMessage
	}

	cc.gameID = gameID
	conn.SetContext(cc)

	return "ok"
}

func (p *API) leave(conn socketio.Conn) string {
	cc, ok := conn.Context().(conContext)
	if !ok {
		logrus.Errorf("socket game: unable to get the context")
		return genericErrorMessage
	}
	conn.Leave(cc.gameID + cc.userID)

	cmd, err := games.NewLeaveGameCommand(cc.gameID, cc.userID)
	if err != nil {
		logrus.Errorf("leave: %v", err)
		return genericErrorMessage
	}

	if err := p.gamesService.Leave(*cmd); err != nil {
		logrus.Errorf("leave: %v", err)
		return genericErrorMessage
	}

	return "ok"
}

func (p *API) vote(conn socketio.Conn, vote string) string {
	cc, ok := conn.Context().(conContext)
	if !ok {
		logrus.Errorf("socket game: unable to get the context")
		return genericErrorMessage
	}

	vote, err := url.QueryUnescape(vote)
	if err != nil {
		logrus.Errorf("socket leave game: failed to unescape vote: %v", err)
		return genericErrorMessage
	}

	card, err := games.NewCard(vote)
	if err != nil {
		logrus.Errorf("vote: %v", err)
		return genericErrorMessage
	}

	cmd, err := games.NewVoteCommand(cc.gameID, cc.userID, *card)
	if err != nil {
		logrus.Errorf("vote: %v", err)
		return genericErrorMessage
	}

	err = p.gamesService.Vote(*cmd)
	if err != nil {
		logrus.Errorf("vote: %v", err)
		return genericErrorMessage
	}

	return ""
}

type updatePayload struct {
	Name      string `json:"name"`
	TicketURL string `json:"ticket_url"`
}

func (p *API) update(conn socketio.Conn, params updatePayload) string {
	cc, ok := conn.Context().(conContext)
	if !ok {
		logrus.Errorf("socket game: unable to get the context")
		return genericErrorMessage
	}

	cmd, err := games.NewUpdateGameCommand(cc.gameID, params.Name, params.TicketURL, cc.userID)
	if err != nil {
		logrus.Errorf("update: %v", err)
		return genericErrorMessage
	}

	if err := p.gamesService.Update(*cmd); err != nil {
		logrus.Errorf("update: %v", err)
		return genericErrorMessage
	}

	return "ok"
}

func (p *API) reveal(conn socketio.Conn) string {
	cc, ok := conn.Context().(conContext)
	if !ok {
		logrus.Errorf("socket game: unable to get the context")
		return genericErrorMessage
	}

	cmd, err := games.NewRevealCardsCommand(cc.gameID, cc.userID)
	if err != nil {
		logrus.Errorf("reveal: %v", err)
		return genericErrorMessage
	}

	if err := p.gamesService.Reveal(*cmd); err != nil {
		logrus.Errorf("reveal: %v", err)
		return genericErrorMessage
	}

	return "ok"
}

func (p *API) unVote(conn socketio.Conn) string {
	cc, ok := conn.Context().(conContext)
	if !ok {
		logrus.Errorf("socket game: unable to get the context")
		return genericErrorMessage
	}

	cmd, err := games.NewUnVoteCommand(cc.gameID, cc.userID)
	if err != nil {
		logrus.Errorf("unvote: %v", err)
		return genericErrorMessage
	}

	if err := p.gamesService.UnVote(*cmd); err != nil {
		logrus.Errorf("unvote: %v", err)
		return genericErrorMessage
	}

	return "ok"
}

func (p *API) restart(conn socketio.Conn) string {
	cc, ok := conn.Context().(conContext)
	if !ok {
		logrus.Errorf("socket game: unable to get the context")
		return genericErrorMessage
	}

	cmd, err := games.NewRestartGameCommand(cc.gameID, cc.userID)
	if err != nil {
		logrus.Errorf("restart: %v", err)
		return genericErrorMessage
	}

	if err := p.gamesService.Restart(*cmd); err != nil {
		logrus.Errorf("restart: %v", err)
		return genericErrorMessage
	}

	return "ok"
}
