package port

import (
	"context"

	game_domain "github.com/GoReactors/backend-learning/internal/application/game/domain"
)

type GameRepository interface {
	Save(ctx context.Context, game *game_domain.Game) (*game_domain.Game, error)
	FindByID(ctx context.Context, id string) (*game_domain.Game, error)
}
