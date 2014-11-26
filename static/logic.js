/// <reference path="typings/jquery/jquery.d.ts" />
var WS_CONNECT = 'ws://localhost:8080/game/websocket';
var App = (function () {
    function App() {
        this.mux = new Mux(WS_CONNECT);
        this.ajax = new Ajax();
    }
    return App;
})();
/**
 * Websocket multiplexer.
 */
var Mux = (function () {
    function Mux(connection) {
        this.socket = new WebSocket(connection);
        this.socket.addEventListener("onopen", function () {
            console.log("hello, world");
        });
    }
    return Mux;
})();
/**
 * Abstraction over the ajax requests.
 */
var Ajax = (function () {
    function Ajax() {
    }
    Ajax.prototype.newGame = function () {
        return $.ajax("/game/new", {
            type: "POST"
        });
    };
    return Ajax;
})();
var mux = new Mux(WS_CONNECT);
//# sourceMappingURL=logic.js.map