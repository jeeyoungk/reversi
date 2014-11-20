MenuItem = React.createClass({
  getInitialState: function() {
    return {clicked: false}
  },
  setSelect: function(value) {
    this.setState({selected: value});
  },
  render: function() {
    var selectedClass = "";
    var onClick = (function() {this.props.onSelect(this);}).bind(this)
    if (this.state.selected) {
      selectedClass = "pure-menu-selected"
    }
    return <li className={selectedClass}><a href="#" onClick={onClick}>{
      this.props.name
    }</a></li>
  }
})
MenuBar = React.createClass({
  render: function() {
    var selected = null;
    function onSelect(selectedMenuItem) {
      if (selected !== null) {
        selected.setSelect(false)
      }
      selected = selectedMenuItem;
      selected.setSelect(true);
    }
    var menuItems = [
      <MenuItem onSelect={onSelect} name="New Game"></MenuItem>,
      <MenuItem onSelect={onSelect} name="Existing Game"></MenuItem>
    ]
    return (
      <div className="pure-menu pure-menu-open pure-menu-horizontal">
        <ul>{ menuItems }</ul>
      </div>
    )
  }
})

React.render(
  <MenuBar/>,
  document.getElementById('content')
);
