function drawInstruments(names) {
  return `
    <select onchange="chooseInstrument(this.value)">
    <option value="">Instruments</option>
    ${names.map(drawInstrumentChoice).join('')}
    </select>`;
}

function drawInstrumentChoice(name) {
  return `<option value="${name}">${name}</option>`;
}

function drawTime(t) {
  return t;
}

function drawStatus(ok) {
  return `
  <span class="badge bg-${ok ? 'success' : 'danger'}">
  ${ok ? 'Connected' : 'Disconnected'}
  </span>
  `
}

function drawInstrument(inst) {
  $('#Title').html(inst.Info?.Name || "Synthesizer");
  return "";
}

function drawHistory(hist) {
  const firstPos = hist.HistoryPosition % hist.History.length;
  const reordered = hist.History
    .slice(firstPos, hist.History.length)
    .concat(hist.History.slice(0, firstPos));

  window.drawGraph("#Graphs", reordered);
}

function drawConstants(consts) {
  const groups = {}
  consts.forEach(c => {
    groups[c.Info.Group] = groups[c.Info.Group] || [];
    groups[c.Info.Group].push(c);
  });
  return Object.keys(groups).map(k => {
    return drawConstantGroup(k, groups[k]);
  }).join('');
}

window.drawInstruments = drawInstruments;
window.drawTime = drawTime;
window.drawInstrument = drawInstrument;
window.drawConstants = drawConstants;

function drawConstantGroup(name, constants) {
  return `
  <div class="constant-group">
  <h2>${name}</h2>
  ${constants.map(drawConstant).join('')}
  </div>
  `
}

function drawConstant(constant) {
  step = (constant.Max - constant.Min) / 100.0
  props = `
      step="${step}"
      min="${constant.Min}"
      max="${constant.Max}"
      value="${constant.Value}"
      onchange="updateConstant('${constant.Info.Name}', this.value)"
      `;
  return `
  <div class="constant">
    <label title="${constant.Info.Name}">${constant.Info.Name}</label>
    <div class="controls">
      <input
        class="slider"
        type="range"
        ${props}
        >
      <input
        type="number"
        ${props}
        >
      <input type="checkbox"
        onchange="toggleFreeze('${constant.Info.Name}', this.value)"
        ${window.freeze[constant.Info.Name] ? 'checked' : ''}>
    </div>
  </div>
  `
}


function drawField(label, value) {
  return `
  <tr>
    <th>${drawLabel(label)}</th>
    <td>${drawValue(value)}</td>
  </tr>
  `;
}

function drawLabel(label) {
  return label
}

function drawValue(value) {
  if (typeof value === 'number') {
    return value;
  }
  if (!value) {
    return `nil`
  }
  if (typeof value === 'object') {
    //return `<a onclick='drawInstrument(${JSON.stringify(value)})'>expand</a>`;
    return `<a onclick="console.log('click')">expand</a>`;
  }
}

