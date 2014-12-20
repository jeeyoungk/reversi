package server

import (
	"sync"
	"testing"
)

// fixtures
func newServerContext() *ServerContext {
	running := &sync.WaitGroup{}
	sc := NewServerContext(running)
	sc.Start()
	return sc
}

func Test_NewGame(t *testing.T) {
	sc := newServerContext()
	defer sc.Stop()

	game1 := sc.NewGameContext()
	if game1.ID != 1 {
		t.Fail()
	}
	game2 := sc.NewGameContext()
	if game2.ID != 2 {
		t.Fail()
	}
}

func Test_GetBoard(t *testing.T) {
	sc := newServerContext()
	defer sc.Stop()

	game := sc.NewGameContext()
	resp := <-game.GetBoard()

	if !resp.Success {
		t.Error("GetBoard failed")
	}
	if resp.Board.Version() != 0 {
		t.Error("Different Version")
	}
}

func Test_Play(t *testing.T) {
	sc := newServerContext()
	defer sc.Stop()

	game := sc.NewGameContext()
	resp := <-game.Play(1, 2, 4)

	if !resp.Success {
		t.Error("GetBoard failed")
	}
	if resp.Board.Version() != 1 {
		t.Error("Different Version")
	}
}

func Test_Wait(t *testing.T) {
	sc := newServerContext()
	defer sc.Stop()

	game := sc.NewGameContext()
	if resp := <-game.Wait(0); !resp.Success {
		t.Error("Wait have failed.")
	}
	future := game.Wait(1)
	if resp := <-game.Play(1, 2, 4); !resp.Success {
		t.Error("GetBoard failed")
	}
	if resp := <-future; !resp.Success {
		t.Error("Future wait have failed.")
	}
}
