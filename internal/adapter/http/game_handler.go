package httpadapter

import (
	"net/http"

	"github.com/gin-gonic/gin"

	game_service "github.com/GoReactors/backend-learning/internal/application/game/service"
	"github.com/GoReactors/backend-learning/internal/port"
	"github.com/GoReactors/backend-learning/pkg/logger"
)

type GameHandler struct {
	svc    *game_service.GameService
	logger logger.Logger
	tracer port.Tracing
}

func NewGameHandler(svc *game_service.GameService, l logger.Logger, tracer port.Tracing) port.RouteRegistrar {
	return &GameHandler{svc: svc, logger: l, tracer: tracer}
}

func (gh *GameHandler) RegisterRoutes(router *gin.RouterGroup) {
	gameGroup := router.Group("/games")
	{
		gameGroup.POST("/", gh.CreateGame)
		gameGroup.GET("/:id", gh.GetGame)
	}
}

type CreateGameRequest struct {
	Title string `json:"title" binding:"required"`
	Mode  string `json:"mode" binding:"required"`
}

func (h *GameHandler) CreateGame(c *gin.Context) {
	ctx := c.Request.Context()
	ctx, span := h.tracer.StartSpan(ctx, "HTTP CreateGame")
	defer span.End()

	var req CreateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request payload", logger.TraceFieldsFromContext(ctx)...)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	game, err := h.svc.CreateGame(ctx, game_service.CreateGameParams{
		Title: req.Title,
		Mode:  req.Mode,
	})
	if err != nil {
		h.logger.Error("failed to create game", logger.TraceFieldsFromContext(ctx)...)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create game"})
		return
	}

	c.JSON(http.StatusCreated, game)
}

func (h *GameHandler) GetGame(c *gin.Context) {
	ctx := c.Request.Context()
	ctx, span := h.tracer.StartSpan(ctx, "HTTP GetGame")
	defer span.End()

	id := c.Param("id")
	if id == "" {
		h.logger.Warn("missing game ID", logger.TraceFieldsFromContext(ctx)...)
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing game ID"})
		return
	}

	game, err := h.svc.GetGame(ctx, id)
	if err != nil {
		h.logger.Warn(
			"game not found",
			append(
				logger.TraceFieldsFromContext(ctx),
				logger.Field{Key: "game_id", Value: id},
			)...,
		)
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}

	c.JSON(http.StatusOK, game)
}
