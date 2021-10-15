package test

import (
	"planningpoker/internal/domain"
	"planningpoker/internal/domain/users"
	"planningpoker/internal/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkflow(t *testing.T) {
	gamesRepo := repository.NewMemoryGameRepository()
	usersRepo := repository.NewMemoryUserRepository()

	gamesService, err := domain.NewGamesService(gamesRepo, usersRepo)
	require.NoError(t, err)
	require.NotNil(t, gamesService)

	usersService, err := users.NewService(usersRepo)
	require.NoError(t, err)
	require.NotNil(t, usersService)

	regCmd, err := users.NewRegisterCommand("John")
	require.NoError(t, err)
	user1, err := usersService.Register(*regCmd)
	require.NoError(t, err)
	require.NotNil(t, user1)

	regCmd, err = users.NewRegisterCommand("Mike")
	require.NoError(t, err)
	user2, err := usersService.Register(*regCmd)
	require.NoError(t, err)
	require.NotNil(t, user2)

	cmd, err := domain.NewCreateGameCommand("a", "b", user1.ID(), newTestCardsDeck(t), false)
	require.NoError(t, err)

	gameID, err := gamesService.Create(*cmd)
	require.NoError(t, err)
	require.NotEmpty(t, gameID)

	joinCmd, err := domain.NewJoinGameCommand(gameID, user2.ID())
	require.NoError(t, err)
	err = gamesService.Join(*joinCmd)
	require.NoError(t, err)

	voteCmd, err := domain.NewVoteCommand(gameID, user1.ID(), "XS")
	require.NoError(t, err)
	err = gamesService.Vote(*voteCmd)
	require.NoError(t, err)

	voteCmd, err = domain.NewVoteCommand(gameID, user2.ID(), "?")
	require.NoError(t, err)
	err = gamesService.Vote(*voteCmd)
	require.NoError(t, err)

	stateCmd, err := domain.NewGameStateCommand(gameID, user1.ID(), time.Second, "")
	require.NoError(t, err)
	state, err := gamesService.GameState(*stateCmd)
	require.NoError(t, err)
	require.NotNil(t, state)

	assert.Len(t, state.Players, 2)
	assert.Equal(t, "started", state.State)
	assert.Equal(t, "XS", string(*state.VotedCard))
	assert.Equal(t, true, state.CanReveal)

	revealCmd, err := domain.NewRevealCardsCommand(gameID, user1.ID())
	require.NoError(t, err)
	err = gamesService.Reveal(*revealCmd)
	require.NoError(t, err)

	stateCmd, err = domain.NewGameStateCommand(gameID, user2.ID(), time.Second, "")
	require.NoError(t, err)
	state, err = gamesService.GameState(*stateCmd)
	require.NoError(t, err)
	require.NotNil(t, state)

	assert.Len(t, state.Players, 2)
	assert.Equal(t, "finished", state.State)
	assert.Equal(t, "?", string(*state.VotedCard))
	assert.Equal(t, false, state.CanReveal)

	t.Logf("%+v", state)
}

func newTestCardsDeck(t *testing.T) domain.CardsDeck {
	types := []string{"XS", "?"}
	cards := make([]domain.Card, len(types))
	for i, v := range types {
		c, err := domain.NewCard(v)
		require.NoError(t, err)
		cards[i] = *c
	}

	deck, err := domain.NewCardsDeck("test", cards)
	require.NoError(t, err)

	return *deck
}
