/**
    URY Tweet Board
    Candidate Interview Night 2021

    Author: Michael Grace
    Date: November 2020, January 2022

    github.com/UniversityRadioYork
 */



fetch("/info").then(d => d.json()).then(j => {
    document.getElementById("hashtag").innerHTML = j.join(", ");
})

// WebSocket Connection

let scheme = window.location.protocol === "https:" ? "wss://" : "ws://"

var ws = new WebSocket(scheme + window.location.host + "/control-ws");
var alert = document.getElementById("server");

ws.onopen = function() {
    alert.innerText = "Connected to Server";
    alert.classList.remove("alert-info", "alert-danger");
    alert.classList.add("alert-success");
};

ws.onclose = function() {
    alert.innerText = "Disconnected from Server";
    alert.classList.remove("alert-info", "alert-success");
    alert.classList.add("alert-danger");
};

ws.onmessage = function(event) {
    if (event.data === "CLEAR") {
        document.getElementById("tweets").innerHTML = "";
        return
    }

    var message = JSON.parse(event.data);

    if (message.action == "REMOVE") {
        document.getElementById(message.id).remove();
        return
    }

    // Now lets put the tweet on the control screen
    var tweet = document.createElement("DIV");
    tweet.classList.add("card");
    tweet.id = message.id.toString()

    var tweetCardBody = document.createElement("DIV");
    tweetCardBody.classList.add("card-body");

    var tweetTitle = document.createElement("H4");
    tweetTitle.innerText = message.name + " - @" + message.user;
    tweetTitle.classList.add("card-title");
    tweetCardBody.appendChild(tweetTitle);

    var tweetBody = document.createElement("P");
    tweetBody.innerText = message.tweet;
    tweetBody.classList.add("card-text");
    tweetCardBody.appendChild(tweetBody);

    var acceptButton = document.createElement("BUTTON");
    acceptButton.classList.add("btn", "btn-primary", "btn-sm");
    acceptButton.innerText = "Accept Tweet";

    acceptButton.onclick = () => {
        ws.send(JSON.stringify({
            "id": message.id,
            "decision": "ACCEPT"
        }));
        document.getElementById(message.id.toString()).remove();
    }

    tweetCardBody.appendChild(acceptButton);

    var rejectButton = document.createElement("BUTTON");
    rejectButton.classList.add("btn", "btn-danger", "btn-sm");
    rejectButton.innerText = "Reject Tweet";

    rejectButton.onclick = () => {
        ws.send(JSON.stringify({
            "id": message.id,
            "decision": "REJECT"
        }));
        document.getElementById(message.id.toString()).remove();
    }

    tweetCardBody.appendChild(rejectButton);

    var blockButton = document.createElement("BUTTON");
    blockButton.classList.add("btn", "btn-warning", "btn-sm");
    blockButton.innerText = "Block User";

    blockButton.onclick = () => {
        if (confirm("Are you sure you want to block @" + message.user + "?")) {
            ws.send(JSON.stringify({
                "id": message.id,
                "decision": "BLOCK"
            }));
        }
        document.getElementById(message.id.toString()).remove();
    }

    tweetCardBody.appendChild(blockButton);

    tweet.appendChild(tweetCardBody);
    document.getElementById("tweets").appendChild(tweet)
}

document.getElementById("clear").onclick = function() {
    ws.send("CLEAR")
}