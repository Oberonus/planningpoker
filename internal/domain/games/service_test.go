package games_test

import (
	"errors"
	"planningpoker/internal/domain/events"
	"planningpoker/internal/domain/games"
	"planningpoker/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		gameRepo games.GameRepository
		eventBus events.EventBus
		expError string
	}{
		"success": {
			gameRepo: gamesRepoStub{},
			eventBus: eventBusStub{},
			expError: "",
		},
		"fail on no game repo": {
			eventBus: eventBusStub{},
			expError: "games repository should be provided",
		},
		"fail on no event bus": {
			gameRepo: gamesRepoStub{},
			expError: "event bus should be provided",
		},
	}
	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			srv, err := games.NewService(tt.gameRepo, tt.eventBus)

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
		gameRepo games.GameRepository
		expError string
	}{
		"success": {
			gameRepo: gamesRepoStub{},
			expError: "",
		},
		"fail on repo error": {
			gameRepo: gamesRepoStub{saveError: errors.New("save failed")},
			expError: "save failed",
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			srv, err := games.NewService(tt.gameRepo, eventBusStub{})
			require.NoError(t, err)

			cmd, err := games.NewCreateGameCommand("foo", "http://example.com", test.User1, test.NewTestDeck(t), true)
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
		gameRepo games.GameRepository
		expError string
	}{
		"success": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).UserJoins(test.User1).Instance()},
			expError: "",
		},
		"fail on error": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).Instance()},
			expError: "user is not a player",
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			srv, err := games.NewService(tt.gameRepo, eventBusStub{})
			require.NoError(t, err)

			cmd, err := games.NewUpdateGameCommand("anything", "new name", "https://ex.com", test.User1)
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
		gameRepo games.GameRepository
		expError string
	}{
		"success": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).UserJoins(test.User1).Instance()},
			expError: "",
		},
		"fail on error": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).Instance()},
			expError: "user is not a player",
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			srv, err := games.NewService(tt.gameRepo, eventBusStub{})
			require.NoError(t, err)

			cmd, err := games.NewRestartGameCommand("anything", test.User1)
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
	card, err := games.NewCard("XS")
	require.NoError(t, err)

	testCases := map[string]struct {
		gameRepo games.GameRepository
		expError string
	}{
		"success": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).UserJoins(test.User1).Instance()},
			expError: "",
		},
		"fail on error": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).Instance()},
			expError: "user is not a player",
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			srv, err := games.NewService(tt.gameRepo, eventBusStub{})
			require.NoError(t, err)

			cmd, err := games.NewVoteCommand("anything", test.User1, *card)
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

func TestGamesService_Ping(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		gameRepo games.GameRepository
		expError string
	}{
		"success": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).UserJoins(test.User1).Instance()},
			expError: "",
		},
		"fail on error": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).Instance()},
			expError: "user is not a player",
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			srv, err := games.NewService(tt.gameRepo, eventBusStub{})
			require.NoError(t, err)

			cmd, err := games.NewPlayerPingCommand("anything", test.User1)
			require.NoError(t, err)

			err = srv.Ping(*cmd)

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
		gameRepo games.GameRepository
		expError string
	}{
		"success": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).UserJoins(test.User1).UserVotes(test.User1, "XS").Instance()},
			expError: "",
		},
		"fail on error": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).Instance()},
			expError: "user is not a player",
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			srv, err := games.NewService(tt.gameRepo, eventBusStub{})
			require.NoError(t, err)

			cmd, err := games.NewUnVoteCommand("anything", test.User1)
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
		gameRepo games.GameRepository
		expError string
	}{
		"success": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).UserJoins(test.User1).Instance()},
			expError: "",
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			srv, err := games.NewService(tt.gameRepo, eventBusStub{})
			require.NoError(t, err)

			cmd, err := games.NewJoinGameCommand("anything", test.User2)
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
		gameRepo games.GameRepository
		expError string
	}{
		"success": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).UserJoins(test.User1).Instance()},
			expError: "",
		},
		"fail on error": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).Instance()},
			expError: "user is not a player",
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			srv, err := games.NewService(tt.gameRepo, eventBusStub{})
			require.NoError(t, err)

			cmd, err := games.NewRevealCardsCommand("anything", test.User1)
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

type gamesRepoStub struct {
	game              *games.Game
	getErr            error
	saveError         error
	activeGames       []games.Game
	getActiveGamesErr error
}

func (g gamesRepoStub) ModifyExclusively(_ string, cb func(game *games.Game) error) error {
	return cb(g.game)
}

func (g gamesRepoStub) Get(string) (*games.Game, error) {
	return g.game, g.getErr
}

func (g gamesRepoStub) Save(*games.Game) error {
	return g.saveError
}

func (g gamesRepoStub) GetActiveGamesByPlayerID(playerID string) ([]games.Game, error) {
	return g.activeGames, g.getActiveGamesErr
}

func newTestServiceGame(t *testing.T) *test.Game {
	return test.NewTestGame(t, test.NewSimpleGame(t, true))
}

type eventBusStub struct{}

func (e eventBusStub) Publish(events.DomainEvent) error {
	return nil
}

func (e eventBusStub) Subscribe(events.Consumer, ...string) {
}
