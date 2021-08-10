package domain

type CreateGameCommand struct {
	UserID    string
	Name      string
	TicketURL string
}

func NewCreateGameCommand(name string, ticketURL string, userID string) (*CreateGameCommand, error) {
	return &CreateGameCommand{
		UserID:    userID,
		Name:      name,
		TicketURL: ticketURL,
	}, nil
}

type UpdateGameCommand struct {
	GameID    string
	UserID    string
	Name      string
	TicketURL string
}

func NewUpdateGameCommand(id string, name string, ticketURL string, userID string) (*UpdateGameCommand, error) {
	return &UpdateGameCommand{
		GameID:    id,
		UserID:    userID,
		Name:      name,
		TicketURL: ticketURL,
	}, nil
}
