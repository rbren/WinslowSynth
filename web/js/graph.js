(function() {
  const margin = {top: 20, right: 0, bottom: 20, left: 30},
      width = 400 - margin.left - margin.right,
      height = 200 - margin.top - margin.bottom;

  window.setUpGraph = function(id, xDomain, yDomain, opts={}) {
    var svg = d3.select(id).append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
      .append("g")
        .attr("transform",
              "translate(" + margin.left + "," + margin.top + ")");

    var x = d3.scaleLinear().range([0, width]);
    var y = d3.scaleLinear().range([height, 0]);
    x.domain(xDomain);
    y.domain(yDomain);

    var valueline = d3.line()
        .x(function(d, idx) { return x(idx); })
        .y(function(d, idx) { return y(d); });

    if (!opts.hideXAxis) {
      svg.append("g")
          .attr("transform", "translate(0," + height + ")")
          .call(d3.axisBottom(x));
    }

    if (!opts.hideYAxis) {
      svg.append("g")
          .call(d3.axisLeft(y));
    }

    // The eventual data
    svg.append("path")
        .attr("class", "line")

    return {x, y, valueline};
  }

  window.drawGraph = function(id, valueline, data) {
    var svg = d3.select(id + " svg");

    svg.selectAll(".line")
        .data([data])
        .join(
          enter => {
            return enter.append("path").attr("class", "line");
          },
          update => {
            return update.attr("d", valueline);
          },
          exit => {
            return
          })
        .attr("d", valueline);
  }

})();

function setUpWaveFormGraph() {
  if (!window.waveFormGraph) {
    window.waveFormGraph = setUpGraph("#WaveFormGraph", [0, 250], [-1.0, 1.0], {hideXAxis: true});
  }
}

function setUpFrequenciesGraph() {
  if (!window.frequenciesGraph) {
    window.frequenciesGraph = setUpGraph("#FrequenciesGraph", [0, 1000], [-2.0, 2.0]);
  }
}

function drawFrequencies(freqs) {
  if (!freqs || freqs.length === 0) return
  const {x, y} = window.frequenciesGraph;
  var valueline = d3.line()
        .x(function(d, idx) { return x(idx); })
        .y(function(d, idx) { return y(d); });
  drawGraph("#FrequenciesGraph", valueline, freqs);
}

function drawWaveForm(freq) {
  if (freq === 0) return;
  const impulsesToShow = 2;
  const impulsesPerSec = freq;
  const samplesPerSec = window.state.Config.SampleRate;
  const samplesPerImpulse = Math.round(samplesPerSec / impulsesPerSec);
  const lastAvailableSample = window.sampleHistoryTime;
  const modulus = lastAvailableSample % samplesPerImpulse;
  const endIdx = window.sampleHistory.length - modulus;
  const startIdx = endIdx - samplesPerImpulse * impulsesToShow;
  const samples = window.sampleHistory.slice(startIdx, endIdx);

  const {x, y} = window.waveFormGraph;
  x.domain([0, samplesPerImpulse * impulsesToShow]);
  var valueline = d3.line()
        .x(function(d, idx) { return x(idx); })
        .y(function(d, idx) { return y(d); });

  drawGraph("#WaveFormGraph", valueline, samples);
}


const drawHistoryIntervalMs = 50;
function startDrawHistoryInterval() {
  const id = "#HistoryGraph";
  const numSamplesInGraph = 10000;
  const {valueline} = setUpGraph(id, [0, numSamplesInGraph], [-1.0, 1.0]);

  let curEndTime = window.sampleHistoryTime;
  return setInterval(() => {
    if (!window.sampleHistory) return;
    const endIdx = window.sampleHistory.length;
    const startIdx = endIdx - numSamplesInGraph;
    const samplesToDraw = window.sampleHistory.slice(startIdx, endIdx);
    drawGraph(id, valueline, samplesToDraw);

    // some debug info
    const newStartTime = window.sampleHistoryTime - (endIdx - startIdx);
    if (curEndTime < newStartTime) {
      console.error('skipped drawing', newStartTime - curEndTime, 'frames');
    }
    curEndTime = window.sampleHistoryTime;
  }, drawHistoryIntervalMs);
}
