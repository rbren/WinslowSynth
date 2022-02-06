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
  window.state = state;

  const instInfo = state.Instrument.Info;
  if (instInfo?.History?.Samples) {
    addHistory(instInfo.History);
  }
  drawWaveForm(state.Frequency);
}

function clearState() {
  window.state = null;
}

function toggleFreeze(group, name, val) {
  console.log('freeze', group, name, val);
  window.freeze[group + '/' + name] = val;
}

function randomize() {
  $('.constant').each(function(i, div) {
    const e = $(this);
    const [group, name] = e.find('label').attr('title').split('/');
    const key = group + '/' + name;
    const freeze = window.freeze[key];
    if (freeze) return;
    const inpt = e.find('input[type="number"]');
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

function reorderHistory(hist) {
  const firstPos = hist.Position
  const reordered = hist.Samples
    .slice(firstPos, hist.Samples.length)
    .concat(hist.Samples.slice(0, firstPos));
  return reordered;
}

function getAvg(arr) {
  return arr.reduce((a, b) => Math.abs(a) + Math.abs(b), 0.0) / arr.length > 0 ? 'X' : 'O';
}

function logHistory(hist, reord) {
  const qLen = reord.length / 4;
  const q1 = reord.slice(0, qLen);
  const q2 = reord.slice(qLen, 2 * qLen);
  const q3 = reord.slice(qLen * 2, qLen * 3);
  const q4 = reord.slice(qLen * 3, qLen * 4);
  console.log(hist.Position, hist.Time, reord.length, getAvg(q1), getAvg(q2), getAvg(q3), getAvg(q4))
}

function addHistory(hist) {
  const reordered = reorderHistory(hist);
  logHistory(hist, reordered);
  window.sampleHistory = window.sampleHistory || [];
  window.sampleHistoryTime = window.sampleHistoryTime || -1;
  const firstNewTime = hist.Time - hist.Samples.length;
  const lastNewTime = hist.Time;
  const lastSeenTime = window.sampleHistoryTime;
  const numNewFrames = lastNewTime - lastSeenTime;
  let oldestNewFrame = hist.Samples.length - numNewFrames;
  if (oldestNewFrame < 0) {
    console.error('skipped', -oldestNewFrame, 'frames');
    oldestNewFrame = 0;
  }
  const newFrames = reordered.slice(oldestNewFrame, reordered.length);
  window.sampleHistory = window.sampleHistory.concat(newFrames);
  const desiredHistoryLength = 48000;
  window.sampleHistory.splice(0, window.sampleHistory.length - desiredHistoryLength);
  window.sampleHistoryTime = hist.Time;
}
