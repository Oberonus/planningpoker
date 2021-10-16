package domain_test

import (
	"errors"
	"planningpoker/internal/domain"
	"planningpoker/internal/domain/events"
	"planningpoker/internal/domain/users"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGame(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		gameRepo domain.GameRepository
		userRepo domain.UsersRepository
		expError string
	}{
		"success": {
			gameRepo: gamesRepoStub{},
			userRepo: usersRepoStub{},
			expError: "",
		},
		"fail on no user repo": {
			gameRepo: gamesRepoStub{saveError: errors.New("save failed")},
			expError: "users repository should be provided",
		},
		"fail on no game repo": {
			userRepo: usersRepoStub{},
			expError: "games repository should be provided",
		},
	}
	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			srv, err := domain.NewGamesService(tt.gameRepo, tt.userRepo, eventBusStub{})

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

func TestGamesService_Create(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		gameRepo domain.GameRepository
		userRepo domain.UsersRepository
		expError string
	}{
		"success": {
			gameRepo: gamesRepoStub{},
			userRepo: usersRepoStub{},
			expError: "",
		},
		"fail on repo error": {
			gameRepo: gamesRepoStub{saveError: errors.New("save failed")},
			userRepo: usersRepoStub{},
			expError: "save failed",
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			srv, err := domain.NewGamesService(tt.gameRepo, tt.userRepo, eventBusStub{})
			require.NoError(t, err)

			cmd, err := domain.NewCreateGameCommand("foo", "http://example.com", User1, newTestDeck(t), true)
			require.NoError(t, err)

			id, err := srv.Create(*cmd)

			if tt.expError != "" {
				assert.EqualError(t, err, tt.expError)
				assert.Empty(t, id)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, id)
			}
		})
	}
}

