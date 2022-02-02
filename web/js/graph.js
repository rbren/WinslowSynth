console.log('graph.js');
(function() {
  console.log('set up graphs');
  const margin = {top: 20, right: 20, bottom: 50, left: 70},
      width = 960 - margin.left - margin.right,
      height = 500 - margin.top - margin.bottom;

  var x = d3.scaleLinear().range([0, width]);
  var y = d3.scaleLinear().range([height, 0]);
  x.domain([0, 480]); // TODO: dynamic number of samples
  y.domain([-5.0, 5.0]);

  var valueline = d3.line()
      .x(function(d, idx) { return x(idx); })
      .y(function(d, idx) { return y(d); });


  window.setUpGraph = function(id) {
    console.log('setup', id);
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
            console.log('enter', enter);
            return enter.append("path").attr("class", "line");
          },
          update => {
            console.log('update', update);
            return update.attr("d", valueline);
          },
          exit => {
            console.log('exit', exit);
          })
        .attr("d", valueline);
  }

  setUpGraph("#Graphs");
})();
