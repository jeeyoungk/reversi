/*
 * Wrapper around the game logic.
 *
 * Includes code for:
 * - managing the concurrency.
 * - managing the player database.
 */
package server

import (
	"fmt"
	"reversi/board"
	"sync"
)

var empty struct{}

type PlayerID int

type Service interface {
	Start() // called at most once.
	Stop()  // called at most once.
}

type ServerContext struct {
	// locks
	gamesMutex     sync.RWMutex
	playersMutex   sync.RWMutex
	games          map[int]*GameContext // accessed by the mutex.
	gameCounter    int                  // accessed by the mutex.
	players        map[PlayerID]*PlayerContext
	playersCounter int
	stopRequest    chan struct{}
	running        *sync.WaitGroup
}

type GameContext struct {
	ID                 int               // game id
	board              *board.Board      // game state
	moveRequest        chan MoveRequest  // queued moves
	queryRequest       chan QueryRequest // queued queries
	waitRequest        chan WaitRequest  // queued waits
	joinRequest        chan JoinRequest  // queued joins
	stopRequest        chan struct{} // queued stops
	queuedWaitRequests []WaitRequest
	players            [2]*PlayerContext // first player name
}

type PlayerContext struct {
	ID   PlayerID
	Name string
}

// constructor
func NewServerContext(running *sync.WaitGroup) *ServerContext {
	return &ServerContext{
		games:       make(map[int]*GameContext),
		players:     make(map[PlayerID]*PlayerContext),
		stopRequest: make(chan struct{}),
		running:     running,
	}
}

func newGameContext(id int) *GameContext {
	const queueSize = 5
	return &GameContext{
		ID:           id,
		board:        board.NewBoard(),
		queryRequest: make(chan QueryRequest, queueSize),
		moveRequest:  make(chan MoveRequest, queueSize),
		waitRequest:  make(chan WaitRequest, queueSize),
		stopRequest:  make(chan struct{}),
	}
}

func newPlayerContext(id PlayerID, name string) *PlayerContext {
	return &PlayerContext{
		ID:   id,
		Name: name,
	}
}

func (sc *ServerContext) NewGameContext() GameContext {
	sc.gamesMutex.Lock()
	defer sc.gamesMutex.Unlock()
	sc.gameCounter++
	gc := newGameContext(sc.gameCounter)
	gc.Start()
	sc.games[sc.gameCounter] = gc
	return *gc
}

func (sc *ServerContext) NewPlayerContext() PlayerContext {
	sc.playersMutex.Lock()
	defer sc.playersMutex.Unlock()
	sc.playersCounter++
	pid := PlayerID(sc.playersCounter)
	name := fmt.Sprintf("player-%d", sc.playersCounter)
	pc := newPlayerContext(pid, name)
	sc.players[pid] = pc
	return *pc
}

func (sc *ServerContext) GetGameCount() int {
	sc.gamesMutex.RLock()
	defer sc.gamesMutex.RUnlock()
	return len(sc.games)
}

func (sc *ServerContext) GetPlayerCount() int {
	sc.playersMutex.RLock()
	defer sc.playersMutex.RUnlock()
	return len(sc.players)
}

func (sc *ServerContext) GetGameContext(id int) (*GameContext, bool) {
	sc.gamesMutex.RLock()
	defer sc.gamesMutex.RUnlock()
	game, ok := sc.games[id]
	return game, ok
}

func (sc *ServerContext) GetPlayerContext(id PlayerID) (*PlayerContext, bool) {
	sc.playersMutex.RLock()
	defer sc.playersMutex.RUnlock()
	player, ok := sc.players[id]
	return player, ok
}

func (gc *GameContext) GetBoard() chan GameResponse {
	req := QueryRequest{response: responseChannel()}
	gc.queryRequest <- req
	return req.response
}

func (gc *GameContext) Wait(version int) chan GameResponse {
	req := WaitRequest{
		version:  version,
		response: responseChannel(),
	}
	gc.waitRequest <- req
	return req.response
}

func (gc *GameContext) Play(player int, row int, col int) chan GameResponse {
	req := MoveRequest{
		row:      row,
		col:      col,
		player:   player,
		response: responseChannel(),
	}
	gc.moveRequest <- req
	return req.response
}

func (gc *GameContext) Join(pc PlayerContext) {

}

func (gc *GameContext) Start() {
	go gc.loop()
}

func (gc *GameContext) loop() {
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
				// version is low enough - return the response now.
				wait.response <- GameResponse{Success: true, Board: *gc.board}
			} else {
				// queue up.
				gc.queuedWaitRequests = append(gc.queuedWaitRequests, wait)
			}
		case <-gc.stopRequest:
			// TODO - when this gets triggered, all the existing messages
			// will not be processed.
			running = false
		}
		if modified {
			gc.triggerWatch(gc.board.Version())
		}
	}
}

func (gc *GameContext) triggerWatch(curVersion int) {
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

// utility functions
func responseChannel() chan GameResponse {
	return make(chan GameResponse, 1)
}
