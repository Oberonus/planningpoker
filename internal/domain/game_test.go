package domain_test

import (
	"planningpoker/internal/domain"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimpleGame(t *testing.T) {
	NewTestGame(t, newSimpleGame(t, false)).
		When().UserJoins(User1).
		And().UserVotes(User1, "XS").
		Then().ShouldHaveVote(User1, "XS").
		When().UserJoins(User2).
		And().UserVotes(User2, "S").
		Then().ShouldHaveVote(User2, "S").
		When().UserReveals(User1).
		Then().GameShouldBeFinished()
}

func TestFailedToReveal(t *testing.T) {
	NewTestGame(t, newSimpleGame(t, false)).
		When().UserJoins(User1).
		And().UserJoins(User2).
		And().UserReveals(User2).
		Then().ShouldFail("user can not reveal cards").
		And().GameShouldBeRunning()
}

func TestEveryoneCanReveal(t *testing.T) {
	NewTestGame(t, newSimpleGame(t, true)).
		When().UserJoins(User1).
		And().UserJoins(User2).
		And().UserReveals(User2).
		Then().GameShouldBeFinished()
}

func TestUserCanJoinTwice(t *testing.T) {
	NewTestGame(t, newSimpleGame(t, true)).
		When().UserJoins(User1).
		And().UserJoins(User1).
		Then().ShouldSucceed()
}

func TestRestartGame(t *testing.T) {
	NewTestGame(t, newSimpleGame(t, true)).
		When().UserJoins(User1).
		And().UserVotes(User1, "XS").
		And().UserReveals(User1).
		Then().GameShouldBeFinished().
		And().ShouldHaveVote(User1, "XS").
		When().UserRestartsGame(User1).
		Then().ShouldHaveNoVote(User1).
		And().GameShouldBeRunning()
}

func TestNonPlayerCanNotRestartAGame(t *testing.T) {
	NewTestGame(t, newSimpleGame(t, true)).
		When().UserJoins(User1).
		And().UserRestartsGame(User2).
		Then().ShouldFail("user is not a player")
}

func TestCanNotVoteWhenGameFinished(t *testing.T) {
	NewTestGame(t, newSimpleGame(t, false)).
		When().UserJoins(User1).
		And().UserReveals(User1).
		And().UserVotes(User1, "XS").
		Then().ShouldFail("can not vote on ended game")
}

func TestCannotVoteWhenNotAPlayer(t *testing.T) {
	NewTestGame(t, newSimpleGame(t, false)).
		When().UserJoins(User1).
		And().UserVotes(User2, "XS").
		Then().ShouldFail("user is not a player")
}

func TestCannotVoteWithWrongCard(t *testing.T) {
	NewTestGame(t, newSimpleGame(t, false)).
		When().UserJoins(User1).
		And().UserVotes(User1, "foo").
		Then().ShouldFail("unknown card")
}

func TestCanUnVote(t *testing.T) {
	NewTestGame(t, newSimpleGame(t, false)).
		When().UserJoins(User1).
		And().UserVotes(User1, "XS").
		Then().ShouldHaveVote(User1, "XS").
		When().UserUnVotes(User1).
		Then().ShouldHaveNoVote(User1)
}

func TestNonPlayerCanNotUnVote(t *testing.T) {
	NewTestGame(t, newSimpleGame(t, false)).
		When().UserJoins(User1).
		And().UserUnVotes(User2).
		Then().ShouldFail("user is not a player")
}

func TestCanNotUnVoteOnFinishedGame(t *testing.T) {
	NewTestGame(t, newSimpleGame(t, false)).
		When().UserJoins(User1).
		And().UserVotes(User1, "XS").
		And().UserReveals(User1).
		When().UserUnVotes(User1).
		Then().ShouldFail("can not un-vote on ended game")
}

func TestCannotRevealWhenNotAPlayer(t *testing.T) {
	NewTestGame(t, newSimpleGame(t, false)).
		When().UserJoins(User1).
		And().UserReveals(User2).
		Then().ShouldFail("user is not a player")
}

func TestPlayerCanUpdateGameName(t *testing.T) {
	NewTestGame(t, newSimpleGame(t, false)).
		When().UserJoins(User1).
		And().UserUpdatesGameName(User1, "new name").
		Then().ShouldHaveGameName("new name")
}

func TestCanNotUpdateGameNameWhenNotAPlayer(t *testing.T) {
	NewTestGame(t, newSimpleGame(t, false)).
		When().UserJoins(User1).
		And().UserUpdatesGameName(User2, "new name").
		Then().ShouldFail("user is not a player")
}

type testGame struct {
	game      *domain.Game
	lastError error
	t         *testing.T
}

var (
	User1 = "user-id-1"
	User2 = "user-id-2"
)

func NewTestGame(t *testing.T, game *domain.Game) *testGame {
	return &testGame{game: game, t: t}
}

func (g *testGame) When() *testGame {
	return g
}

func (g *testGame) And() *testGame {
	return g
}

func (g *testGame) Then() *testGame {
	return g
}

func (g *testGame) Instance() *domain.Game {
	return g.game
}

func (g *testGame) UserJoins(uid string) *testGame {
	cmd, err := domain.NewJoinGameCommand(g.game.ID, uid)
	require.NoError(g.t, err)
	g.lastError = g.game.Join(*cmd)
	return g
}

func (g *testGame) UserVotes(uid string, cardName string) *testGame {
	card, err := domain.NewCard(cardName)
	require.NoError(g.t, err)
	cmd, err := domain.NewVoteCommand(g.game.ID, uid, *card)
	require.NoError(g.t, err)
	g.lastError = g.game.Vote(*cmd)
	return g
}

func (g *testGame) UserUnVotes(uid string) *testGame {
	cmd, err := domain.NewUnVoteCommand(g.game.ID, uid)
	require.NoError(g.t, err)
	g.lastError = g.game.UnVote(*cmd)
	return g
}

func (g *testGame) UserRestartsGame(uid string) *testGame {
	cmd, err := domain.NewRestartGameCommand(g.game.ID, uid)
	require.NoError(g.t, err)
	g.lastError = g.game.Restart(*cmd)
	return g
}

func (g *testGame) ShouldHaveVote(uid string, cardName string) *testGame {
	card, err := domain.NewCard(cardName)
	require.NoError(g.t, err)
	require.Equal(g.t, card, g.game.Players[uid].VotedCard)
	return g
}

func (g *testGame) ShouldHaveNoVote(uid string) *testGame {
	require.Nil(g.t, g.game.Players[uid].VotedCard)
	return g
}

func (g *testGame) UserReveals(uid string) *testGame {
	cmd, err := domain.NewRevealCardsCommand(g.game.ID, uid)
	require.NoError(g.t, err)
	g.lastError = g.game.Reveal(*cmd)
	return g
}

func (g *testGame) GameShouldBeFinished() *testGame {
	require.Equal(g.t, domain.GameStateFinished, g.game.State)
	return g
}

func (g *testGame) GameShouldBeRunning() *testGame {
	require.Equal(g.t, domain.GameStateStarted, g.game.State)
	return g
}

func (g *testGame) ShouldFail(part string) *testGame {
	require.Error(g.t, g.lastError)
	require.Contains(g.t, g.lastError.Error(), part)
	return g
}

func (g *testGame) ShouldSucceed() *testGame {
	require.NoError(g.t, g.lastError)
	return g
}

func (g *testGame) UserUpdatesGameName(uid, name string) *testGame {
	cmd, err := domain.NewUpdateGameCommand(g.game.ID, name, g.game.TicketURL, uid)
	require.NoError(g.t, err)
	g.lastError = g.game.Update(*cmd)
	return g
}

func (g *testGame) ShouldHaveGameName(name string) *testGame {
	require.Equal(g.t, name, g.game.Name)
	return g
}

func newSimpleGame(t *testing.T, everybodyCanReveal bool) *domain.Game {
	cmd, err := domain.NewCreateGameCommand("", "", "", newTestDeck(t), everybodyCanReveal)
	require.NoError(t, err)
	return domain.NewGame(*cmd)
}

func newTestDeck(t *testing.T) domain.CardsDeck {
	card1, err := domain.NewCard("XS")
	require.NoError(t, err)
	card2, err := domain.NewCard("S")
	require.NoError(t, err)
	deck, err := domain.NewCardsDeck("T-shirt", []domain.Card{*card1, *card2})
	require.NoError(t, err)

	return *deck
}
