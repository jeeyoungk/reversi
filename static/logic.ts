/// <reference path="typings/jquery/jquery.d.ts" />

var WS_CONNECT = 'ws://localhost:8080/game/websocket';

/**
 * Websocket multiplexer.
 */
class Mux {
    socket: WebSocket
    constructor (connection: string) {
        this.socket = new WebSocket(connection);
        this.socket.addEventListener("onopen", function() {
            console.log("hello, world")
        });
    }
}

/**
 * Abstraction over the ajax requests.
 */
class Ajax {
    newGame() {
        return $.ajax("/game/new", {
            type: "POST"
        })
    }
}

var mux = new Mux(WS_CONNECT);