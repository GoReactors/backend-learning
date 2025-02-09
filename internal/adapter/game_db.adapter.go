package adapter

import (
	"fmt"

	game_domain "github.com/GoReactors/backend-learning/internal/application/game/domain"
)

type GameRepositoryAdapter struct {
	games map[string]*game_domain.Game
}

func NewGameRepositoryAdapter() *GameRepositoryAdapter {
	return &GameRepositoryAdapter{
		games: make(map[string]*game_domain.Game),
	}
}

func (repo *GameRepositoryAdapter) Get(id string) (game_domain.Game, error) {
	game, exists := repo.games[id]
	if !exists {
		return game_domain.Game{}, fmt.Errorf("game not found")
	}
	return *game, nil
}

func (repo *GameRepositoryAdapter) Save(game *game_domain.Game) error {
	repo.games[game.ID] = game
	return nil
}
