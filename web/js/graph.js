(function() {
  const margin = {top: 20, right: 0, bottom: 20, left: 30},
      width = 400 - margin.left - margin.right,
      height = 200 - margin.top - margin.bottom;

  window.setUpGraph = function(id, xDomain, yDomain) {
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

    // Add the x Axis
    svg.append("g")
        .attr("transform", "translate(0," + height + ")")
        .call(d3.axisBottom(x));

    // Add the y Axis
    svg.append("g")
        .call(d3.axisLeft(y));

    // The eventual data
    svg.append("path")
        .attr("class", "line")

    return valueline;
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

  setUpGraph("#WaveFormGraph", [0, 440.0], [-1.0, 1.0]);
})();

const drawHistoryIntervalMs = 50;
function startDrawHistoryInterval() {

  const samplesPerSecond = window.state.Config.SampleRate;
  const newSamplesPerInterval = (samplesPerSecond / 1000) * drawHistoryIntervalMs;

  console.log('drawing', newSamplesPerInterval, 'every', drawHistoryIntervalMs);

  const id = "#HistoryGraph";
  const numSamplesInGraph = newSamplesPerInterval * 1;
  const valueline = setUpGraph(id, [0, numSamplesInGraph], [-1.0, 1.0]);

  const startTime = window.sampleHistoryTime - window.sampleHistory.length;
  let curTime = startTime;
  let samplesToDraw = [];
  return setInterval(() => {
    const firstAvailableTime = window.sampleHistoryTime - window.sampleHistory.length;
    if (firstAvailableTime > curTime) {
      console.error('skip drawing', firstAvailableTime - curTime, 'frames');
      curTime = firstAvailableTime;
    }
    const firstSampleIdx = curTime - firstAvailableTime;
    let newSamples = window.sampleHistory.slice(firstSampleIdx, firstSampleIdx + newSamplesPerInterval);
    samplesToDraw = samplesToDraw.concat(newSamples);
    if (samplesToDraw.length > numSamplesInGraph) {
      samplesToDraw.splice(0, samplesToDraw.length - numSamplesInGraph);
    }
    drawGraph(id, valueline, samplesToDraw);
    curTime += newSamples.length;
    $("#DebugTime").html(curTime);
  }, drawHistoryIntervalMs);
}
