package web

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// upgrades the normal HTTP connection to a websocket connection.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (ctx *HTTPServerContext) WebsocketHandler(rw http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	for {
		// TODO - implement timeout & closing.
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Returning...")
			return
		}
		if err = conn.WriteMessage(messageType, p); err != nil {
			log.Println("Returning...")
			return
		}
	}
}
