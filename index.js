var connection = null;
var clientID = 0;

function setUsername() {
  console.log("***SETUSERNAME");
  var msg = {
    name: document.getElementById("name").value,
    date: Date.now(),
    id: clientID,
    type: "username",
  };
  connection.send(JSON.stringify(msg));
}

function addCommands() {
  let c = document.getElementById("commands");
  var commandMap = new Map([
    ["C", "Cold"],
    ["H", "Hot"],
    ["P", "Milk powder"],
    ["a", "🍏"],
    ["b", "🍌"],
    ["c", "Chocolate"],
    ["d", "drink"],
    ["e", "eat"],
    ["e", "jaffle"],
    ["j", "jump"],
    ["l", "🩸 bag"],
    ["m", "drink magic potion"],
    ["p", "🍑"],
    ["r", "run!!"],
    ["t", "t-shirt"],
    ["w", "drink water"],
    ["y", "eat ryvita"],
    ["z", "zombie!"],
  ]);

  commandMap.forEach((name, key) => {
    var b = document.createElement("BUTTON");
    b.textContent = name;
    b.onclick = function () {
      connection.send(key);
    };
    c.appendChild(b);
  });
}

function updateStat(name, value) {
  let stat = document.getElementById(name);
  stat.textContent = name + ": " + value;
  stat.className = value.split(" ")[0];
}

function appendLog(line) {
  const para = document.createElement("p");
  console.log(line);
  para.textContent = line;
  document.getElementById("logs").prepend(para);
}

function connect() {
  var serverUrl;
  var scheme = "ws";

  // If this is an HTTPS connection, we have to use a secure WebSocket
  // connection too, so add another "s" to the scheme.

  if (document.location.protocol === "https:") {
    scheme += "s";
  }

  serverUrl = scheme + "://" + document.location.host + "/ws";

  connection = new WebSocket(serverUrl);
  console.log("***CREATED WEBSOCKET");

  connection.onopen = function (evt) {
    const reInfo = /^(\w+)(| \w+): (.+)/;
    connection.onmessage = function (evt) {
      console.log("***ONMESSAGE");
      console.log(evt.data);
      let matche = evt.data.match(reInfo);
      if (matche) {
        console.log("hi");
        if (matche[1] === "STAT") {
          updateStat(matche[2].substring(1), matche[3]);
        } else if (matche[1] === "LOG") {
          appendLog(matche[3]);
        }
      }
    };
    console.log("***CREATED ONMESSAGE");
    let name = document.getElementById("name").value;
    connection.send("AUTH " + name);
    addCommands();
  };
}

function send() {
  console.log("***SEND");
  var msg = {
    text: document.getElementById("text").value,
    type: "message",
    id: clientID,
    date: Date.now(),
  };
  connection.send(JSON.stringify(msg));
  document.getElementById("text").value = "";
}
