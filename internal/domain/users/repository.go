package users

type Repository interface {
	Get(id string) (*User, error)
	Save(user User) error
}
