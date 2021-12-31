package users

// RegisterCommand is a command to register new user.
type RegisterCommand struct {
	Name string
}

// NewRegisterCommand creates a new command instance.
func NewRegisterCommand(name string) (*RegisterCommand, error) {
	return &RegisterCommand{
		Name: name,
	}, nil
}

// UpdateCommand is a command to update user details.
type UpdateCommand struct {
	ID   string
	Name string
}

// NewUpdateCommand creates a new command instance.
func NewUpdateCommand(id, name string) (*UpdateCommand, error) {
	return &UpdateCommand{
		ID:   id,
		Name: name,
	}, nil
}

// AuthByIDCommand is a command to authenticate user by ID.
// In more complex system it will be a separate service, but for current functionality it is enough.
type AuthByIDCommand struct {
	ID string
}

// NewAuthByIDCommand creates a new command instance.
func NewAuthByIDCommand(id string) (*AuthByIDCommand, error) {
	return &AuthByIDCommand{
		ID: id,
	}, nil
}
