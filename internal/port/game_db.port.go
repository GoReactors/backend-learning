package port

import game_domain "github.com/GoReactors/backend-learning/internal/application/game/domain"

type GameRepositoryPort interface {
	Get(id string) (game_domain.Game, error)
	Save(*game_domain.Game) error
}
