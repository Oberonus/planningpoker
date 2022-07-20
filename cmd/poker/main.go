package main

import (
	"log"

	"planningpoker/internal/domain/state"

	"github.com/sirupsen/logrus"

	"planningpoker/internal/domain/games"
	"planningpoker/internal/domain/users"
	"planningpoker/internal/infra/async"
	"planningpoker/internal/infra/auth"
	"planningpoker/internal/infra/eventbus"
	"planningpoker/internal/infra/http"
	"planningpoker/internal/infra/repository"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	logrus.Infof("starting the service")

	eventBus := eventbus.NewInternalBus()

	gamesRepo := repository.NewMemoryGameRepository(eventBus)
	usersRepo := repository.NewMemoryUserRepository(eventBus)

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
