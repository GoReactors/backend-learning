package repositoryadapter

import (
	"context"
	"errors"
	"sync"

	game_domain "github.com/GoReactors/backend-learning/internal/application/game/domain"
	"github.com/GoReactors/backend-learning/internal/port"
	"go.opentelemetry.io/otel"
)

type InMemoryGameRepository struct {
	mu    sync.RWMutex
	store map[string]*game_domain.Game
}

func NewInMemoryGameRepository() port.GameRepository {
	return &InMemoryGameRepository{
		store: make(map[string]*game_domain.Game),
	}
}

// Save adds or updates a game in the store
func (r *InMemoryGameRepository) Save(ctx context.Context, game *game_domain.Game) (*game_domain.Game, error) {
	ctx, span := otel.Tracer("game-repo").Start(ctx, "InMemoryGameRepository.Save")
	defer span.End()

	r.mu.Lock()
	defer r.mu.Unlock()

	r.store[game.ID] = game
	return game, nil
}

// FindByID fetches a game by ID
func (r *InMemoryGameRepository) FindByID(ctx context.Context, id string) (*game_domain.Game, error) {
	ctx, span := otel.Tracer("game-repo").Start(ctx, "InMemoryGameRepository.FindByID")
	defer span.End()

	r.mu.RLock()
	defer r.mu.RUnlock()

	game, ok := r.store[id]
	if !ok {
		return nil, errors.New("game not found")
	}

	return game, nil
}
