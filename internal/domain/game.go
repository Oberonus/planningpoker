package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Player struct {
	VotedCard string
	CanReveal bool
	LastPing  time.Time
}

var (
	CfgAllCanReveal = true
)

const (
	GameStateStarted  = "started"
	GameStateFinished = "finished"

	StalePlayerTTL = 10 * time.Second
)

var TShirtCards = []string{"XXS", "XS", "S", "M", "L", "XL", "XXL", "?"}

type Game struct {
	ID       string
	Cards    []string
	Players  map[string]*Player
	State    string
	ChangeID string
}

func NewTShirtGame() *Game {
	return &Game{
		ID:      strings.Replace(uuid.New().String(), "-", "", -1),
		Cards:   TShirtCards,
		Players: make(map[string]*Player),
		State:   GameStateStarted,
	}
}

func (g *Game) Join(uid string) error {
	if g.IsPlayer(uid) {
		return nil
	}

	// first joiner will be an admin
	canReveal := CfgAllCanReveal || len(g.Players) == 0

	g.Players[uid] = &Player{
		VotedCard: "",
		CanReveal: canReveal,
		LastPing:  time.Now(),
	}

	g.setChanged()

	return nil
}

func (g *Game) Restart(uid string) error {
	_, ok := g.Players[uid]
	if !ok {
		return errors.New("user is not a player")
	}
	g.State = GameStateStarted

	for _, p := range g.Players {
		p.VotedCard = ""
	}

	g.setChanged()

	return nil
}

func (g *Game) Vote(uid string, vote string) error {
	if !g.IsPlayer(uid) {
		return errors.New("user is not a player")
	}

	if g.State != GameStateStarted {
		return errors.New("can not vote on ended game")
	}

	found := false
	for _, card := range g.Cards {
		if card == vote {
			found = true
			break
		}
	}
	if !found {
		return errors.New("unknown card")
	}

	g.Players[uid].VotedCard = vote

	g.setChanged()

	return nil
}

func (g *Game) Reveal(uid string) error {
	p, ok := g.Players[uid]
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

func (g *Game) UnVote(uid string) error {
	if !g.IsPlayer(uid) {
		return errors.New("user is not a player")
	}

	if g.State != GameStateStarted {
		return errors.New("can not un-vote on ended game")
	}

	g.Players[uid].VotedCard = ""

	g.setChanged()

	return nil
}

func (g *Game) Leave(uid string) error {
	if !g.IsPlayer(uid) {
		return errors.New("user is not a player")
	}
	delete(g.Players, uid)
	g.setChanged()
	return nil
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
		if p.LastPing.Add(StalePlayerTTL).After(time.Now()) || p.VotedCard != "" {
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

func (g *Game) CurrentState(userID string, allUsers map[string]*User) *GameState {
	state := &GameState{
		Cards:     g.Cards,
		Players:   make([]PlayerState, 0, len(g.Players)),
		State:     g.State,
		VotedCard: g.Players[userID].VotedCard,
		CanReveal: g.Players[userID].CanReveal,
		ChangeID:  g.ChangeID,
	}

	for uid, p := range g.Players {
		userName := "Unknown"
		if user, ok := allUsers[uid]; ok {
			userName = user.Name
		}

		votedCard := p.VotedCard
		// mask real votes if game is running
		if g.State == GameStateStarted && votedCard != "" {
			votedCard = "*"
		}

		state.Players = append(state.Players, PlayerState{
			Name:      userName,
			VotedCard: votedCard,
		})
	}

	return state
}

func (g *Game) setChanged() {
	g.ChangeID = uuid.NewString()
}
