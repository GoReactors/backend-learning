package main

import (
	"fmt"

	"github.com/GoReactors/backend-learning/config"
	httpadapter "github.com/GoReactors/backend-learning/internal/adapter/http"
	repositoryadapter "github.com/GoReactors/backend-learning/internal/adapter/repository"
	"github.com/GoReactors/backend-learning/internal/adapter/tracing"
	game_service "github.com/GoReactors/backend-learning/internal/application/game/service"
	"github.com/GoReactors/backend-learning/pkg/logger"
)

func main() {
	// Load ENVs
	cfg := config.LoadConfig()

	// Init Logger
	logger.Initialize(cfg.LoggingConfig)
	defer logger.FirstSyncInMain(logger.Global())
	log := logger.Global().With(logger.Field{Key: "component", Value: "main"})

	// Init Tracer
	tracing.InitTracer(cfg, logger.Global().With(logger.Field{Key: "component", Value: "tracing"}))

	// Init Services
	gameRepository := repositoryadapter.NewInMemoryGameRepository()
	gameService := game_service.NewGameService(
		gameRepository,
		logger.Global().With(logger.Field{
			Key:   "component",
			Value: "game_service",
		}),
	)

	// start HTTP server
	router := httpadapter.NewServer(cfg, gameService, log)
	router.Run(fmt.Sprintf(":%v", cfg.GinAppPort))
}
