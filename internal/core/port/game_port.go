package port

import "github.com/GoReactors/backend-learning/internal/core/domain"

type GamesRepository interface {
	Get(id string) (domain.Game, error)
	Save(domain.Game) error
}

type GamesService interface {
	Get(id string) (domain.Game, error)
	Create(name string, size uint, bombs uint) (domain.Game, error)
}
