package state_test

import (
	"errors"
	"planningpoker/internal/domain/games"
	"planningpoker/internal/domain/state"
	"planningpoker/internal/domain/users"
	"planningpoker/test"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		gameRepo  state.GameRepository
		usersRepo state.UsersRepository
		expError  string
	}{
		"success": {
			gameRepo:  gamesRepoStub{},
			usersRepo: usersRepoStub{},
			expError:  "",
		},
		"fail on no game repo": {
			usersRepo: usersRepoStub{},
			expError:  "games repository should be provided",
		},
		"fail on no user repo": {
			gameRepo: gamesRepoStub{},
			expError: "users repository should be provided",
		},
	}
	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			srv, err := state.NewService(tt.gameRepo, tt.usersRepo)

			if tt.expError != "" {
				assert.EqualError(t, err, tt.expError)
				assert.Nil(t, srv)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, srv)
			}
		})
	}
}

func TestGamesService_GameState(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		gameRepo state.GameRepository
		userRepo state.UsersRepository
		expError string
	}{
		"success": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).UserJoins(test.User1).Instance()},
			userRepo: usersRepoStub{},
			expError: "",
		},
		"fail on error": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).Instance()},
			userRepo: usersRepoStub{},
			expError: "user is not a player",
		},
		"fail on games repo error": {
			gameRepo: gamesRepoStub{
				game:   newTestServiceGame(t).UserJoins(test.User1).Instance(),
				getErr: errors.New("get failed"),
			},
			userRepo: usersRepoStub{},
			expError: "get failed",
		},
		"fail on users repo error": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).UserJoins(test.User1).Instance()},
			userRepo: usersRepoStub{getManyErr: errors.New("users failed")},
			expError: "users failed",
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			srv, err := state.NewService(tt.gameRepo, tt.userRepo)
			require.NoError(t, err)

			cmd, err := state.NewGameStateCommand("anything", test.User1, time.Millisecond, "")
			require.NoError(t, err)

			state, err := srv.GameState(*cmd)

			if tt.expError != "" {
				assert.EqualError(t, err, tt.expError)
				assert.Nil(t, state)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, state)
				assert.Len(t, state.Players, 1)
			}
		})
	}
}

type gamesRepoStub struct {
	game   *games.Game
	getErr error
}

func (g gamesRepoStub) ModifyExclusively(_ string, cb func(game *games.Game) error) error {
	return cb(g.game)
}

func (g gamesRepoStub) Get(string) (*games.Game, error) {
	return g.game, g.getErr
}

type usersRepoStub struct {
	getManyErr error
	manyUsers  []users.User
}

func (u usersRepoStub) GetMany([]string) ([]users.User, error) {
	return u.manyUsers, u.getManyErr
}

func newTestServiceGame(t *testing.T) *test.Game {
	return test.NewTestGame(t, test.NewSimpleGame(t, true))
}
