function setState(state) {
  window.state = state;
  $("#data--time").text(state.Time);
}
