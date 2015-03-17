var width = 1500;
var height = 900;

var force = d3.layout.force()
    .charge(-120)
    .linkDistance(60)
    .size([width, height]);

var svg = d3.select("body").append("svg")
    .attr("width", width)
    .attr("height", height);


var usersByID = {};
function loadUsers() {
  d3.json("/users.json", function(error, data) {
    usersByID = data;
    loadMessages();
  });
}

function loadMessages() {
  d3.json("/messages.json", function(error, data) {
      var nodes = [], nodesByWord = {}, nodesByUser = {};
      var links = [];

      data.forEach(function(d) {
          d.text.split(" ").forEach(function(word) {
              // Skip if it's not an emoji.
              if(word.search(/^:\w+:$/) == -1) {
                return;
              }
              word = word.slice(1, -1);

              switch(word) {
                case "simple_smile": word = "smile"; break;
                case "tddface": word = "trollface"; break;
                case "upvote": word = "+1"; break;
              }

              // Add node.
              if(nodesByWord[word] == null) {
                nodesByWord[word] = nodes.length;
                nodes.push({
                  type: "word",
                  name: word,
                  image: "http://www.tortue.me/emoji/" + word + ".png"
                });
              }

              // Retrieve user.
              var username;
              var image;
              if(usersByID[d.user] != null) {
                username = usersByID[d.user].name;
                image = usersByID[d.user].profile.image_24;
              }

              // Add user.
              if(nodesByUser[d.user] == null) {
                nodesByUser[d.user] = nodes.length;
                nodes.push({type: "user", name: username, image: image});
              }

              // Generate link.
              links.push({
                source: nodesByUser[d.user],
                target: nodesByWord[word]
              });
          }) 
      })

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
      .enter().append("image")
        .attr("class", "node")
        .call(force.drag);

    node.append("title")
        .text(function(d) { return d.name; });

    force.on("tick", function() {
      link.attr("x1", function(d) { return d.source.x; })
          .attr("y1", function(d) { return d.source.y; })
          .attr("x2", function(d) { return d.target.x; })
          .attr("y2", function(d) { return d.target.y; });

      node.attr("x", function(d) { return d.x - 12; })
          .attr("y", function(d) { return d.y - 12; })
          .attr("width", function(d) { return 24; })
          .attr("height", function(d) { return 24; })
          .attr("xlink:href", function(d) { return d.image; })
          .attr("fill", function(d) {
            if(d.type == "word") {
              return "#000000";
            } else {
              return "#00FF00";
            }
          });
    });
  })
}

// Load users first then messages.
loadUsers();
