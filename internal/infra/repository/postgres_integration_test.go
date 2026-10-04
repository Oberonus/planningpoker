//go:build integration
// +build integration

package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"planningpoker/internal/domain/events"
	"planningpoker/internal/domain/games"
	"planningpoker/internal/domain/users"
	"planningpoker/internal/infra/repository"
)

func testDSN(t *testing.T) string {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")

	if dsn == "" {
		t.Fatal("integration tests require TEST_DATABASE_URL pointing to a disposable database")
	}

	return dsn
}

func openDB(t *testing.T) *sql.DB {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	db, err := repository.OpenPostgres(ctx, testDSN(t))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	return db
}

// Exercise first-start migrations with simultaneous application instances.
func TestPostgresConcurrentStartup(t *testing.T) {
	dsn := testDSN(t)
	errs := make(chan error, 6)

	for i := 0; i < cap(errs); i++ {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			db, err := repository.OpenPostgres(ctx, dsn)

			if err == nil {
				err = db.Close()
			}
			errs <- err
		}()
	}

	for i := 0; i < cap(errs); i++ {
		require.NoError(t, <-errs)
	}
	var version int
	require.NoError(t, openDB(t).QueryRow("SELECT MAX(version) FROM schema_migrations").Scan(&version))
	assert.Equal(t, 1, version)
}

func fixtureGame(t *testing.T) *games.Game {
	t.Helper()

	card, err := games.NewCard("5")
	require.NoError(t, err)
	deck, err := games.NewCardsDeck("test", []games.Card{*card})
	require.NoError(t, err)
	cmd, err := games.NewCreateGameCommand("Persistent game", "https://example.com/ticket", "owner", *deck, true)
	require.NoError(t, err)

	game := games.NewGame(*cmd)
	require.NoError(t, game.Join(games.JoinGameCommand{
		GameID: game.ID(),
		UserID: "owner",
	}))
	require.NoError(t, game.Vote(games.VoteCommand{
		GameID:     game.ID(),
		UserID:     "owner",
		Vote:       *card,
		Confidence: games.ConfidenceNormal,
	}))

	return game
}

func TestPostgresPersistenceAfterReconnect(t *testing.T) {
	db := openDB(t)
	bus := &observedBus{}
	gamesRepo := repository.NewPostgresGameRepository(db, bus)
	usersRepo := repository.NewPostgresUserRepository(db, bus)
	user, err := users.NewUser("Persistent user")
	require.NoError(t, err)
	require.NoError(t, usersRepo.Save(*user))
	game := fixtureGame(t)
	require.NoError(t, gamesRepo.Save(game))
	require.NoError(t, db.Close())

	db = openDB(t)
	gamesRepo = repository.NewPostgresGameRepository(db, bus)
	usersRepo = repository.NewPostgresUserRepository(db, bus)
	loadedUser, err := usersRepo.Get(user.ID())
	require.NoError(t, err)
	require.Equal(t, user.Name(), loadedUser.Name())

	loaded, err := gamesRepo.Get(game.ID())
	require.NoError(t, err)
	require.NotNil(t, loaded)
	assert.Equal(t, game.Name(), loaded.Name())
	assert.Equal(t, game.TicketURL(), loaded.TicketURL())
	assert.Equal(t, game.CardsDeck(), loaded.CardsDeck())
	assert.Equal(t, game.Players(), loaded.Players())
	assert.Equal(t, game.EveryoneCanReveal(), loaded.EveryoneCanReveal())
	assert.Empty(t, loaded.GetEvents())

	missing, err := gamesRepo.Get("missing-'game")
	require.NoError(t, err)
	assert.Nil(t, missing)

	missingUser, err := usersRepo.Get("missing-'user")
	require.NoError(t, err)
	assert.Nil(t, missingUser)

	other, err := users.NewUser("Second user")
	require.NoError(t, err)
	require.NoError(t, usersRepo.Save(*other))
	list, err := usersRepo.GetMany([]string{other.ID(), "missing", user.ID()})
	require.NoError(t, err)
	require.Len(t, list, 2)
	assert.Equal(t, other.ID(), list[0].ID())
	assert.Equal(t, user.ID(), list[1].ID())
	require.NoError(t, other.NameAs("Renamed user"))
	require.NoError(t, usersRepo.Save(*other))
	loadedUser, err = usersRepo.Get(other.ID())
	require.NoError(t, err)
	assert.Equal(t, "Renamed user", loadedUser.Name())

	active, err := gamesRepo.GetActiveGamesByPlayerID("owner")
	require.NoError(t, err)

	found := false

	for _, g := range active {
		found = found || g.ID() == game.ID()
	}

	assert.True(t, found)
	require.NoError(t, gamesRepo.ModifyExclusively(game.ID(), func(g *games.Game) error {
		return g.Reveal(games.RevealCardsCommand{
			GameID: g.ID(),
			UserID: "owner",
		})
	}))

	active, err = gamesRepo.GetActiveGamesByPlayerID("owner")
	require.NoError(t, err)

	for _, g := range active {
		assert.NotEqual(t, game.ID(), g.ID())
	}
}