func TestGamesService_Update(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		gameRepo domain.GameRepository
		userRepo domain.UsersRepository
		expError string
	}{
		"success": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).UserJoins(User1).Instance()},
			userRepo: usersRepoStub{},
			expError: "",
		},
		"fail on error": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).Instance()},
			userRepo: usersRepoStub{},
			expError: "user is not a player",
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			srv, err := domain.NewGamesService(tt.gameRepo, tt.userRepo, eventBusStub{})
			require.NoError(t, err)

			cmd, err := domain.NewUpdateGameCommand("anything", "new name", "https://ex.com", User1)
			require.NoError(t, err)

			err = srv.Update(*cmd)

			if tt.expError != "" {
				assert.EqualError(t, err, tt.expError)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func TestGamesService_Restart(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		gameRepo domain.GameRepository
		userRepo domain.UsersRepository
		expError string
	}{
		"success": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).UserJoins(User1).Instance()},
			userRepo: usersRepoStub{},
			expError: "",
		},
		"fail on error": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).Instance()},
			userRepo: usersRepoStub{},
			expError: "user is not a player",
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			srv, err := domain.NewGamesService(tt.gameRepo, tt.userRepo, eventBusStub{})
			require.NoError(t, err)

			cmd, err := domain.NewRestartGameCommand("anything", User1)
			require.NoError(t, err)

			err = srv.Restart(*cmd)

			if tt.expError != "" {
				assert.EqualError(t, err, tt.expError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGamesService_Vote(t *testing.T) {
	t.Parallel()
	card, err := domain.NewCard("XS")
	require.NoError(t, err)

	testCases := map[string]struct {
		gameRepo domain.GameRepository
		userRepo domain.UsersRepository
		expError string
	}{
		"success": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).UserJoins(User1).Instance()},
			userRepo: usersRepoStub{},
			expError: "",
		},
		"fail on error": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).Instance()},
			userRepo: usersRepoStub{},
			expError: "user is not a player",
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			srv, err := domain.NewGamesService(tt.gameRepo, tt.userRepo, eventBusStub{})
			require.NoError(t, err)

			cmd, err := domain.NewVoteCommand("anything", User1, *card)
			require.NoError(t, err)

			err = srv.Vote(*cmd)

			if tt.expError != "" {
				assert.EqualError(t, err, tt.expError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGamesService_UnVote(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		gameRepo domain.GameRepository
		userRepo domain.UsersRepository
		expError string
	}{
		"success": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).UserJoins(User1).UserVotes(User1, "XS").Instance()},
			userRepo: usersRepoStub{},
			expError: "",
		},
		"fail on error": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).Instance()},
			userRepo: usersRepoStub{},
			expError: "user is not a player",
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			srv, err := domain.NewGamesService(tt.gameRepo, tt.userRepo, eventBusStub{})
			require.NoError(t, err)

			cmd, err := domain.NewUnVoteCommand("anything", User1)
			require.NoError(t, err)

			err = srv.UnVote(*cmd)

			if tt.expError != "" {
				assert.EqualError(t, err, tt.expError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGamesService_Join(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		gameRepo domain.GameRepository
		userRepo domain.UsersRepository
		expError string
	}{
		"success": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).UserJoins(User1).Instance()},
			userRepo: usersRepoStub{},
			expError: "",
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			srv, err := domain.NewGamesService(tt.gameRepo, tt.userRepo, eventBusStub{})
			require.NoError(t, err)

			cmd, err := domain.NewJoinGameCommand("anything", User2)
			require.NoError(t, err)

			err = srv.Join(*cmd)

			if tt.expError != "" {
				assert.EqualError(t, err, tt.expError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGamesService_Reveal(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		gameRepo domain.GameRepository
		userRepo domain.UsersRepository
		expError string
	}{
		"success": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).UserJoins(User1).Instance()},
			userRepo: usersRepoStub{},
			expError: "",
		},
		"fail on error": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).Instance()},
			userRepo: usersRepoStub{},
			expError: "user is not a player",
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			srv, err := domain.NewGamesService(tt.gameRepo, tt.userRepo, eventBusStub{})
			require.NoError(t, err)

			cmd, err := domain.NewRevealCardsCommand("anything", User1)
			require.NoError(t, err)

			err = srv.Reveal(*cmd)

			if tt.expError != "" {
				assert.EqualError(t, err, tt.expError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGamesService_GameState(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		gameRepo domain.GameRepository
		userRepo domain.UsersRepository
		expError string
	}{
		"success": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).UserJoins(User1).Instance()},
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
				game:   newTestServiceGame(t).UserJoins(User1).Instance(),
				getErr: errors.New("get failed"),
			},
			userRepo: usersRepoStub{},
			expError: "get failed",
		},
		"fail on users repo error": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).UserJoins(User1).Instance()},
			userRepo: usersRepoStub{getManyErr: errors.New("users failed")},
			expError: "users failed",
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			srv, err := domain.NewGamesService(tt.gameRepo, tt.userRepo, eventBusStub{})
			require.NoError(t, err)

			cmd, err := domain.NewGameStateCommand("anything", User1, time.Millisecond, "")
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
	game              *domain.Game
	getErr            error
	saveError         error
	activeGames       []domain.Game
	getActiveGamesErr error
}

func (g gamesRepoStub) ModifyExclusively(_ string, cb func(game *domain.Game) error) error {
	return cb(g.game)
}

func (g gamesRepoStub) Get(string) (*domain.Game, error) {
	return g.game, g.getErr
}

func (g gamesRepoStub) Save(*domain.Game) error {
	return g.saveError
}

func (g gamesRepoStub) GetActiveGamesByPlayerID(playerID string) ([]domain.Game, error) {
	return g.activeGames, g.getActiveGamesErr
}

type usersRepoStub struct {
	getManyErr error
	manyUsers  []users.User
}

func (u usersRepoStub) GetMany([]string) ([]users.User, error) {
	return u.manyUsers, u.getManyErr
}

func newTestServiceGame(t *testing.T) *testGame {
	return NewTestGame(t, newSimpleGame(t, true))
}

type eventBusStub struct{}

func (e eventBusStub) Publish(events.DomainEvent) error {
	return nil
}

func (e eventBusStub) Subscribe(events.Consumer, ...string) {
	return
}
