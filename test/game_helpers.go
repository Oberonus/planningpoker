package test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"planningpoker/internal/domain/games"
)

// Game is a testing wrapper around domain game aggregate.
type Game struct {
	game      *games.Game
	lastError error
	t         *testing.T
}

var (
	// User1 is a dummy id for testing user.
	User1 = "user-id-1"
	// User2 is a dummy id for testing user.
	User2 = "user-id-2"
)

// NewTestGame creates a new testing game.
func NewTestGame(t *testing.T, game *games.Game) *Game {
	return &Game{game: game, t: t}
}

// When is just a glue for when/then/and flow.
func (g *Game) When() *Game {
	return g
}

// And is just a glue for when/then/and flow.
func (g *Game) And() *Game {
	return g
}

// Then is just a glue for when/then/and flow.
func (g *Game) Then() *Game {
	return g
}

// Instance returns the real game aggregate instance.
func (g *Game) Instance() *games.Game {
	return g.game
}

// UserJoins preform user joining the game by id.
func (g *Game) UserJoins(uid string) *Game {
	cmd, err := games.NewJoinGameCommand(g.game.ID(), uid)
	require.NoError(g.t, err)
	g.lastError = g.game.Join(*cmd)
	return g
}

// UserVotes performs a user vote.
func (g *Game) UserVotes(uid, cardName string) *Game {
	card, err := games.NewCard(cardName)
	require.NoError(g.t, err)
	cmd, err := games.NewVoteCommand(g.game.ID(), uid, *card)
	require.NoError(g.t, err)
	g.lastError = g.game.Vote(*cmd)
	return g
}

// UserLeaves performs a user leave from game.
func (g *Game) UserLeaves(uid string) *Game {
	cmd, err := games.NewLeaveGameCommand(g.game.ID(), uid)
	require.NoError(g.t, err)
	g.lastError = g.game.Leave(*cmd)
	return g
}

// UserUnVotes performs a user unvote.
func (g *Game) UserUnVotes(uid string) *Game {
	cmd, err := games.NewUnVoteCommand(g.game.ID(), uid)
	require.NoError(g.t, err)
	g.lastError = g.game.UnVote(*cmd)
	return g
}

// UserRestartsGame performs a game restart.
func (g *Game) UserRestartsGame(uid string) *Game {
	cmd, err := games.NewRestartGameCommand(g.game.ID(), uid)
	require.NoError(g.t, err)
	g.lastError = g.game.Restart(*cmd)
	return g
}

// ShouldHaveVote asserts that specific user voted with specific card.
func (g *Game) ShouldHaveVote(uid string, cardName string) *Game {
	card, err := games.NewCard(cardName)
	require.NoError(g.t, err)
	require.Equal(g.t, card, g.game.Players()[uid].VotedCard)
	return g
}

// ShouldHaveNoVote asserts that a user did not vote.
func (g *Game) ShouldHaveNoVote(uid string) *Game {
	require.Nil(g.t, g.game.Players()[uid].VotedCard)
	return g
}

// UserReveals performs game end with cards revealing.
func (g *Game) UserReveals(uid string) *Game {
	cmd, err := games.NewRevealCardsCommand(g.game.ID(), uid)
	require.NoError(g.t, err)
	g.lastError = g.game.Reveal(*cmd)
	return g
}

// GameShouldBeFinished asserts that game is finished.
func (g *Game) GameShouldBeFinished() *Game {
	require.Equal(g.t, games.GameStateFinished, g.game.State())
	return g
}

// GameShouldBeRunning asserts that game is in progress.
func (g *Game) GameShouldBeRunning() *Game {
	require.Equal(g.t, games.GameStateStarted, g.game.State())
	return g
}

// ShouldFail asserts that there was an error during the previous step.
func (g *Game) ShouldFail(part string) *Game {
	require.Error(g.t, g.lastError)
	require.Contains(g.t, g.lastError.Error(), part)
	return g
}

// ShouldSucceed asserts that there were no error during the previous step.
func (g *Game) ShouldSucceed() *Game {
	require.NoError(g.t, g.lastError)
	return g
}

// UserUpdatesGameName performs game name update.
func (g *Game) UserUpdatesGameName(uid, name string) *Game {
	cmd, err := games.NewUpdateGameCommand(g.game.ID(), name, g.game.TicketURL(), uid)
	require.NoError(g.t, err)
	g.lastError = g.game.Update(*cmd)
	return g
}

// ShouldHaveGameName asserts that the game has specific name.
func (g *Game) ShouldHaveGameName(name string) *Game {
	require.Equal(g.t, name, g.game.Name())
	return g
}

// NewSimpleGame creates a simple testing game.
func NewSimpleGame(t *testing.T, everybodyCanReveal bool) *games.Game {
	cmd, err := games.NewCreateGameCommand("", "", "", NewTestDeck(t), everybodyCanReveal)
	require.NoError(t, err)
	return games.NewGame(*cmd)
}

// NewTestDeck creates a simple testing cards deck.
func NewTestDeck(t *testing.T) games.CardsDeck {
	card1, err := games.NewCard("XS")
	require.NoError(t, err)
	card2, err := games.NewCard("S")
	require.NoError(t, err)
	deck, err := games.NewCardsDeck("T-shirt", []games.Card{*card1, *card2})
	require.NoError(t, err)

	return *deck
}
