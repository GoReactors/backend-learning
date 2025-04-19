package game_domain

import (
	"github.com/google/uuid"
)

type Game struct {
	ID    string
	Title string
	Mode  string
}

func NewGame(title, mode string) Game {
	return Game{
		ID:    uuid.New().String(),
		Title: title,
		Mode:  mode,
	}
}
