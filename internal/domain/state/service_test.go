package state_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"planningpoker/internal/domain/events"
	"planningpoker/internal/domain/games"
	"planningpoker/internal/domain/state"
	"planningpoker/internal/domain/users"
	"planningpoker/test"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		gameRepo  state.GameRepository
		usersRepo state.UsersRepository
		publisher state.Publisher
		eventBus  events.EventBus
		expError  string
	}{
		"success": {
			gameRepo:  gamesRepoStub{},
			usersRepo: usersRepoStub{},
			publisher: publisherStub{},
			eventBus:  eventBusStub{},
			expError:  "",
		},
		"fail on no event bus": {
			gameRepo:  gamesRepoStub{},
			usersRepo: usersRepoStub{},
			publisher: publisherStub{},
			expError:  "event bus should be provided",
		},
		"fail on no publisher": {
			gameRepo:  gamesRepoStub{},
			usersRepo: usersRepoStub{},
			eventBus:  eventBusStub{},
			expError:  "publisher should be provided",
		},
		"fail on no game repo": {
			usersRepo: usersRepoStub{},
			publisher: publisherStub{},
			eventBus:  eventBusStub{},
			expError:  "games repository should be provided",
		},
		"fail on no user repo": {
			gameRepo:  gamesRepoStub{},
			publisher: publisherStub{},
			eventBus:  eventBusStub{},
			expError:  "users repository should be provided",
		},
	}

	for name, tt := range testCases {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			srv, err := state.NewService(tt.gameRepo, tt.usersRepo, tt.publisher, tt.eventBus)

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
		gameRepo  state.GameRepository
		userRepo  state.UsersRepository
		eventBus  events.EventBus
		publisher state.Publisher
		expError  string
	}{
		"success": {
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).UserJoins(test.User1).Instance()},
			userRepo: usersRepoStub{},
			expError: "",
		},
		"fail on games repo error": {
			gameRepo: gamesRepoStub{
				game:   newTestServiceGame(t).UserJoins(test.User1).Instance(),
				getErr: errors.New("get failed"),
			},
			userRepo: usersRepoStub{},
			expError: "get game: get failed",
		},
		"fail on missing game": {
			gameRepo: gamesRepoStub{},
			userRepo: usersRepoStub{},
			expError: "game not found",
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

			srv, err := state.NewService(tt.gameRepo, tt.userRepo, publisherStub{}, eventBusStub{})
			require.NoError(t, err)

			st, err := srv.GameState("anything")

			if tt.expError != "" {
				assert.EqualError(t, err, tt.expError)
				assert.Nil(t, st)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, st)
				assert.Len(t, st.Players, 1)
			}
		})
	}
}

func TestGameStateEventHandlesRepositoryFailures(t *testing.T) {
	for _, tc := range []struct {
		name     string
		gameRepo gamesRepoStub
		userRepo usersRepoStub
	}{
		{
			name:     "game read fails",
			gameRepo: gamesRepoStub{getErr: errors.New("database unavailable")},
		},
		{
			name: "game missing",
		},
		{
			name:     "user read fails",
			gameRepo: gamesRepoStub{game: newTestServiceGame(t).UserJoins(test.User1).Instance()},
			userRepo: usersRepoStub{getManyErr: errors.New("database unavailable")},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			bus := &capturingEventBus{}
			_, err := state.NewService(tc.gameRepo, tc.userRepo, unexpectedPublisher{t: t}, bus)
			require.NoError(t, err)
			require.NotNil(t, bus.consumer)

			event := events.NewDomainEventBuilder(events.EventTypeGameUpdated).ForAggregate("game").Build()

			assert.NotPanics(t, func() { bus.consumer(event) })
		})
	}
}

type capturingEventBus struct {
	eventBusStub
	consumer events.Consumer
}

func (b *capturingEventBus) Subscribe(consumer events.Consumer, _ ...string) {
	b.consumer = consumer
}

type unexpectedPublisher struct {
	t *testing.T
}

func (p unexpectedPublisher) SendToPlayer(state.GameState, string) error {
	p.t.Error("a failed state read must not publish")
	return nil
}

type gamesRepoStub struct {
	game   *games.Game
	getErr error
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

type publisherStub struct {
	err error
}

func (p publisherStub) SendToPlayer(gameState state.GameState, userID string) error {
	return p.err
}

type eventBusStub struct{}

func (e eventBusStub) Publish(event events.DomainEvent) error {
	return nil
}

func (e eventBusStub) Subscribe(consumer events.Consumer, eventTypes ...string) {
}
