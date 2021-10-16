package state

import (
	"planningpoker/internal/domain/games"
	"planningpoker/internal/domain/users"
)

type GameRepository interface {
	ModifyExclusively(id string, cb func(game *games.Game) error) error
	Get(id string) (*games.Game, error)
}

type UsersRepository interface {
	GetMany(ids []string) ([]users.User, error)
}
