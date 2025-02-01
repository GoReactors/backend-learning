package game_service

import (
	"errors"

	"github.com/GoReactors/backend-learning/internal/core/domain"
)

func (srv *GameService) Create(name string, size uint, bombs uint) (domain.Game, error) {
	if bombs >= size*size {
		return domain.Game{}, errors.New("the number of bombs is invalid")
	}

	game := domain.NewGame(srv.uidGen.UID().String(), name, size, bombs)

	if err := srv.gamesRepository.Save(game); err != nil {
		return domain.Game{}, errors.New("create game into repository has failed")
	}

	return game, nil
}
