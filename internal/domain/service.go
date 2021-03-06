package domain

import "errors"

type GamesService struct {
	gamesRepo GameRepository
	usersRepo UsersRepository
}

func NewGamesService(gr GameRepository, ur UsersRepository) (*GamesService, error) {
	if gr == nil {
		return nil, errors.New("games repository should be provided")
	}
	if ur == nil {
		return nil, errors.New("users repository should be provided")
	}

	return &GamesService{
		gamesRepo: gr,
		usersRepo: ur,
	}, nil
}

func (s *GamesService) Create(userID string) (string, error) {
	game := NewTShirtGame()
	if err := game.Join(userID); err != nil {
		return "", err
	}
	if err := s.gamesRepo.Save(game); err != nil {
		return "", err
	}

	return game.ID, nil
}

func (s *GamesService) Join(gameID, userID string) error {
	return s.gamesRepo.ModifyExclusively(gameID, func(game *Game) error {
		return game.Join(userID)
	})
}

func (s *GamesService) Restart(gameID, userID string) error {
	return s.gamesRepo.ModifyExclusively(gameID, func(game *Game) error {
		return game.Restart(userID)
	})
}

func (s *GamesService) Vote(gameID, userID, vote string) error {
	return s.gamesRepo.ModifyExclusively(gameID, func(game *Game) error {
		return game.Vote(userID, vote)
	})
}

func (s *GamesService) UnVote(gameID, userID string) error {
	return s.gamesRepo.ModifyExclusively(gameID, func(game *Game) error {
		return game.UnVote(userID)
	})
}

func (s *GamesService) Leave(gameID, userID, vote string) error {
	return s.gamesRepo.ModifyExclusively(gameID, func(game *Game) error {
		return game.Leave(userID)
	})
}

func (s *GamesService) Reveal(gameID, userID string) error {
	return s.gamesRepo.ModifyExclusively(gameID, func(game *Game) error {
		return game.Reveal(userID)
	})
}

func (s *GamesService) GameState(gameID, userID string) (*GameState, error) {
	var gotGame *Game
	err := s.gamesRepo.ModifyExclusively(gameID, func(game *Game) error {
		gotGame = game
		return game.Ping(userID)
	})
	if err != nil {
		return nil, err
	}

	userIDs := make([]string, 0)
	for id := range gotGame.Players {
		userIDs = append(userIDs, id)
	}

	users, err := s.usersRepo.GetMany(userIDs)
	if err != nil {
		return nil, err
	}

	mapUsers := make(map[string]*User)
	for _, u := range users {
		mapUsers[u.ID] = u
	}

	return gotGame.CurrentState(userID, mapUsers), nil
}
