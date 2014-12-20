// entities between client & server.
// these objects are serialized between the server & the client.
package web

import (
	"github.com/jeeyoungk/reversi/server"
)

type GameEntity struct {
	ID     int    `json:"id"`
	Board  string `json:"board"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type PlayerEntity struct {
	ID   server.PlayerID `json:"id"`
	Name string          `json:"name"`
}

// Create a new GameEntity from the cotext.
func (ge *GameEntity) FromContext(gc server.GameContext) {
	ge.ID = gc.ID
	ge.Board = (<-gc.GetBoard()).Board.ToString()
	ge.Width = 8
	ge.Height = 8
}

func (pe *PlayerEntity) FromContext(pc server.PlayerContext) {
	pe.ID = pc.ID
	pe.Name = pc.Name
}
