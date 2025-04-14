package game_service

import (
	game_domain "github.com/GoReactors/backend-learning/internal/application/game/domain"
	"github.com/GoReactors/backend-learning/internal/port"
	"github.com/GoReactors/backend-learning/pkg/logger"
)

type GameService struct {
	repo   port.GameRepositoryPort
	logger logger.Logger
}

func NewGameService(repo port.GameRepositoryPort, logger logger.Logger) *GameService {
	return &GameService{
		repo:   repo,
		logger: logger,
	}
}

func (s *GameService) Create(name string) (game_domain.Game, error) {
	game := game_domain.NewGame(name)
	err := s.repo.Save(&game)
	if err != nil {
		return game_domain.Game{}, err
	}
	s.logger.Info("Game created", logger.Field{Key: "name", Value: name})
	return game, nil
}

func (s *GameService) FindOne(id string) (game_domain.Game, error) {
	game, err := s.repo.Get(id)
	if err != nil {
		return game_domain.Game{}, err
	}
	return game, nil
}
