package users

type RegisterCommand struct {
	Name string
}

func NewRegisterCommand(name string) (*RegisterCommand, error) {
	return &RegisterCommand{
		Name: name,
	}, nil
}

type UpdateCommand struct {
	ID   string
	Name string
}

func NewUpdateCommand(id, name string) (*UpdateCommand, error) {
	return &UpdateCommand{
		ID:   id,
		Name: name,
	}, nil
}

type AuthByIDCommand struct {
	ID string
}

func NewAuthByIDCommand(id string) (*AuthByIDCommand, error) {
	return &AuthByIDCommand{
		ID: id,
	}, nil
}
