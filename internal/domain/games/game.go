// Package games contains domain level game logic.
package games

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	// GameStateStarted represents running game state.
	GameStateStarted = "started"
	// GameStateFinished represents finished game state.
	GameStateFinished = "finished"

	// StalePlayerTTL is a timeout when a player will be automatically removed from a game in case of no Pings.
	StalePlayerTTL = 10 * time.Second
)

// Game is a domain aggregate that represents one single game.
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

// Player is an entity of a game player with state.
type Player struct {
	VotedCard *Card
	CanReveal bool
	LastPing  time.Time
}

// NewGame creates a new game aggregate instance.
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

// NewRaw instantiates a game aggregate from raw data.
// It should never be used in any logic except aggregate hydration from any serialized format (db, etc...)
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

// ID returns a game id.
func (g Game) ID() string {
	return g.id
}

// Name returns a game name.
func (g Game) Name() string {
	return g.name
}

// TicketURL returns a ticket URL (e.g. link to JIRA) related to the game.
func (g Game) TicketURL() string {
	return g.ticketURL
}

// CardsDeck returns a deck of cards used in game.
func (g Game) CardsDeck() CardsDeck {
	return g.cardsDeck
}

// Players returns all game players.
func (g Game) Players() map[string]*Player {
	return g.players
}

// State returns a game current state.
func (g Game) State() string {
	return g.state
}

// EveryoneCanReveal returns true if any player can reveal cards.
func (g Game) EveryoneCanReveal() bool {
	return g.everyoneCanReveal
}

// ChangeID returns internal change ID which shows that the aggregate was changed.
func (g Game) ChangeID() string {
	return g.changeID
}

// Update updates game generic data.
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

// Join adds a new player to the game.
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

// Restart resets the game state.
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

// Vote performs a player voting.
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

// Reveal opens all cards and stops the game.
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

// UnVote removes a vote for a passenger.
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

// ForceChanged marks the aggregate as changes (dirty state).
func (g *Game) ForceChanged() {
	g.setChanged()
}

// IsPlayer checks if specific user is a player.
func (g *Game) IsPlayer(uid string) bool {
	_, ok := g.players[uid]
	return ok
}

// Ping maintains players state and does a players cleanup if needed.
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
