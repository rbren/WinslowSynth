const STATE_KEYS = ["Time", "Instrument", "Instruments", "Constants"];
function setState(state) {
  window.freeze = window.freeze || {};
  STATE_KEYS.forEach(key => {
    if (!window.state || JSON.stringify(window.state[key]) !== JSON.stringify(state[key])) {
      var drawKey = "draw" + key;
      console.log('draw',drawKey)
      $("#" + key).html(window[drawKey](state[key]));
    }
  })
  findHistories(state.Instrument);
  window.state = state;
}

function clearState() {
  window.state = null;
}

function toggleFreeze(name, val) {
  console.log('freeze', name, val);
  window.freeze[name] = !window.freeze[name];
}

function randomize() {
  $('.constant').each(function(i, div) {
    const e = $(this);
    const name = e.find('label').text();
    const freeze = window.freeze[name];
    if (freeze) return;
    const inpt = e.find('input[type="number"]');
    console.log('name', name, inpt.val(), freeze);
    const min = parseFloat(inpt.attr('min'))
    const max = parseFloat(inpt.attr('max'))
    const rand = min + Math.random() * (max - min);
    updateConstant(name, rand);
  });
}

function findHistories(inst) {
  if (typeof inst !== 'object') return [];
  if (inst.Info?.History) {
    console.log('hist', inst.Info.Name, inst.Info.Group, inst.Info.History);
  }
  for (let key in inst) {
    if (key === 'Info') continue;
    return findHistories(inst[key]);
  }
}
