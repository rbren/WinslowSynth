const STATE_KEYS = ["Time", "Instrument", "Instruments", "Constants"];

function setState(state) {
  window.freeze = window.freeze || {};
  state.Constants = findConstants(state.Instrument);
  STATE_KEYS.forEach(key => {
    if (!window.state || JSON.stringify(window.state[key]) !== JSON.stringify(state[key])) {
      var drawKey = "draw" + key;
      $("#" + key).html(window[drawKey](state[key]));
    }
  });
  const instInfo = state.Instrument.Info;
  if (instInfo && instInfo.History) {
    drawHistory(instInfo);
  }
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
    const [group, name] = e.find('label').attr('title').split('/');
    const freeze = window.freeze[name];
    if (freeze) return;
    const inpt = e.find('input[type="number"]');
    console.log('name', name, inpt.val(), freeze);
    const min = parseFloat(inpt.attr('min'))
    const max = parseFloat(inpt.attr('max'))
    const rand = min + Math.random() * (max - min);
    updateConstant(group, name, rand);
  });
}

function findHistories(inst) {
  if (!inst || typeof inst !== 'object') return [];
  let all = [];
  if (inst.Info?.History) {
    all.push(inst.Info);
  }
  for (let key in inst) {
    if (key === 'Info') continue;
    all = all.concat(findHistories(inst[key]));
  }
  return all;
}

function isConst(inst) {
  if (!inst) return false;
  if (!inst.Info?.Name) return false;
  if (inst.Info.Name === 'Frequency') return false;
  if (inst.Value === undefined) return false;
  if (inst.Min === undefined) return false;
  if (inst.Max === undefined) return false;
  return true;
}

function findConstants(inst) {
  if (!inst || typeof inst !== 'object') return [];
  let all = [];
  if (isConst(inst)) {
    all.push(inst);
  }
  for (let key in inst) {
    if (key === 'Info') continue;
    all = all.concat(findConstants(inst[key]));
  }
  return all;
}
