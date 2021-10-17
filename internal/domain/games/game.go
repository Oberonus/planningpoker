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
	ID                string
	Name              string
	TicketURL         string
	CardsDeck         CardsDeck
	Players           map[string]*Player
	State             string
	ChangeID          string
	EveryoneCanReveal bool
}

func NewGame(cmd CreateGameCommand) *Game {
	return &Game{
		ID:                strings.Replace(uuid.New().String(), "-", "", -1),
		Name:              cmd.Name,
		TicketURL:         cmd.TicketURL,
		CardsDeck:         cmd.CardsDeck,
		Players:           make(map[string]*Player),
		State:             GameStateStarted,
		EveryoneCanReveal: cmd.EveryoneCanReveal,
	}
}

func (g *Game) Update(cmd UpdateGameCommand) error {
	_, ok := g.Players[cmd.UserID]
	if !ok {
		return errors.New("user is not a player")
	}

	g.Name = cmd.Name
	g.TicketURL = cmd.TicketURL
	g.setChanged()

	return nil
}

func (g *Game) Join(cmd JoinGameCommand) error {
	if g.IsPlayer(cmd.UserID) {
		return nil
	}

	canReveal := g.EveryoneCanReveal || len(g.Players) == 0

	g.Players[cmd.UserID] = &Player{
		VotedCard: nil,
		CanReveal: canReveal,
		LastPing:  time.Now(),
	}

	g.setChanged()

	return nil
}

func (g *Game) Restart(cmd RestartGameCommand) error {
	_, ok := g.Players[cmd.UserID]
	if !ok {
		return errors.New("user is not a player")
	}
	g.State = GameStateStarted

	for _, p := range g.Players {
		p.VotedCard = nil
	}

	g.setChanged()

	return nil
}

func (g *Game) Vote(cmd VoteCommand) error {
	if !g.IsPlayer(cmd.UserID) {
		return errors.New("user is not a player")
	}

	if g.State != GameStateStarted {
		return errors.New("can not vote on ended game")
	}

	if !g.CardsDeck.IsInDeck(cmd.Vote) {
		return errors.New("unknown card")
	}

	g.Players[cmd.UserID].VotedCard = &cmd.Vote

	g.setChanged()

	return nil
}

func (g *Game) Reveal(cmd RevealCardsCommand) error {
	p, ok := g.Players[cmd.UserID]
	if !ok {
		return errors.New("user is not a player")
	}

	if !p.CanReveal {
		return errors.New("user can not reveal cards")
	}

	g.State = GameStateFinished

	g.setChanged()

	return nil
}

func (g *Game) UnVote(cmd UnVoteCommand) error {
	if !g.IsPlayer(cmd.UserID) {
		return errors.New("user is not a player")
	}

	if g.State != GameStateStarted {
		return errors.New("can not un-vote on ended game")
	}

	g.Players[cmd.UserID].VotedCard = nil

	g.setChanged()

	return nil
}

func (g *Game) ForceChanged() {
	g.setChanged()
}

func (g *Game) IsPlayer(uid string) bool {
	_, ok := g.Players[uid]
	return ok
}

// Ping additionally does a players cleanup. Who does not ping for several seconds, will be removed.
func (g *Game) Ping(uid string) error {
	p, ok := g.Players[uid]
	if !ok {
		return errors.New("user is not a player")
	}
	p.LastPing = time.Now()

	// cleanup stale players who did not vote
	newPlayers := make(map[string]*Player)
	hasRevealer := false
	for id, p := range g.Players {
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

	if len(g.Players) != len(newPlayers) {
		g.setChanged()
	}

	g.Players = newPlayers

	return nil
}

func (g *Game) setChanged() {
	g.ChangeID = uuid.NewString()
}
