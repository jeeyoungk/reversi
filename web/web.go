package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"

	"github.com/gorilla/schema"
	"reversi/server"
)

var decoder = schema.NewDecoder()

// Create a new GameEntity from the context.
func (ge *GameEntity) FromContext(gc *server.GameContext) {
	ge.ID = gc.ID
	ge.Board = (<-gc.GetBoard()).Board.ToString()
	ge.Width = 8
	ge.Height = 8
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
	port := 8080
	log.Printf("Starting. Port:%d\n", port)
	mux := http.NewServeMux()
	mux.HandleFunc("/game/", ctx.GetGameHandler)
	// create a new game.
	mux.HandleFunc("/player/new", ctx.NewPlayerHandler)
	mux.HandleFunc("/game/new", ctx.NewGameHandler)
	// join a game.
	mux.HandleFunc("/game/join", ctx.JoinGameHandler)
	// play a new move.
	mux.HandleFunc("/game/play", ctx.PlayMoveHandler)
	mux.HandleFunc("/_admin/state", ctx.GameStateHandler)
	mux.HandleFunc("/websocket", ctx.WebsocketHandler)
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	ctx.ServerContext.Start()
	s := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
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

func (ctx *HTTPServerContext) GameStateHandler(rw http.ResponseWriter, r *http.Request) {
	type Response struct {
		Games int
	}
	response := Response{
		Games: ctx.ServerContext.GetGameCount(),
	}
	bytes, _ := json.Marshal(response)
	rw.Write(bytes)
}

func (ctx *HTTPServerContext) NewPlayerHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	pc := ctx.ServerContext.NewPlayerContext()
	response := PlayerEntity{}
	bytes, _ := json.Marshal(response)
	rw.Write(bytes)

}
func (ctx *HTTPServerContext) NewGameHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	gc := ctx.ServerContext.NewGameContext()
	response := GameEntity{}
	response.FromContext(gc)
	bytes, _ := json.Marshal(response)
	rw.Write(bytes)
}

func (ctx *HTTPServerContext) JoinGameHandler(rw http.ResponseWriter, r *http.Request) {
	type RequestForm struct {
		Id int `schema:id`
	}
	type Response struct {
	}
	if r.Method != "POST" {
		writeError(rw, http.StatusMethodNotAllowed, nil, "Invalid method.")
		return
	}
	if err := r.ParseForm(); err != nil {
		writeError(rw, http.StatusMethodNotAllowed, err, "Error while parsing form.")
		return
	}
	reqForm := &RequestForm{}
	if err := decoder.Decode(reqForm, r.PostForm); err != nil {
		writeError(rw, http.StatusBadRequest, err, "Error parsing form.")
		return
	}
	if reqForm.Id != 0 {
		// join existing game, as a player.
	}
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
	type PostForm struct {
		Id     int `schema:id`
		Row    int `schema:row`
		Col    int `schema:col`
		Player int `schema:player`
	}
	if r.Method != "POST" {
		writeError(rw, http.StatusMethodNotAllowed, nil, "Invalid method.")
		return
	}
	if err := r.ParseForm(); err != nil {
		writeError(rw, http.StatusMethodNotAllowed, err, "Error while parsing form.")
		return
	}

	postForm := &PostForm{}
	if err := decoder.Decode(postForm, r.PostForm); err != nil {
		writeError(rw, http.StatusBadRequest, err, "Error parsing form.")
		return
	}
	gc, ok := ctx.ServerContext.GetGameContext(postForm.Id)
	if !ok {
		writeError(rw, http.StatusBadRequest, nil, "Invalid Game.")
		return
	}

	response := <-gc.Play(postForm.Player, postForm.Row, postForm.Col)
	if response.Success {
		rw.Write([]byte("success"))
	} else {
		rw.Write([]byte("failure"))
	}
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

func writeError(rw http.ResponseWriter, status int, err error, message string) {
	type ErrorResponse struct {
		Error   string `json:"debug"`
		Message string `json:"message"`
	}
	rw.WriteHeader(status)
	errorString := ""
	if err != nil {
		errorString = err.Error()
	}
	bytes, _ := json.Marshal(ErrorResponse{
		Error:   errorString,
		Message: message,
	})
	rw.Write(bytes)
}

// utility method to get int or error.
func getIntOrError(rw http.ResponseWriter, r *http.Request, key string) (int, error) {
	value := r.URL.Query().Get(key)
	converted, err := strconv.Atoi(value)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(fmt.Sprintf(key)))
	}
	return converted, err
}
