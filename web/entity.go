// entities between client & server.
package web

import (
	"reversi/server"
)

type GameEntity struct {
	ID     int    `json:"id"`
	Board  string `json:"board"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type PlayerEntity struct {
	ID   int    `json:"id"`
	Name string `json:"string"`
}

// Create a new GameEntity from the context.
func (ge *GameEntity) FromContext(gc *server.GameContext) {
	ge.ID = gc.ID
	ge.Board = (<-gc.GetBoard()).Board.ToString()
	ge.Width = 8
	ge.Height = 8
}
