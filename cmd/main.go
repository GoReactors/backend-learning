package main

import (
	"fmt"

	"github.com/GoReactors/backend-learning/config"
	httpadapter "github.com/GoReactors/backend-learning/internal/adapter/http"
	repositoryadapter "github.com/GoReactors/backend-learning/internal/adapter/repository"
	tracingadapter "github.com/GoReactors/backend-learning/internal/adapter/tracing"
	game_service "github.com/GoReactors/backend-learning/internal/application/game/service"
	"github.com/GoReactors/backend-learning/pkg/logger"
	pkgtracing "github.com/GoReactors/backend-learning/pkg/tracing"
)

func main() {
	// Load ENVs
	cfg := config.LoadConfig()

	// Init Logger
	logger.Initialize(cfg.LoggingConfig)
	defer logger.FirstSyncInMain(logger.Global())
	log := logger.Global().With(logger.Field{Key: "component", Value: "main"})

	// Init Tracer
	pkgtracing.InitOTELTracer(cfg, logger.Global().With(logger.Field{Key: "component", Value: "tracing"}))

	// Init Services
	gameRepository := repositoryadapter.NewInMemoryGameRepository(tracingadapter.NewTracer(tracingadapter.GAME_REPO))
	gameService := game_service.NewGameService(
		gameRepository,
		logger.Global().With(logger.Field{
			Key:   "component",
			Value: "game_service",
		}),
		tracingadapter.NewTracer(tracingadapter.GAME_SERVICE),
	)

	// start HTTP server
	router := httpadapter.NewServer(cfg, gameService, log)
	router.Run(fmt.Sprintf(":%v", cfg.GinAppPort))
}
