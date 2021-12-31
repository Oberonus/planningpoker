package users

// Repository is a contract to fetch and update users.
type Repository interface {
	Get(id string) (*User, error)
	Save(user User) error
}
