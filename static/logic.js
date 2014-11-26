/// <reference path="typings/jquery/jquery.d.ts" />
var logic;
(function (logic) {
    /**
     * Represents a game board.
     */
    var Board = (function () {
        function Board(data) {
            this.data = data;
        }
        Board.prototype.getCell = function (row, col) {
            var cell = this.data.charAt(row * 8 + col);
            switch (cell) {
                case "X":
                    return new Cell(1);
                case "O":
                    return new Cell(2);
                default:
                    return new Cell(0);
            }
        };
        return Board;
    })();
    logic.Board = Board;
    /**
     * Represents a cell of game.
     */
    var Cell = (function () {
        function Cell(player) {
            this.player = player;
        }
        return Cell;
    })();
})(logic || (logic = {}));
//# sourceMappingURL=logic.js.map