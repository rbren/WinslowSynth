function drawInstruments(names) {
  return `<ul>${names.map(drawInstrumentChoice).join('')}</ul>`;
}

function drawTime(t) {
  return t;
}

function drawInstrument(inst) {
  const fields = ['Amplitude', 'Frequency', 'Phase', 'Bias'];
  return `
    <table>
    ${fields.map(f => drawField(f, inst[f])).join('')}
    </table>
  `;
}

function drawConstants(consts) {
  return consts.map(drawConstant).join('');
}

window.drawInstruments = drawInstruments;
window.drawTime = drawTime;
window.drawInstrument = drawInstrument;
window.drawConstants = drawConstants;

function drawInstrumentChoice(name) {
  return `<li><a href="#" onclick="chooseInstrument('${name}')">${name}</a></li>`;
}

function drawConstant(constant) {
  step = (constant.Max - constant.Min) / 100.0
  props = `
      step="${step}"
      min="${constant.Min}"
      max="${constant.Max}"
      value="${constant.Value}"
      onchange="updateConstant('${constant.Name}', this.value)"
      `;
  return `
  <div>
    <label>${constant.Name}</label>
    <br>
    <input
      class="slider"
      type="range"
      ${props}
      >
    <input
      type="number"
      ${props}
      >
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

