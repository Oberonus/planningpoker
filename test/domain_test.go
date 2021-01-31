package test

import (
	"planningpoker/internal/domain"
	"planningpoker/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkflow(t *testing.T) {
	gamesRepo := repository.NewMemoryGameRepository()
	usersRepo := repository.NewMemoryUserRepository()

	gamesService, err := domain.NewGamesService(gamesRepo, usersRepo)
	require.NoError(t, err)
	require.NotNil(t, gamesService)

	usersService, err := domain.NewUsersService(usersRepo)
	require.NoError(t, err)
	require.NotNil(t, usersService)

	user1, err := usersService.Register("John")
	require.NoError(t, err)
	require.NotNil(t, user1)

	user2, err := usersService.Register("Mike")
	require.NoError(t, err)
	require.NotNil(t, user2)

	gameID, err := gamesService.Create(user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, gameID)

	err = gamesService.Join(gameID, user2.ID)
	require.NoError(t, err)

	//err = gamesService.Vote(gameID, user1.ID, "XS")
	//require.NoError(t, err)
	//
	//err = gamesService.Vote(gameID, user2.ID, "?")
	//require.NoError(t, err)

	state, err := gamesService.GameState(gameID, user1.ID)
	require.NoError(t, err)
	require.NotNil(t, state)

	assert.Len(t, state.Players, 2)
	assert.Equal(t, "started", state.State)
	//assert.Equal(t, "XS", state.VotedCard)
	assert.Equal(t, true, state.CanReveal)

	err = gamesService.Reveal(gameID, user1.ID)
	require.NoError(t, err)

	state, err = gamesService.GameState(gameID, user2.ID)
	require.NoError(t, err)
	require.NotNil(t, state)

	assert.Len(t, state.Players, 2)
	assert.Equal(t, "finished", state.State)
	//assert.Equal(t, "?", state.VotedCard)
	assert.Equal(t, false, state.CanReveal)

	t.Logf("%+v", state)
}
