package game_service

import (
	"github.com/GoReactors/backend-learning/internal/core/port"
	"github.com/codemodus/uidgen"
)

type GameService struct {
	gamesRepository port.GamesRepository
	uidGen          uidgen.UIDGen
}

func New(gamesRepository port.GamesRepository, uidGen uidgen.UIDGen) *GameService {
	return &GameService{
		gamesRepository: gamesRepository,
		uidGen:          uidGen,
	}
}
