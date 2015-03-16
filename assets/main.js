var width = 960,
    height = 500;

var force = d3.layout.force()
    .charge(-120)
    .linkDistance(30)
    .size([width, height]);

var svg = d3.select("body").append("svg")
    .attr("width", width)
    .attr("height", height);

var emojis = {};

d3.json("/messages.json", function(error, data) {

    var links = [];
    data.forEach(function(d) {
        var words = d.text.split(" ");
        words.forEach(function(word) {
            if(word.search(/^:\w+:$/) != -1) {
                if(emojis[word] == null) {
                    emojis[word] = {text: word, messages:[]}
                }
                emojis[word].messages.push(d)
                links.push({source: word, target: d});
            }
        }) 
    })

    var nodes = []
    for(var emoji in emojis) {
        nodes.push(emoji);
    }

  force
      .nodes(nodes)
      .links(links)
      .start();

    var link = svg.selectAll(".link")
      .data(links)
    .enter().append("line")
      .attr("class", "link")
      .style("stroke-width", function(d) { return Math.sqrt(d.value); });

  var node = svg.selectAll(".node")
      .data(nodes)
    .enter().append("circle")
      .attr("class", "node")
      .attr("r", 5)
      .style("fill", function(d) { return "#FF0000"; })
      .call(force.drag);

  force.on("tick", function() {
    link.attr("x1", function(d) { return d.source.x; })
        .attr("y1", function(d) { return d.source.y; })
        .attr("x2", function(d) { return d.target.x; })
        .attr("y2", function(d) { return d.target.y; });

})
