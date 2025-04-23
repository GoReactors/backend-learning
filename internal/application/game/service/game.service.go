package game_service

import (
	"context"

	game_domain "github.com/GoReactors/backend-learning/internal/application/game/domain"
	"github.com/GoReactors/backend-learning/internal/port"
	tracingport "github.com/GoReactors/backend-learning/internal/port/tracer"
	"github.com/GoReactors/backend-learning/pkg/logger"
	"go.opentelemetry.io/otel/codes"
)

type GameService struct {
	repo   port.GameRepository
	logger logger.Logger
	tracer tracingport.Tracer
}

type CreateGameParams struct {
	Title string
	Mode  string
}

func NewGameService(repo port.GameRepository, log logger.Logger, tracer tracingport.Tracer) *GameService {
	return &GameService{
		repo:   repo,
		logger: log,
		tracer: tracer,
	}
}

func (s *GameService) CreateGame(ctx context.Context, params CreateGameParams) (*game_domain.Game, error) {
	ctx, span := s.tracer.StartSpan(ctx, "GameService.CreateGame")
	defer span.End()

	s.tracer.SetAttribute(span, "game.title", params.Title)
	s.tracer.SetAttribute(span, "game.mode", params.Mode)

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
	ctx, span := s.tracer.StartSpan(ctx, "GameService.GetGame")
	defer span.End()

	s.tracer.SetAttribute(span, "game.id", id)

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
