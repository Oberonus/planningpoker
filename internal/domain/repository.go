package domain

type GameRepository interface {
	ModifyExclusively(id string, cb func(game *Game) error) error
	Get(id string) (*Game, error)
	Save(game *Game) error
}

type UsersRepository interface {
	Get(id string) (*User, error)
	GetMany(ids []string) ([]*User, error)
	Save(user *User) error
}
