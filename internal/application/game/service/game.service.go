package game_service

import (
	"context"

	game_domain "github.com/GoReactors/backend-learning/internal/application/game/domain"
	"github.com/GoReactors/backend-learning/internal/port"
	"github.com/GoReactors/backend-learning/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type GameService struct {
	repo   port.GameRepository
	logger logger.Logger
}

type CreateGameParams struct {
	Title string
	Mode  string
}

func NewGameService(repo port.GameRepository, log logger.Logger) *GameService {
	return &GameService{
		repo:   repo,
		logger: log,
	}
}

func (s *GameService) CreateGame(ctx context.Context, params CreateGameParams) (*game_domain.Game, error) {
	ctx, span := otel.Tracer("game-service").Start(ctx, "GameService.CreateGame")
	defer span.End()

	span.SetAttributes(
		attribute.String("game.title", params.Title),
		attribute.String("game.mode", params.Mode),
	)

	game := game_domain.NewGame(params.Title, params.Mode)
	saved, err := s.repo.Save(ctx, &game)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to save game")
		s.logger.Error("failed to save game", logger.TraceFieldsFromContext(ctx)...)
		return nil, err
	}

	span.SetStatus(codes.Ok, "game created")
	s.logger.InfoCtx(ctx, "game created", logger.Field{Key: "game_id", Value: saved.ID})

	return saved, nil
}

func (s *GameService) GetGame(ctx context.Context, id string) (*game_domain.Game, error) {
	ctx, span := otel.Tracer("game-service").Start(ctx, "GameService.GetGame")
	defer span.End()

	span.SetAttributes(attribute.String("game.id", id))

	game, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "game not found")
		s.logger.Warn("game not found", logger.Field{Key: "game_id", Value: id})
		return nil, err
	}

	span.SetStatus(codes.Ok, "game retrieved")
	s.logger.InfoCtx(ctx, "game retrieved", logger.Field{Key: "game_id", Value: id})

	return game, nil
}
