package main

import (
	"github.com/GoReactors/backend-learning/config"
	"github.com/GoReactors/backend-learning/internal/adapter"
	game_service "github.com/GoReactors/backend-learning/internal/application/game/service"
	"github.com/GoReactors/backend-learning/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file" + err.Error())
	}

	cfg := config.LoadConfig()

	logger.Initialize(cfg.LoggingConfig)
	defer func() {
		if err := logger.Global().Sync(); err != nil {
			// Don't panic if sync fails (common when writing to stdout)
			logger.Global().Error("failed to sync logger", logger.Field{Key: "error", Value: err})
		}
	}()
	log := logger.Global().With(logger.Field{Key: "component", Value: "main"})

	gameRepository := adapter.NewGameRepositoryAdapter()
	gameService := game_service.NewGameService(
		gameRepository,
		logger.Global().With(logger.Field{
			Key:   "component",
			Value: "game_service",
		}),
	)
	gameAPIAdapter := adapter.NewGameAPIAdapter(gameService)

	log.Info("starting application initialization")
	gameAPIAdapter.Run(cfg)
}
