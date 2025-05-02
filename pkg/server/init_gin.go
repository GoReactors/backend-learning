package server

import (
	"github.com/GoReactors/backend-learning/config"
	"github.com/GoReactors/backend-learning/internal/port"
	"github.com/GoReactors/backend-learning/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewServer(config config.Config, l logger.Logger, routeRegistrars []port.RouteRegistrar) *gin.Engine {
	router := gin.New()

	// Middlewares
	router.Use(gin.Recovery())
	router.Use(otelgin.Middleware(config.ServiceName))
	router.Use(logger.GinZapMiddleware(l, config.ServiceName, config.Environment))

	// Register Routes
	api := router.Group("/api")
	for _, registrar := range routeRegistrars {
		registrar.RegisterRoutes(api)
	}

	return router
}
