// Package games contains domain level game logic.
package games

import (
	"errors"
	"strings"

	"planningpoker/internal/domain"
	"planningpoker/internal/domain/events"

	"github.com/google/uuid"
)

const (
	// GameStateStarted represents running game state.
	GameStateStarted = "started"
	// GameStateFinished represents finished game state.
	GameStateFinished = "finished"
)

// Game is a domain aggregate that represents one single game.
type Game struct {
	domain.BaseAggregate
	id                string
	name              string
	ticketURL         string
	cardsDeck         CardsDeck
	players           map[string]*Player
	state             string
	everyoneCanReveal bool
}

// Player is an entity of a game player with state.
type Player struct {
	VotedCard *Card
	CanReveal bool
	Active    bool
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
func NewRaw(id, name, ticketURL string, deck CardsDeck, players map[string]*Player, state string, ecr bool) *Game {
	return &Game{
		id:                id,
		name:              name,
		ticketURL:         ticketURL,
		cardsDeck:         deck,
		players:           players,
		state:             state,
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
	g.setChanged()

	if g.IsPlayer(cmd.UserID) {
		g.players[cmd.UserID].Active = true
		return nil
	}

	canReveal := g.everyoneCanReveal || len(g.players) == 0

	g.players[cmd.UserID] = &Player{
		VotedCard: nil,
		CanReveal: canReveal,
		Active:    true,
	}

	return nil
}

// Leave marks a player as inactive or removes them from players depending on voting state.
func (g *Game) Leave(cmd LeaveGameCommand) error {
	p, ok := g.players[cmd.UserID]
	if !ok {
		return nil
	}

	g.setChanged()

	// if the player is voted, we don't want to delete the data until cards not revealed.
	if p.VotedCard != nil {
		p.Active = false
		return nil
	}

	delete(g.players, cmd.UserID)

	return nil
}

// Restart resets the game state.
func (g *Game) Restart(cmd RestartGameCommand) error {
	_, ok := g.players[cmd.UserID]
	if !ok {
		return errors.New("user is not a player")
	}
	g.state = GameStateStarted

	// cleanup non-active players and remove votes
	players := g.players
	for id, p := range g.players {
		if !p.Active {
			delete(players, id)
			continue
		}
		players[id].VotedCard = nil
	}
	g.players = players

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

func (g *Game) setChanged() {
	g.AddEvent(events.NewDomainEventBuilder(events.EventTypeGameUpdated).ForAggregate(g.id).Build())
}
