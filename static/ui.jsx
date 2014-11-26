// the size of the board.
var N = 8;
var SIZE = 32;

MenuItem = React.createClass({
    getInitialState: function () {
        return {clicked: false}
    },
    setSelect: function (value) {
        this.setState({selected: value});
    },
    render: function () {
        var selectedClass = "";
        var onClick = (function () {
            this.props.onSelect(this);
        }).bind(this);
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
            board: new logic.Board("")
        };
    },
    render: function () {
        console.log("rendering...");
        var result = [];
        function bindOnclick(row, col) {
            return function() {
                console.log(row, col, 'clicked');
            }
        }
        for (var row = 0; row < 8; row++) {
            for (var col = 0; col < 8; col++) {
                var cell = this.state.board.getCell(row, col);
                var fill = "player-" + cell.player;
                result.push(<Piece row={row} col={col} fill={fill} onClick={bindOnclick(row, col)}/>);
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
                var board = new logic.Board(parsed.board);
                self.refs.board.setState({board: board});
            });
            onSelect(selectedMenuItem);
        }

        return (
            <div>
                <div className="pure-menu pure-menu-open pure-menu-horizontal">
                    <ul>
                        <MenuItem onSelect={onNewGame} name="New Game"></MenuItem>
                        <MenuItem onSelect={onSelect} name="Existing Game"></MenuItem>
                    </ul>
                </div>
                <Board ref="board"/>
            </div>
        )
    }
});

React.render(
    <Main/>,
    document.getElementById('content')
);
