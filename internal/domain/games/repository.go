package games

// GameRepository is a repository contract to fetch/persist games.
type GameRepository interface {
	ModifyExclusively(id string, cb func(game *Game) error) error
	Get(id string) (*Game, error)
	Save(game *Game) error
	GetActiveGamesByPlayerID(playerID string) ([]Game, error)
}