func TestPostgresRollbackAndCommitEvents(t *testing.T) {
	db := openDB(t)
	bus := &observedBus{}
	r := repository.NewPostgresGameRepository(db, bus)
	game := fixtureGame(t)
	require.NoError(t, r.Save(game))

	bus.count = 0
	errRejected := errors.New("rejected change")
	err := r.ModifyExclusively(game.ID(), func(g *games.Game) error {
		require.NoError(t, g.Update(games.UpdateGameCommand{
			GameID: g.ID(),
			UserID: "owner",
			Name:   "Rejected",
		}))

		return errRejected
	})
	require.ErrorIs(t, err, errRejected)
	loaded, err := r.Get(game.ID())
	require.NoError(t, err)
	assert.Equal(t, game.Name(), loaded.Name())
	assert.Zero(t, bus.count)
	// A synchronous subscriber must be able to see committed state using another connection.
	bus.observe = func(e events.DomainEvent) {
		loaded, err := r.Get(e.AggregateID())
		require.NoError(t, err)
		assert.Equal(t, "Committed", loaded.Name())
	}
	require.NoError(t, r.ModifyExclusively(game.ID(), func(g *games.Game) error {
		return g.Update(games.UpdateGameCommand{
			GameID: g.ID(),
			UserID: "owner",
			Name:   "Committed",
		})
	}))
	assert.Equal(t, 1, bus.count)

	err = r.ModifyExclusively("missing", func(*games.Game) error {
		t.Fatal("callback must not run for a missing game")
		return nil
	})
	require.Error(t, err)
}

func TestPostgresConcurrentVotesAcrossConnections(t *testing.T) {
	bus := &observedBus{}
	firstRepo := repository.NewPostgresGameRepository(openDB(t), bus)
	secondRepo := repository.NewPostgresGameRepository(openDB(t), bus)
	game := fixtureGame(t)
	players := make([]string, 20)

	for i := range players {
		user, err := users.NewUser("Player")
		require.NoError(t, err)

		players[i] = user.ID()
		require.NoError(t, game.Join(games.JoinGameCommand{
			GameID: game.ID(),
			UserID: user.ID(),
		}))
	}

	require.NoError(t, firstRepo.Save(game))

	card, err := games.NewCard("5")
	require.NoError(t, err)

	errs := make(chan error, len(players))
	start := make(chan struct{})

	for i, id := range players {
		r := firstRepo
		if i%2 == 0 {
			r = secondRepo
		}

		go func(r *repository.PostgresGameRepository, id string) {
			<-start
			errs <- r.ModifyExclusively(game.ID(), func(g *games.Game) error {
				return g.Vote(games.VoteCommand{
					GameID:     g.ID(),
					UserID:     id,
					Vote:       *card,
					Confidence: games.ConfidenceNormal,
				})
			})
		}(r, id)
	}

	close(start)

	for range players {
		require.NoError(t, <-errs)
	}

	loaded, err := firstRepo.Get(game.ID())
	require.NoError(t, err)

	for _, id := range players {
		require.NotNil(t, loaded.Players()[id].VotedCard, "vote for %s was lost", id)
		assert.Equal(t, "5", loaded.Players()[id].VotedCard.Type())
	}
}

type observedBus struct {
	mu      sync.Mutex
	count   int
	observe func(events.DomainEvent)
}

func (b *observedBus) Publish(e events.DomainEvent) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.count++

	if b.observe != nil {
		b.observe(e)
	}

	return nil
}

func (*observedBus) Subscribe(events.Consumer, ...string) {}
