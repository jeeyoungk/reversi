/// <reference path="typings/jquery/jquery.d.ts" />
module logic {
    /**
     * Represents a game board.
     */
    export class Board {
        private data:string;

        constructor(data) {
            this.data = data;
        }

        getCell(row:number, col:number) {
            var cell = this.data.charAt(row * 8 + col);
            switch (cell) {
                case "X":
                    return new Cell(1);
                case "O":
                    return new Cell(2);
                default:
                    return new Cell(0);
            }
        }
    }

    /**
     * Represents a cell of game.
     */
    class Cell {
        public player:number;

        constructor(player:number) {
            this.player = player;
        }
    }
}