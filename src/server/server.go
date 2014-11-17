package server

import (
	"board"
	"fmt"
	"sync"
)

var empty struct{}

type Service interface {
	Start() // called at most once.
	Stop()  // called at most once.
}

type ServerContext struct {
	gamesMutex  sync.RWMutex
	games       map[int]*GameContext // accessed by the mutex.
	gameCounter int                  // accessed by the mutex.
	stopRequest chan struct{}
	running     *sync.WaitGroup
}

type GameContext struct {
	ID                 int
	board              *board.Board
	moveRequest        chan MoveRequest
	queryRequest       chan QueryRequest
	waitRequest        chan WaitRequest
	stopRequest        chan struct{}
	queuedWaitRequests []WaitRequest
}

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

type GameResponse struct {
	Success bool        // true if the request succeeded.
	Board   board.Board // board state, can be safely shared.
}

func NewServerContext(running *sync.WaitGroup) *ServerContext {
	return &ServerContext{
		games:       make(map[int]*GameContext),
		stopRequest: make(chan struct{}),
		running:     running,
	}
}

func newGameContext(id int) *GameContext {
	const size = 5
	return &GameContext{
		ID:           id,
		board:        board.NewBoard(),
		queryRequest: make(chan QueryRequest, size),
		moveRequest:  make(chan MoveRequest, size),
		waitRequest:  make(chan WaitRequest, size),
		stopRequest:  make(chan struct{}),
	}
}

func (sc *ServerContext) NewGameContext() *GameContext {
	sc.gamesMutex.Lock()
	defer sc.gamesMutex.Unlock()
	sc.gameCounter++
	gc := newGameContext(sc.gameCounter)
	gc.Start()
	sc.games[sc.gameCounter] = gc
	return gc
}

func (sc *ServerContext) GetGameContext(id int) (*GameContext, bool) {
	sc.gamesMutex.RLock()
	defer sc.gamesMutex.RUnlock()
	game, ok := sc.games[id]
	return game, ok
}

func (gc *GameContext) GetBoard() GameResponse {
	req := QueryRequest{response: make(chan GameResponse, 1)}
	gc.queryRequest <- req
	return <-req.response
}

func (gc *GameContext) Play(player int, row int, col int) GameResponse {
	req := MoveRequest{
		row:      row,
		col:      col,
		player:   player,
		response: make(chan GameResponse, 1),
	}
	gc.moveRequest <- req
	return <-req.response
}

func (gc *GameContext) Start() {
	go func() {
		running := true
		for running {
			// in this loop, you have the full access.
			modified := false
			select {
			case query := <-gc.queryRequest:
				query.response <- GameResponse{Success: true, Board: *gc.board}
			case move := <-gc.moveRequest:
				played := gc.board.Play(move.player, move.row, move.col)
				move.response <- GameResponse{Success: played, Board: *gc.board}
				modified = played
			case wait := <-gc.waitRequest:
				if wait.version <= gc.board.Version() {
					// version is low enough - trigger it now.
					wait.response <- GameResponse{Success: true, Board: *gc.board}
				} else {
					// queue up.
					gc.queuedWaitRequests = append(gc.queuedWaitRequests, wait)
				}
				return
			case <-gc.stopRequest:
				// TODO - when this gets triggered, all the existing messages
				// will not be processed.
				running = false
			}
			if modified {
				curVersion := gc.board.Version()
				trigger := false
				// see if we need to trigger.
				for _, wait := range gc.queuedWaitRequests {
					if wait.version <= curVersion {
						trigger = true
						break
					}
				}
				if trigger {
					// Actually trigger.
					newQueued := make([]WaitRequest, len(gc.queuedWaitRequests))
					for _, wait := range gc.queuedWaitRequests {
						if wait.version <= curVersion {
							wait.response <- GameResponse{Success: true, Board: *gc.board}
						} else {
							newQueued = append(newQueued, wait)
						}
					}
					gc.queuedWaitRequests = newQueued
				}
			}
		}
	}()
}

func (gc *GameContext) Stop() {
	close(gc.stopRequest)
}

func (sc *ServerContext) Start() {
	sc.running.Add(1)
	go func() {
		defer sc.running.Done()
		for {
			select {
			case <-sc.stopRequest:
				return
			}
		}
	}()
}

func (sc *ServerContext) Stop() {
	fmt.Println("Stopping..")
	close(sc.stopRequest)
}
