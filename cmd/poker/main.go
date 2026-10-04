package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"planningpoker/internal/domain/games"
	"planningpoker/internal/domain/state"
	"planningpoker/internal/domain/users"
	"planningpoker/internal/infra/async"
	"planningpoker/internal/infra/auth"
	"planningpoker/internal/infra/eventbus"
	"planningpoker/internal/infra/http"
	"planningpoker/internal/infra/repository"
)

func main() {
	logrus.Infof("starting the service")

	eventBus := eventbus.NewInternalBus()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	db, err := repository.OpenPostgres(ctx, os.Getenv("DATABASE_URL"))

	cancel()

	if err != nil {
		log.Fatalf("unable to initialize storage: %v", err)
	}

	defer func() { _ = db.Close() }()

	gamesRepo := repository.NewPostgresGameRepository(db, eventBus)
	usersRepo := repository.NewPostgresUserRepository(db, eventBus)

	gamesService, err := games.NewService(gamesRepo, eventBus)
	if err != nil {
		log.Fatalf("unable to create games service: %v", err)
	}

	usersService, err := users.NewService(usersRepo)
	if err != nil {
		log.Fatalf("unable to create users service: %v", err)
	}

	authenticator := auth.NewUserAuthenticator(usersService)

	api, err := http.NewAPI(usersService, authenticator)
	if err != nil {
		log.Fatalf("unable to create http API: %v", err)
	}

	api.SetReadinessCheck(db.PingContext)

	fe := http.NewFrontend()

	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	asyncAPI := async.NewAPI(gamesService, authenticator)

	_, err = state.NewService(gamesRepo, usersRepo, asyncAPI, eventBus)
	if err != nil {
		log.Fatalf("unable to create game state service: %v", err)
	}

	api.SetupRoutes(r)
	asyncAPI.SetupRoutes(r)
	fe.SetupRoutes(r)

	if err := r.Run(); err != nil {
		log.Fatalf("failed service: %v", err)
	}
}
