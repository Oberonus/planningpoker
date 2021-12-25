package test

import (
	"planningpoker/internal/domain/games"
	"testing"

	"github.com/stretchr/testify/require"
)

type Game struct {
	game      *games.Game
	lastError error
	t         *testing.T
}

var (
	User1 = "user-id-1"
	User2 = "user-id-2"
)

func NewTestGame(t *testing.T, game *games.Game) *Game {
	return &Game{game: game, t: t}
}

func (g *Game) When() *Game {
	return g
}

func (g *Game) And() *Game {
	return g
}

func (g *Game) Then() *Game {
	return g
}

func (g *Game) Instance() *games.Game {
	return g.game
}

func (g *Game) UserJoins(uid string) *Game {
	cmd, err := games.NewJoinGameCommand(g.game.ID(), uid)
	require.NoError(g.t, err)
	g.lastError = g.game.Join(*cmd)
	return g
}

func (g *Game) UserVotes(uid string, cardName string) *Game {
	card, err := games.NewCard(cardName)
	require.NoError(g.t, err)
	cmd, err := games.NewVoteCommand(g.game.ID(), uid, *card)
	require.NoError(g.t, err)
	g.lastError = g.game.Vote(*cmd)
	return g
}

func (g *Game) UserUnVotes(uid string) *Game {
	cmd, err := games.NewUnVoteCommand(g.game.ID(), uid)
	require.NoError(g.t, err)
	g.lastError = g.game.UnVote(*cmd)
	return g
}

func (g *Game) UserRestartsGame(uid string) *Game {
	cmd, err := games.NewRestartGameCommand(g.game.ID(), uid)
	require.NoError(g.t, err)
	g.lastError = g.game.Restart(*cmd)
	return g
}

func (g *Game) ShouldHaveVote(uid string, cardName string) *Game {
	card, err := games.NewCard(cardName)
	require.NoError(g.t, err)
	require.Equal(g.t, card, g.game.Players()[uid].VotedCard)
	return g
}

func (g *Game) ShouldHaveNoVote(uid string) *Game {
	require.Nil(g.t, g.game.Players()[uid].VotedCard)
	return g
}

func (g *Game) UserReveals(uid string) *Game {
	cmd, err := games.NewRevealCardsCommand(g.game.ID(), uid)
	require.NoError(g.t, err)
	g.lastError = g.game.Reveal(*cmd)
	return g
}

func (g *Game) GameShouldBeFinished() *Game {
	require.Equal(g.t, games.GameStateFinished, g.game.State())
	return g
}

func (g *Game) GameShouldBeRunning() *Game {
	require.Equal(g.t, games.GameStateStarted, g.game.State())
	return g
}

func (g *Game) ShouldFail(part string) *Game {
	require.Error(g.t, g.lastError)
	require.Contains(g.t, g.lastError.Error(), part)
	return g
}

func (g *Game) ShouldSucceed() *Game {
	require.NoError(g.t, g.lastError)
	return g
}

func (g *Game) UserUpdatesGameName(uid, name string) *Game {
	cmd, err := games.NewUpdateGameCommand(g.game.ID(), name, g.game.TicketURL(), uid)
	require.NoError(g.t, err)
	g.lastError = g.game.Update(*cmd)
	return g
}

func (g *Game) ShouldHaveGameName(name string) *Game {
	require.Equal(g.t, name, g.game.Name())
	return g
}

func NewSimpleGame(t *testing.T, everybodyCanReveal bool) *games.Game {
	cmd, err := games.NewCreateGameCommand("", "", "", NewTestDeck(t), everybodyCanReveal)
	require.NoError(t, err)
	return games.NewGame(*cmd)
}

func NewTestDeck(t *testing.T) games.CardsDeck {
	card1, err := games.NewCard("XS")
	require.NoError(t, err)
	card2, err := games.NewCard("S")
	require.NoError(t, err)
	deck, err := games.NewCardsDeck("T-shirt", []games.Card{*card1, *card2})
	require.NoError(t, err)

	return *deck
}
