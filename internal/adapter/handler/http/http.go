package http_handler

import (
	"github.com/GoReactors/backend-learning/internal/core/port"
	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	gamesService port.GamesService
}

func NewHTTPHandler(gamesService port.GamesService) *HTTPHandler {
	return &HTTPHandler{
		gamesService: gamesService,
	}
}

func (hdl *HTTPHandler) Get(c *gin.Context) {
	game, err := hdl.gamesService.Get(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, game)
}
