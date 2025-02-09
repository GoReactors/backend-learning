package port

import game_domain "github.com/GoReactors/backend-learning/internal/application/game/domain"

type GameAPIPort interface {
	Create(name string) (game_domain.Game, error)
	FindOne(id string) (game_domain.Game, error)
}