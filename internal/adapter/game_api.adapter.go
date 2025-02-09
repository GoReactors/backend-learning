package adapter

import (
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/GoReactors/backend-learning/config"
	"github.com/GoReactors/backend-learning/internal/port"
)

const gamesPrefix = "/games"

type GameAPIAdapter struct {
	app  *gin.Engine
	port port.GameAPIPort
}

func NewGameAPIAdapter(port port.GameAPIPort) *GameAPIAdapter {
	adapter := &GameAPIAdapter{
		app:  gin.Default(),
		port: port,
	}

	adapter.app.GET(gamesPrefix+"/:id", adapter.findOne)
	adapter.app.POST(gamesPrefix, adapter.create)

	return adapter
}

func (adapter *GameAPIAdapter) Run(cfg config.Config) {
	adapter.app.Run(":" + strconv.Itoa(cfg.GinAppPort))
}

func (adapter *GameAPIAdapter) findOne(c *gin.Context) {
	id := c.Param("id")
	game, err := adapter.port.FindOne(id)
	if (err != nil) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, game)
}

func (adapter *GameAPIAdapter) create(c *gin.Context) {
	var request struct {
		Name string `json:"name"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var gameName string = request.Name
	createdGame, err := adapter.port.Create(gameName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdGame)
}
