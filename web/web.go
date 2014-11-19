package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reversi/server"
	"strconv"
	"sync"
)

type GameEntity struct {
	ID    int    `json: id`
	Board string `json: string`
}

// Create a new GameEntity from the context.
func (ge *GameEntity) FromContext(gc *server.GameContext) {
	ge.ID = gc.ID
	ge.Board = gc.GetBoard().Board.ToString()
}

// Create a new HTTPServer.
type HTTPServerContext struct {
	sc      *server.ServerContext
	running *sync.WaitGroup
}

func WrapLogger(handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(rw, r)
	})
}

func (ctx *HTTPServerContext) Start() {
	fmt.Println("Starting...")
	mux := http.NewServeMux()
	mux.HandleFunc("/game/", ctx.GetGameHandler)
	mux.HandleFunc("/game/new", ctx.NewGameHandler)
	mux.HandleFunc("/game/play", ctx.PlayMoveHandler)
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	ctx.sc.Start()
	s := http.Server{
		Addr:    ":8080",
		Handler: WrapLogger(mux),
	}
	go func() {
		log.Fatal(s.ListenAndServe())
	}()
}

func (ctx *HTTPServerContext) Stop() {
	ctx.running.Done()
	ctx.sc.Stop()
}

func (ctx *HTTPServerContext) NewGameHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	gc := ctx.sc.NewGameContext()
	response := GameEntity{}
	response.FromContext(gc)
	bytes, _ := json.Marshal(response)
	rw.Write(bytes)
}

func (ctx *HTTPServerContext) GetGameHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id, err := getIntOrError(rw, r, "id")
	if err != nil {
		return
	}
	gc, ok := ctx.sc.GetGameContext(id)
	if !ok {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	response := GameEntity{}
	response.FromContext(gc)
	bytes, _ := json.Marshal(response)
	rw.Write(bytes)
}

func (ctx *HTTPServerContext) PlayMoveHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id, err := getIntOrError(rw, r, "id")
	if err != nil {
		return
	}
	row, err := getIntOrError(rw, r, "row")
	if err != nil {
		return
	}
	col, err := getIntOrError(rw, r, "col")
	if err != nil {
		return
	}
	player, err := getIntOrError(rw, r, "player")
	if err != nil {
		return
	}
	gc, ok := ctx.sc.GetGameContext(id)
	if !ok {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	// can't play directly!
	response := gc.Play(player, row, col)
	if response.Success {
		rw.Write([]byte("success"))
	} else {
		rw.Write([]byte("failure"))
	}
}

func getIntOrError(rw http.ResponseWriter, r *http.Request, key string) (int, error) {
	value := r.URL.Query().Get(key)
	converted, err := strconv.Atoi(value)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(key))
	}
	return converted, err
}

func Start() {
	// TODO - make them into variables.
	var running sync.WaitGroup
	ctx := HTTPServerContext{
		sc:      server.NewServerContext(&running),
		running: &running,
	}
	ctx.Start()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			ctx.Stop()
		}
	}()
	running.Wait()
}
