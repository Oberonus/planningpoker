package games

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Player struct {
	VotedCard *Card
	CanReveal bool
	LastPing  time.Time
}

const (
	GameStateStarted  = "started"
	GameStateFinished = "finished"

	StalePlayerTTL = 10 * time.Second
)

type Game struct {
	id                string
	name              string
	ticketURL         string
	cardsDeck         CardsDeck
	players           map[string]*Player
	state             string
	changeID          string
	everyoneCanReveal bool
}

func NewGame(cmd CreateGameCommand) *Game {
	return &Game{
		id:                strings.ReplaceAll(uuid.New().String(), "-", ""),
		name:              cmd.Name,
		ticketURL:         cmd.TicketURL,
		cardsDeck:         cmd.CardsDeck,
		players:           make(map[string]*Player),
		state:             GameStateStarted,
		everyoneCanReveal: cmd.EveryoneCanReveal,
	}
}

func NewRaw(id, name, ticketURL string, deck CardsDeck, players map[string]*Player, state, cid string, ecr bool) *Game {
	return &Game{
		id:                id,
		name:              name,
		ticketURL:         ticketURL,
		cardsDeck:         deck,
		players:           players,
		state:             state,
		changeID:          cid,
		everyoneCanReveal: ecr,
	}
}

func (g Game) ID() string {
	return g.id
}

func (g Game) Name() string {
	return g.name
}

func (g Game) TicketURL() string {
	return g.ticketURL
}

func (g Game) CardsDeck() CardsDeck {
	return g.cardsDeck
}

func (g Game) Players() map[string]*Player {
	return g.players
}

func (g Game) State() string {
	return g.state
}

func (g Game) EveryoneCanReveal() bool {
	return g.everyoneCanReveal
}

func (g Game) ChangeID() string {
	return g.changeID
}

func (g *Game) Update(cmd UpdateGameCommand) error {
	_, ok := g.players[cmd.UserID]
	if !ok {
		return errors.New("user is not a player")
	}

	g.name = cmd.Name
	g.ticketURL = cmd.TicketURL
	g.setChanged()

	return nil
}

func (g *Game) Join(cmd JoinGameCommand) error {
	if g.IsPlayer(cmd.UserID) {
		return nil
	}

	canReveal := g.everyoneCanReveal || len(g.players) == 0

	g.players[cmd.UserID] = &Player{
		VotedCard: nil,
		CanReveal: canReveal,
		LastPing:  time.Now(),
	}

	g.setChanged()

	return nil
}

func (g *Game) Restart(cmd RestartGameCommand) error {
	_, ok := g.players[cmd.UserID]
	if !ok {
		return errors.New("user is not a player")
	}
	g.state = GameStateStarted

	for _, p := range g.players {
		p.VotedCard = nil
	}

	g.setChanged()

	return nil
}

func (g *Game) Vote(cmd VoteCommand) error {
	if !g.IsPlayer(cmd.UserID) {
		return errors.New("user is not a player")
	}

	if g.state != GameStateStarted {
		return errors.New("can not vote on ended game")
	}

	if !g.cardsDeck.IsInDeck(cmd.Vote) {
		return errors.New("unknown card")
	}

	g.players[cmd.UserID].VotedCard = &cmd.Vote

	g.setChanged()

	return nil
}

func (g *Game) Reveal(cmd RevealCardsCommand) error {
	p, ok := g.players[cmd.UserID]
	if !ok {
		return errors.New("user is not a player")
	}

	if !p.CanReveal {
		return errors.New("user can not reveal cards")
	}

	g.state = GameStateFinished

	g.setChanged()

	return nil
}

func (g *Game) UnVote(cmd UnVoteCommand) error {
	if !g.IsPlayer(cmd.UserID) {
		return errors.New("user is not a player")
	}

	if g.state != GameStateStarted {
		return errors.New("can not un-vote on ended game")
	}

	g.players[cmd.UserID].VotedCard = nil

	g.setChanged()

	return nil
}

func (g *Game) ForceChanged() {
	g.setChanged()
}

func (g *Game) IsPlayer(uid string) bool {
	_, ok := g.players[uid]
	return ok
}

// Ping additionally does a players cleanup. Who does not ping for several seconds, will be removed.
func (g *Game) Ping(uid string) error {
	p, ok := g.players[uid]
	if !ok {
		return errors.New("user is not a player")
	}
	p.LastPing = time.Now()

	// cleanup stale players who did not vote
	newPlayers := make(map[string]*Player)
	hasRevealer := false
	for id, p := range g.players {
		if p.LastPing.Add(StalePlayerTTL).After(time.Now()) || p.VotedCard != nil {
			newPlayers[id] = p
			if p.CanReveal {
				hasRevealer = true
			}
		}
	}

	// choose random admin since master left
	if !hasRevealer {
		for k := range newPlayers {
			newPlayers[k].CanReveal = true
			break
		}
	}

	if len(g.players) != len(newPlayers) {
		g.setChanged()
	}

	g.players = newPlayers

	return nil
}

func (g *Game) setChanged() {
	g.changeID = uuid.NewString()
}
