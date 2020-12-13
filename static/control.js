/**
    URY Tweet Board
    Candidate Interview Night 2021

    Author: Michael Grace
    Date: November 2020

    github.com/UniversityRadioYork
 */

document.getElementById("clear").onclick = function() {
    document.getElementById("tweets").innerHTML = "";
}

var xhttp = new XMLHttpRequest();
xhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
        server_data = JSON.parse(this.responseText);
        document.getElementById("hashtag").innerHTML = server_data.hashtag;

        // WebSocket Connection
        var ws = new WebSocket(server_data.ws_conn + "/control");
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
            var message = JSON.parse(event.data);
            // Now lets put the tweet on the control screen
            var tweet = document.createElement("DIV");
            tweet.classList.add("card");
            tweet.id = message.id.toString()

            var tweetCardBody = document.createElement("DIV");
            tweetCardBody.classList.add("card-body");

            var tweetTitle = document.createElement("H4");
            tweetTitle.innerText = message.title;
            tweetTitle.classList.add("card-title");
            tweetCardBody.appendChild(tweetTitle);

            var tweetBody = document.createElement("P");
            tweetBody.innerText = message.body;
            tweetBody.classList.add("card-text");
            tweetCardBody.appendChild(tweetBody);

            var acceptButton = document.createElement("BUTTON");
            acceptButton.classList.add("btn", "btn-primary", "btn-sm");
            acceptButton.innerText = "Accept Tweet";
            acceptButton.onclick = function() {
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
            rejectButton.onclick = function() {
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
            blockButton.onclick = function() {
                if (confirm("Are you sure you want to block @" + message.title.split("@")[1] + "?")) {
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

    }
};
xhttp.open("GET", "/info", true);
xhttp.send();