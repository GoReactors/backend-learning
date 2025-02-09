package main

import (
	"log"

	"github.com/GoReactors/backend-learning/config"
	"github.com/GoReactors/backend-learning/internal/adapter"
	game_service "github.com/GoReactors/backend-learning/internal/application/game/service"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.LoadConfig()

	gameRepository := adapter.NewGameRepositoryAdapter()
	gameService := game_service.NewGameService(gameRepository)
	gameAPIAdapter := adapter.NewGameAPIAdapter(gameService)
	gameAPIAdapter.Run(cfg)
}
