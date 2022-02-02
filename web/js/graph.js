(function() {
  const margin = {top: 20, right: 0, bottom: 20, left: 20},
      width = 400 - margin.left - margin.right,
      height = 200 - margin.top - margin.bottom;

  var x = d3.scaleLinear().range([0, width]);
  var y = d3.scaleLinear().range([height, 0]);
  x.domain([0, 480]); // TODO: dynamic number of samples
  y.domain([-5.0, 5.0]);

  var valueline = d3.line()
      .x(function(d, idx) { return x(idx); })
      .y(function(d, idx) { return y(d); });


  window.setUpGraph = function(id) {
    var svg = d3.select(id).append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
      .append("g")
        .attr("transform",
              "translate(" + margin.left + "," + margin.top + ")");

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
  }

  window.drawGraph = function(id, data) {
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

  setUpGraph("#Graphs");
})();