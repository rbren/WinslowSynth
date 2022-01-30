function setState(state) {
  var curInst = window.state?.Instrument;
  var curConst = window.state?.Constants;
  window.state = state;
  $("#data--time").text(state.Time);
  if (JSON.stringify(curInst) !== JSON.stringify(state.Instrument)) {
    $("#instrument").html(drawInstrument(state.Instrument));
  }
  if (JSON.stringify(curConst) !== JSON.stringify(state.Constants)) {
    $("#constants").html(drawConstants(state.Constants));
  }
}

function updateConstant(name, val) {
  console.log('update', name, val);
  ws.send(JSON.stringify({
    "Key": name,
    "Value": parseFloat(val),
    "Action": "set",
  }));
}

function drawConstants(consts) {
  return consts.map(drawConstant).join('');
}

function drawConstant(constant) {
  return `
  <div>
  <label>${constant.Name}</label>
  <br>
  <input
    class="slider"
    type="range"
    min="${constant.Min}"
    max="${constant.Max}"
    value="${constant.Value}"
    onchange="updateConstant('${constant.Name}', this.value)"
    >
    </div>
  `
}

function drawInstrument(inst) {
  const fields = ['Amplitude', 'Frequency', 'Phase', 'Bias'];
  return `
    <table>
    ${fields.map(f => drawField(f, inst[f])).join('')}
    </table>
  `;
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
