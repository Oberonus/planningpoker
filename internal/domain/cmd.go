package domain

import "time"

type CreateGameCommand struct {
	UserID            string
	Name              string
	TicketURL         string
	CardsDeck         CardsDeck
	EveryoneCanReveal bool
}

func NewCreateGameCommand(name, ticketURL, userID string, deck CardsDeck, everyoneCanReveal bool) (*CreateGameCommand, error) {
	return &CreateGameCommand{
		UserID:            userID,
		Name:              name,
		TicketURL:         ticketURL,
		CardsDeck:         deck,
		EveryoneCanReveal: everyoneCanReveal,
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

type JoinGameCommand struct {
	GameID string
	UserID string
}

func NewJoinGameCommand(gameID, userID string) (*JoinGameCommand, error) {
	return &JoinGameCommand{
		GameID: gameID,
		UserID: userID,
	}, nil
}

type VoteCommand struct {
	GameID string
	UserID string
	Vote   Card
}

func NewVoteCommand(gameID, userID string, card Card) (*VoteCommand, error) {
	return &VoteCommand{
		GameID: gameID,
		UserID: userID,
		Vote:   card,
	}, nil
}

type UnVoteCommand struct {
	GameID string
	UserID string
}

func NewUnVoteCommand(gameID, userID string) (*UnVoteCommand, error) {
	return &UnVoteCommand{
		GameID: gameID,
		UserID: userID,
	}, nil
}

type GameStateCommand struct {
	GameID       string
	UserID       string
	WaitFor      time.Duration
	LastChangeID string
}

func NewGameStateCommand(gameID, userID string, waitFor time.Duration, lastChangeID string) (*GameStateCommand, error) {
	return &GameStateCommand{
		GameID:       gameID,
		UserID:       userID,
		WaitFor:      waitFor,
		LastChangeID: lastChangeID,
	}, nil
}

type RestartGameCommand struct {
	GameID string
	UserID string
}

func NewRestartGameCommand(gameID, userID string) (*RestartGameCommand, error) {
	return &RestartGameCommand{
		GameID: gameID,
		UserID: userID,
	}, nil
}

type RevealCardsCommand struct {
	GameID string
	UserID string
}

func NewRevealCardsCommand(gameID, userID string) (*RevealCardsCommand, error) {
	return &RevealCardsCommand{
		GameID: gameID,
		UserID: userID,
	}, nil
}
