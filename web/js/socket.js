console.log("socket");
window.addEventListener("load", function(evt) {
  console.log("load");
  var ws = new WebSocket("ws://" + window.location.host + "/connect");
  ws.onopen = function(evt) {
    console.log("OPEN");
  }
  ws.onclose = function(evt) {
    console.log("CLOSE");
    ws = null;
  }
  ws.onmessage = function(evt) {
    //console.log("RESPONSE: " + evt.data);
    window.setState(JSON.parse(evt.data))
  }
  ws.onerror = function(evt) {
    console.log("ERROR: " + evt.data);
  }

  const pressedKeys = {}

  document.onkeydown = function(event) {
    console.log("keydown", event.key);
    if (!ws) return true;
    if (pressedKeys[event.key]) return true;
    pressedKeys[event.key] = true;
    console.log("DOWN: ", event.key)
    ws.send(JSON.stringify({
      "Key": event.key,
      "Action": "down",
     }))
    return true;
  }

  document.onkeyup = function(event) {
    if (!ws) return true;
    pressedKeys[event.key] = false;
    console.log("UP: ", event.key)
    ws.send(JSON.stringify({
      "Key": event.key,
      "Action": "up",
     }))
    return true;
  }
});