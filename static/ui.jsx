// the size of the board.
var N = 8;

MenuItem = React.createClass({
  getInitialState: function() {
    return {clicked: false}
  },
  setSelect: function(value) {
    this.setState({selected: value});
  },
  render: function() {
    var selectedClass = "";
    var onClick = (function() {this.props.onSelect(this);}).bind(this);
    if (this.state.selected) {
      selectedClass = "pure-menu-selected"
    }
    return <li className={selectedClass}><a href="#" onClick={onClick}>{
      this.props.name
    }</a></li>
  }
});

Piece = React.createClass({
  render: function() {
    return (
        <svg>
            <circle></circle>
        </svg>
    )
  }
});
Board = React.createClass({
  getInitialState: function() {
    return {
      board: []
    }
  },
  render: function() {
    var result = [];
    for (var row = 0; row < 8; row++) {
      for (var col = 0; col < 8; col++) {
        var divStyle = {
          top: row * 32,
          left: col * 32
        };
        result.push(<div className="cell" style={divStyle}>X</div>);
      }
    }
    return <div className="board">{result}</div>;
  }
});

MenuBar = React.createClass({
  render: function() {
    var selected = null;
    function onSelect(selectedMenuItem) {
      if (selected !== null) {
        selected.setSelect(false)
      }
      selected = selectedMenuItem;
      selected.setSelect(true);
    };

    function onNewGame(selectedMenuItem) {
      $.ajax("/game/new", {
        type: "POST"
      }).then(function(objectString) {
        // render the object.
        var parsed = JSON.parse(objectString);

        debugger
        console.log("Hello, world");
      });
      onSelect(selectedMenuItem);
    }

    var menuItems = [
      <MenuItem onSelect={onNewGame} name="New Game"></MenuItem>,
      <MenuItem onSelect={onSelect} name="Existing Game"></MenuItem>
    ];

    return (
      <div className="pure-menu pure-menu-open pure-menu-horizontal">
        <ul>{ menuItems }</ul>
      </div>
    )
  }
})

Main = React.createClass({
  render: function() {
    return (
    <div>
      <MenuBar/>
      <Board/>
    </div>
    )
  }
});

React.render(
  <Main/>,
  document.getElementById('content')
);
