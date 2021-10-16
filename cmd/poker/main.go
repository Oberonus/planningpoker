package main

import (
	"log"
	"planningpoker/internal/domain/games"
	"planningpoker/internal/domain/state"
	"planningpoker/internal/domain/users"
	"planningpoker/internal/infra/eventbus"
	"planningpoker/internal/infra/http"
	"planningpoker/internal/infra/repository"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	gamesRepo := repository.NewMemoryGameRepository()
	usersRepo := repository.NewMemoryUserRepository()
	eventBus := eventbus.NewInternalBus()

	gamesService, err := games.NewService(gamesRepo, eventBus)
	if err != nil {
		log.Fatalf("unable to create games service: %v", err)
	}

	usersService, err := users.NewService(usersRepo, eventBus)
	if err != nil {
		log.Fatalf("unable to create users service: %v", err)
	}

	stateService, err := state.NewService(gamesRepo, usersRepo)
	if err != nil {
		log.Fatalf("unable to create game state service: %v", err)
	}

	api, err := http.NewAPI(gamesService, usersService, stateService)
	if err != nil {
		log.Fatalf("unable to create http API: %v", err)
	}

	fe := http.NewFrontend()

	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	api.SetupRoutes(r)
	fe.SetupRoutes(r)

	if err := r.Run(); err != nil {
		log.Fatalf("failed service: %v", err)
	}
}
