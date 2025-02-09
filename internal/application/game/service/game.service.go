package game_service

import (
	game_domain "github.com/GoReactors/backend-learning/internal/application/game/domain"
	"github.com/GoReactors/backend-learning/internal/port"
)

type GameService struct {
	repo port.GameRepositoryPort
}

func NewGameService(repo port.GameRepositoryPort) *GameService {
	return &GameService{
		repo: repo,
	}
}

func (s *GameService) Create(name string) (game_domain.Game, error) {
	game := game_domain.NewGame(name)
	err := s.repo.Save(&game)
	if err != nil {
		return game_domain.Game{}, err
	}
	return game, nil
}

func (s *GameService) FindOne(id string) (game_domain.Game, error) {
	game, err := s.repo.Get(id)
	if err != nil {
		return game_domain.Game{}, err
	}
	return game, nil
}
