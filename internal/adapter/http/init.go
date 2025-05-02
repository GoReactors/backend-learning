package httpadapter

import (
	"github.com/GoReactors/backend-learning/config"
	tracingadapter "github.com/GoReactors/backend-learning/internal/adapter/tracing"
	game_service "github.com/GoReactors/backend-learning/internal/application/game/service"
	"github.com/GoReactors/backend-learning/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewServer(config config.Config, gameService *game_service.GameService, l logger.Logger) *gin.Engine {
	router := gin.New()

	// Middlewares
	router.Use(gin.Recovery())
	router.Use(otelgin.Middleware(config.ServiceName))
	router.Use(logger.GinZapMiddleware(l, config.ServiceName, config.Environment))

	// Register Routes
	RegisterRoutesGame(router, gameService, l, tracingadapter.NewTracer(tracingadapter.GAME_HTTP_HANDLER))

	return router
}
