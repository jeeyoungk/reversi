// author : jeeyoung kim
// front-end code for the Reversi game.
var N = 8;     // size of the board.
var SIZE = 32; // size of each circle.

MenuItem = React.createClass({
    getInitialState: function () {
        return {clicked: false}
    },
    setSelect: function (value) {
        this.setState({selected: value});
    },
    render: function () {
        var selectedClass = "";
        var onClick = () => {
            this.props.onSelect(this);
        };
        if (this.state.selected) {
            selectedClass = "pure-menu-selected"
        }
        return <li className={selectedClass}>
            <a href="#" onClick={onClick}>{
                this.props.name
                }</a>
        </li>
    }
});

Piece = React.createClass({
    render: function () {
        var r = SIZE / 2;
        var cx = this.props.col * SIZE + r;
        var cy = this.props.row * SIZE + r;
        return (
            <circle cx={cx} cy={cy} r={r} className={this.props.fill} onClick={this.props.onClick}></circle>
        )
    }
});

Board = React.createClass({
    getInitialState: function () {
        return {
            board: new BoardModel({id: null, board: "", width: 0, height: 0})
        };
    },
    render: function () {
        var result = [];
        var self = this;
        function bindOnclick(row, col) {
            var onError = (response) => {
                console.log('error');
            };
            var onSuccess = (response) => {
                console.log('success');
            };
            return function() {
                var board = self.state.board;
                $.ajax("/game/play", {
                    data: {
                        id: board.id,
                        row: row,
                        col: col,
                        player: 1
                    },
                    type: "POST"
                }).then(onSuccess, onError);
                console.log(row, col, 'clicked');
            }
        }
        for (var row = 0; row < this.state.board.width; row++) {
            for (var col = 0; col < this.state.board.height; col++) {
                var cell = this.state.board.getCell(row, col);
                var fill = "player-" + cell.getPlayer();
                result.push(<Piece key={row + ":" + col} row={row} col={col} fill={fill} onClick={bindOnclick(row, col)}/>);
            }
        }
        return <svg className="board">{result}</svg>;
    }
});

Main = React.createClass({
    render: function () {
        var self = this;
        var selected = null;

        function onSelect(selectedMenuItem) {
            if (selected !== null) {
                selected.setSelect(false)
            }
            selected = selectedMenuItem;
            selected.setSelect(true);
        }

        function onNewGame(selectedMenuItem) {
            $.ajax("/game/new", {
                type: "POST"
            }).then(function (objectString) {
                // render the object.
                var parsed = JSON.parse(objectString);
                var board = new BoardModel(parsed);
                self.refs.board.setState({board: board});
            });
            onSelect(selectedMenuItem);
        }

        return (
            <div>
                <div className="pure-menu pure-menu-open pure-menu-horizontal">
                    <ul>
                        <MenuItem onSelect={onNewGame} name="New Game"></MenuItem>
                    </ul>
                </div>
                <Board ref="board"/>
            </div>
        )
    }
});

class BoardModel {
    constructor(options) {
        this.id =     options.id;
        this.width =  options.width  || 0;
        this.height = options.height || 0;
        this.board = [];
        for (var index = 0; index < options.board.length; index++) {
            var cell = new CellModel(options.board.charAt(index));
            this.board.push(cell);
        }
    }

    getCell(row, column) {
        return this.board[row * 8 + column];
    }
}

class CellModel {
    constructor(value) {
        this.value = value;
    }

    getPlayer() {
        if (this.value == "X") {
            return 1;
        }
        if (this.value == "O") {
            return 2;
        }
        return 0;
    }
}

React.render(
    <Main/>,
    document.getElementById('content')
);
