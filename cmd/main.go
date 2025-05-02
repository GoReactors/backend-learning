package main

import (
	"fmt"

	"github.com/GoReactors/backend-learning/config"
	httpadapter "github.com/GoReactors/backend-learning/internal/adapter/http"
	repositoryadapter "github.com/GoReactors/backend-learning/internal/adapter/repository"
	tracingadapter "github.com/GoReactors/backend-learning/internal/adapter/tracing"
	game_service "github.com/GoReactors/backend-learning/internal/application/game/service"
	"github.com/GoReactors/backend-learning/internal/port"
	"github.com/GoReactors/backend-learning/pkg/logger"
	"github.com/GoReactors/backend-learning/pkg/server"
	pkgtracing "github.com/GoReactors/backend-learning/pkg/tracing"
)

func main() {
	// Load ENVs
	cfg := config.LoadConfig()

	// Init Logger
	logger.Initialize(cfg.LoggingConfig)
	defer logger.FirstSyncInMain(logger.Global())
	log := logger.Global().With(logger.Field{Key: "component", Value: "main"})

	// Init OTEL Tracer
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
	// -- Init Handlers
	gameHandler := httpadapter.NewGameHandler(gameService, log, tracingadapter.NewTracer(tracingadapter.GAME_HTTP_HANDLER))

	// -- Init Server
	router := server.NewServer(cfg, log, []port.RouteRegistrar{
		gameHandler,
	})

	// -- Run Server
	router.Run(fmt.Sprintf(":%v", cfg.GinAppPort))
}
