console.log("socket");
window.addEventListener("load", function(evt) {
  console.log("load");
  var ws = window.ws = new WebSocket("ws://" + window.location.host + "/connect");
  ws.onopen = function(evt) {
    $('#Status').html(drawStatus(true));
    initialize();
  }
  ws.onclose = function(evt) {
    ws = null;
    $('#Status').html(drawStatus(false));
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

function updateConstant(group, name, val) {
  console.log('update', group, name, val);
  clearState();
  ws.send(JSON.stringify({
    Key: group + "/" + name,
    Value: parseFloat(val),
    Action: "set",
  }));
}

function chooseInstrument(name) {
  console.log("choose", name);
  clearState();
  ws.send(JSON.stringify({
    Key: name,
    Action: "choose"
  }));
}
