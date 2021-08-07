package domain

type PlayerState struct {
	Name      string
	VotedCard string
}

type GameState struct {
	ChangeID  string
	Cards     []string
	Players   []PlayerState
	State     string
	VotedCard string
	CanReveal bool
}
