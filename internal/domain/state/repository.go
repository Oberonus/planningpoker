package state

import (
	"planningpoker/internal/domain/games"
	"planningpoker/internal/domain/users"
)

// GameRepository is a contract to fetch games data.
type GameRepository interface {
	Get(id string) (*games.Game, error)
}

// UsersRepository is a contract to fetch users data.
type UsersRepository interface {
	GetMany(ids []string) ([]users.User, error)
}
