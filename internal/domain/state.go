package domain

type PlayerState struct {
	Name      string
	VotedCard string
}

type GameState struct {
	Name      string
	TicketURL string
	ChangeID  string
	Cards     []string
	Players   []PlayerState
	State     string
	VotedCard string
	CanReveal bool
}
