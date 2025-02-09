package game_domain

import (
	"github.com/google/uuid"
)

type Game struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewGame(name string) Game {
	return Game{
		ID:   uuid.New().String(),
		Name: name,
	}
}
