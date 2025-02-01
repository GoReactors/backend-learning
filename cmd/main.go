package main

import (
	http_handler "github.com/GoReactors/backend-learning/internal/adapter/handler/http"
	"github.com/GoReactors/backend-learning/internal/adapter/repository/game_repository"
	"github.com/GoReactors/backend-learning/internal/core/service/game_service"
	"github.com/codemodus/uidgen"
	"github.com/gin-gonic/gin"
)

func main() {
	gamesRepository := game_repository.NewMemKVS()
	gamesService := game_service.New(gamesRepository, *uidgen.New(0, uidgen.VARCHAR26))
	gamesHandler := http_handler.NewHTTPHandler(gamesService)

	router := gin.New()
	router.GET("/games/:id", gamesHandler.Get)
	// router.POST("/games", gamesHandler.Create)

	router.Run(":8080")
}
