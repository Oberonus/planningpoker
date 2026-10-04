//go:build component

package component_test

import (
	"context"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"planningpoker/internal/domain/games"
	"planningpoker/internal/domain/state"
	"planningpoker/internal/domain/users"
	"planningpoker/internal/infra/async"
	"planningpoker/internal/infra/auth"
	"planningpoker/internal/infra/eventbus"
	api "planningpoker/internal/infra/http"
	"planningpoker/internal/infra/repository"
)

// Run the public-client suite against the real adapters and PostgreSQL. Keeping
// the server in this process lets Go collect coverage across the whole stack.
func TestService(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	require.NotEmpty(t, dsn, "component tests require a disposable TEST_DATABASE_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := repository.OpenPostgres(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })

	bus := eventbus.NewInternalBus()
	gamesRepo := repository.NewPostgresGameRepository(db, bus)
	usersRepo := repository.NewPostgresUserRepository(db, bus)
	gamesService, err := games.NewService(gamesRepo, bus)
	require.NoError(t, err)
	usersService, err := users.NewService(usersRepo)
	require.NoError(t, err)

	authenticator := auth.NewUserAuthenticator(usersService)
	httpAPI, err := api.NewAPI(usersService, authenticator)
	require.NoError(t, err)
	httpAPI.SetReadinessCheck(db.PingContext)

	socketAPI := async.NewAPI(gamesService, authenticator)

	t.Cleanup(func() { require.NoError(t, socketAPI.Close()) })

	_, err = state.NewService(gamesRepo, usersRepo, socketAPI, bus)
	require.NoError(t, err)

	router := gin.New()
	router.Use(gin.Recovery())
	httpAPI.SetupRoutes(router)
	socketAPI.SetupRoutes(router)
	server := httptest.NewServer(router)
	t.Cleanup(server.Close)

	clientCtx, clientCancel := context.WithTimeout(context.Background(), time.Minute)
	defer clientCancel()

	command := exec.CommandContext(clientCtx, "node", "--test", "../../web/tests/component/service.test.js")

	command.Env = append(os.Environ(), "POKER_TEST_URL="+server.URL)

	output, err := command.CombinedOutput()
	t.Log(string(output))
	require.NoError(t, err, "black-box service tests failed")
}
