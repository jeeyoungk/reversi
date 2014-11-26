package web

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reversi/server"
	"strconv"
	"sync"
)

// Create a new GameEntity from the context.
func (ge *GameEntity) FromContext(gc *server.GameContext) {
	ge.ID = gc.ID
	ge.Board = (<-gc.GetBoard()).Board.ToString()
}

// Create a new HTTPServer.
type HTTPServerContext struct {
	ServerContext *server.ServerContext
	running       *sync.WaitGroup
}

func WrapLogger(handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(rw, r)
	})
}

func (ctx *HTTPServerContext) Start() {
	log.Println("Starting...")
	mux := http.NewServeMux()
	mux.HandleFunc("/game/", ctx.GetGameHandler)
	mux.HandleFunc("/game/new", ctx.NewGameHandler)
	mux.HandleFunc("/game/play", ctx.PlayMoveHandler)
	mux.HandleFunc("/game/websocket", ctx.WebsocketHandler)
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	ctx.ServerContext.Start()
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
	ctx.ServerContext.Stop()
}

func (ctx *HTTPServerContext) NewGameHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	gc := ctx.ServerContext.NewGameCoontext()
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
	gc, ok := ctx.ServerContext.GetGameContext(id)
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
	gc, ok := ctx.ServerContext.GetGameContext(id)
	if !ok {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	// can't play directly!
	response := <-gc.Play(player, row, col)
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
		ServerContext: server.NewServerContext(&running),
		running:       &running,
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
