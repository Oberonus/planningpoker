package test_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"planningpoker/internal/domain/games"
	"planningpoker/internal/domain/state"
	"planningpoker/internal/domain/users"
	"planningpoker/internal/infra/eventbus"
	"planningpoker/internal/infra/repository"
)

func TestWorkflow(t *testing.T) {
	eventBus := eventbus.NewInternalBus()
	gamesRepo := repository.NewMemoryGameRepository(eventBus)
	usersRepo := repository.NewMemoryUserRepository(eventBus)

	gamesService, err := games.NewService(gamesRepo, eventBus)
	require.NoError(t, err)
	require.NotNil(t, gamesService)

	usersService, err := users.NewService(usersRepo)
	require.NoError(t, err)
	require.NotNil(t, usersService)

	stateService, err := state.NewService(gamesRepo, usersRepo)
	require.NoError(t, err)
	require.NotNil(t, stateService)

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

	cmd, err := games.NewCreateGameCommand("a", "b", user1.ID(), newTestCardsDeck(t), false)
	require.NoError(t, err)

	gameID, err := gamesService.Create(*cmd)
	require.NoError(t, err)
	require.NotEmpty(t, gameID)

	joinCmd, err := games.NewJoinGameCommand(gameID, user2.ID())
	require.NoError(t, err)
	err = gamesService.Join(*joinCmd)
	require.NoError(t, err)

	voteCmd, err := games.NewVoteCommand(gameID, user1.ID(), "XS")
	require.NoError(t, err)
	err = gamesService.Vote(*voteCmd)
	require.NoError(t, err)

	voteCmd, err = games.NewVoteCommand(gameID, user2.ID(), "?")
	require.NoError(t, err)
	err = gamesService.Vote(*voteCmd)
	require.NoError(t, err)

	st, err := stateService.GameState(gameID)
	require.NoError(t, err)
	require.NotNil(t, st)

	assert.Len(t, st.Players, 2)

	revealCmd, err := games.NewRevealCardsCommand(gameID, user1.ID())
	require.NoError(t, err)
	err = gamesService.Reveal(*revealCmd)
	require.NoError(t, err)

	st, err = stateService.GameState(gameID)
	require.NoError(t, err)
	require.NotNil(t, st)

	assert.Len(t, st.Players, 2)
	assert.Equal(t, "finished", st.State)
}

func newTestCardsDeck(t *testing.T) games.CardsDeck {
	types := []string{"XS", "?"}
	cards := make([]games.Card, len(types))
	for i, v := range types {
		c, err := games.NewCard(v)
		require.NoError(t, err)
		cards[i] = *c
	}

	deck, err := games.NewCardsDeck("test", cards)
	require.NoError(t, err)

	return *deck
}
