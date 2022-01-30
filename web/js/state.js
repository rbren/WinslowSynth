function setState(state) {
  var curInst = window.state?.Instrument;
  window.state = state;
  $("#data--time").text(state.Time);
  if (JSON.stringify(curInst) !== JSON.stringify(state.Instrument)) {
    $("#instrument").html(drawInstrument(state.Instrument))
  }
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
