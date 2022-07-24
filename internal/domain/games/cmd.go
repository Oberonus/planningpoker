package games

// CreateGameCommand is a game creation command.
type CreateGameCommand struct {
	UserID            string
	Name              string
	TicketURL         string
	CardsDeck         CardsDeck
	EveryoneCanReveal bool
}

// NewCreateGameCommand creates a new command instance.
func NewCreateGameCommand(name, ticketURL, userID string, deck CardsDeck, everyoneCanReveal bool) (*CreateGameCommand, error) {
	return &CreateGameCommand{
		UserID:            userID,
		Name:              name,
		TicketURL:         ticketURL,
		CardsDeck:         deck,
		EveryoneCanReveal: everyoneCanReveal,
	}, nil
}

// UpdateGameCommand is a game update command.
type UpdateGameCommand struct {
	GameID    string
	UserID    string
	Name      string
	TicketURL string
}

// NewUpdateGameCommand creates a new command instance.
func NewUpdateGameCommand(id string, name string, ticketURL string, userID string) (*UpdateGameCommand, error) {
	return &UpdateGameCommand{
		GameID:    id,
		UserID:    userID,
		Name:      name,
		TicketURL: ticketURL,
	}, nil
}

// JoinGameCommand is a join game command.
type JoinGameCommand struct {
	GameID string
	UserID string
}

// NewJoinGameCommand creates a new command instance.
func NewJoinGameCommand(gameID, userID string) (*JoinGameCommand, error) {
	return &JoinGameCommand{
		GameID: gameID,
		UserID: userID,
	}, nil
}

// VoteCommand is a user voting command.
type VoteCommand struct {
	GameID     string
	UserID     string
	Vote       Card
	Confidence string
}

// NewVoteCommand creates a new command instance.
func NewVoteCommand(gameID, userID string, card Card, confidence string) (*VoteCommand, error) {
	return &VoteCommand{
		GameID:     gameID,
		UserID:     userID,
		Vote:       card,
		Confidence: confidence,
	}, nil
}

// UnVoteCommand is a user un-voting command.
type UnVoteCommand struct {
	GameID string
	UserID string
}

// NewUnVoteCommand creates a new command instance.
func NewUnVoteCommand(gameID, userID string) (*UnVoteCommand, error) {
	return &UnVoteCommand{
		GameID: gameID,
		UserID: userID,
	}, nil
}

// RestartGameCommand is a game restart command.
type RestartGameCommand struct {
	GameID string
	UserID string
}

// NewRestartGameCommand creates a new command instance.
func NewRestartGameCommand(gameID, userID string) (*RestartGameCommand, error) {
	return &RestartGameCommand{
		GameID: gameID,
		UserID: userID,
	}, nil
}

// RevealCardsCommand is a stop-game-and-reveal-cards command.
type RevealCardsCommand struct {
	GameID string
	UserID string
}

// NewRevealCardsCommand creates a new command instance.
func NewRevealCardsCommand(gameID, userID string) (*RevealCardsCommand, error) {
	return &RevealCardsCommand{
		GameID: gameID,
		UserID: userID,
	}, nil
}

// LeaveGameCommand is a leave game command.
type LeaveGameCommand struct {
	GameID string
	UserID string
}

// NewLeaveGameCommand forces a player to leave the game.
func NewLeaveGameCommand(gameID, userID string) (*LeaveGameCommand, error) {
	return &LeaveGameCommand{
		GameID: gameID,
		UserID: userID,
	}, nil
}
