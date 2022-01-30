const STATE_KEYS = ["Time", "Instrument", "Instruments", "Constants"];
function setState(state) {
  STATE_KEYS.forEach(key => {
    if (!window.state || JSON.stringify(window.state[key]) !== JSON.stringify(state[key])) {
      var drawKey = "draw" + key;
      console.log('draw',drawKey)
      $("#" + key).html(window[drawKey](state[key]));
    }
  })
  window.state = state;
}

function clearState() {
  window.state = null;
}

function randomize() {
  $('.constant').each(function(i, div) {
    const e = $(this);
    const name = e.find('label').text();
    console.log('name', name);
    const inpt = e.find('input[type="number"]');
    const min = parseFloat(inpt.attr('min'))
    const max = parseFloat(inpt.attr('max'))
    const rand = min + Math.random() * (max - min);
    updateConstant(name, rand);
  });
}

function updateConstant(name, val) {
  console.log('update', name, val);
  clearState();
  ws.send(JSON.stringify({
    Key: name,
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
