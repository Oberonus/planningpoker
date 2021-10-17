package state

import (
	"planningpoker/internal/domain/games"
	"planningpoker/internal/domain/users"
)

type GameRepository interface {
	Get(id string) (*games.Game, error)
}

type UsersRepository interface {
	GetMany(ids []string) ([]users.User, error)
}
