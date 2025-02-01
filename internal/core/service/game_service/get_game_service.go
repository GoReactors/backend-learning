package game_service

import (
	"errors"

	"github.com/GoReactors/backend-learning/internal/core/domain"
)

func (srv *GameService) Get(id string) (domain.Game, error) {
	game, err := srv.gamesRepository.Get(id)
	if err != nil {
		return domain.Game{}, errors.New("get game from repository has failed")
	}

	return game, nil
}
