package server

import (
	"github.com/jeeyoungk/reversi/board"
)

type MoveRequest struct {
	row      int
	col      int
	player   int
	response chan GameResponse
}

type QueryRequest struct {
	response chan GameResponse
}

// waits for the given version to occur.
type WaitRequest struct {
	version  int
	response chan GameResponse
}

// Joins the given game.
type JoinRequest struct {
	player PlayerID
	order  int
}

/**
 * Response after any interaction with the game board.
 */
type GameResponse struct {
	Success bool        // true if the request succeeded.
	Board   board.Board // board state, can be safely shared.
}
