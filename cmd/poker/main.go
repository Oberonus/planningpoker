package main

import (
	"log"
	"planningpoker/internal/domain"
	"planningpoker/internal/domain/users"
	"planningpoker/internal/http"
	"planningpoker/internal/repository"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	gamesRepo := repository.NewMemoryGameRepository()
	usersRepo := repository.NewMemoryUserRepository()

	gamesService, err := domain.NewGamesService(gamesRepo, usersRepo)
	if err != nil {
		log.Fatalf("unable to create games service: %v", err)
	}

	usersService, err := users.NewService(usersRepo)
	if err != nil {
		log.Fatalf("unable to create users service: %v", err)
	}

	api, err := http.NewAPI(gamesService, usersService)
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
