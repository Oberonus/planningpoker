package domain

type PlayerState struct {
	Name      string
	VotedCard string
}

type GameState struct {
	Cards     []string
	Players   []PlayerState
	State     string
	VotedCard string
	CanReveal bool
}
