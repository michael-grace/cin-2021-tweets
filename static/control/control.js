/**
    URY Tweet Board
    Candidate Interview Night 2021

    Author: Michael Grace
    Date: November 2020, January 2022

    github.com/UniversityRadioYork
 */

let scheme = window.location.protocol === "https:" ? "wss://" : "ws://"

var alert = document.getElementById("server");

const handleWs = () => {

    let ws = new WebSocket(scheme + window.location.host + "/control-ws");

    const createBlockedUser = (user) => {
        let userCard = document.createElement("DIV");
        userCard.classList.add("card");
        userCard.id = user

        let userCardBody = document.createElement("DIV")
        userCardBody.classList.add("card-body")

        let userName = document.createElement("P")
        userName.innerText = user;
        userName.classList.add("card-text")

        userCardBody.appendChild(userName)

        let unblockButton = document.createElement("BUTTON");
        unblockButton.classList.add("btn", "btn-warning", "btn-sm");
        unblockButton.innerText = "Unblock";

        unblockButton.onclick = () => {
            ws.send(JSON.stringify({
                "content": user,
                "action": "UNBLOCK"
            }))
        }

        userCardBody.appendChild(unblockButton)
        userCard.appendChild(userCardBody)
        document.getElementById("blocked").appendChild(userCard)
    }

    fetch("/info").then(d => d.json()).then(j => {
        document.getElementById("hashtag").innerHTML = j.hashtags.join(", ");
        j.blockedUsers.forEach((user) => {
            createBlockedUser(user)
        })
    })

    ws.onopen = () => {
        alert.innerText = "Connected to Server";
        alert.classList.remove("alert-info", "alert-danger");
        alert.classList.add("alert-success");
    };

    ws.onclose = () => {
        alert.innerText = "Disconnected from Server";
        alert.classList.remove("alert-info", "alert-success");
        alert.classList.add("alert-danger");
        console.log("disconnected...retrying in 1 sec")
        setTimeout(() => { handleWs() }, 1000)
    };

    ws.onmessage = (event) => {

        var message = JSON.parse(event.data);

        if (message.action === "CLEAR_CONTROL") {
            document.getElementById("tweets").innerHTML = "";
            return
        } else if (message.action == "REMOVE") {
            document.getElementById(message.id).remove();
            return
        } else if (message.action == "UNBLOCK") {
            document.getElementById(message.user).remove()
            return
        } else if (message.action == "BLOCK") {
            createBlockedUser(message.user)
            return
        } else if (message.action == "RECENT") {
            return
        } else if (message.action == "UNRECENT") {
            return
        }

        // CONSIDER TWEET

        // Now lets put the tweet on the control screen
        var tweet = document.createElement("DIV");
        tweet.classList.add("card");
        tweet.id = message.tweet.id.toString()

        var tweetCardBody = document.createElement("DIV");
        tweetCardBody.classList.add("card-body");

        var tweetTitle = document.createElement("H4");
        tweetTitle.innerText = message.tweet.name + " - @" + message.tweet.user;
        tweetTitle.classList.add("card-title");
        tweetCardBody.appendChild(tweetTitle);

        var tweetBody = document.createElement("P");
        tweetBody.innerText = message.tweet.tweet;
        tweetBody.classList.add("card-text");
        tweetCardBody.appendChild(tweetBody);

        var acceptButton = document.createElement("BUTTON");
        acceptButton.classList.add("btn", "btn-primary", "btn-sm");
        acceptButton.innerText = "Accept Tweet";

        acceptButton.onclick = () => {
            ws.send(JSON.stringify({
                "content": message.tweet.id,
                "action": "ACCEPT"
            }));
        }

        tweetCardBody.appendChild(acceptButton);

        var rejectButton = document.createElement("BUTTON");
        rejectButton.classList.add("btn", "btn-danger", "btn-sm");
        rejectButton.innerText = "Reject Tweet";

        rejectButton.onclick = () => {
            ws.send(JSON.stringify({
                "content": message.tweet.id,
                "action": "REJECT"
            }));
        }

        tweetCardBody.appendChild(rejectButton);

        var blockButton = document.createElement("BUTTON");
        blockButton.classList.add("btn", "btn-warning", "btn-sm");
        blockButton.innerText = "Block User";

        blockButton.onclick = () => {
            if (confirm("Are you sure you want to block @" + message.tweet.user + "?")) {
                ws.send(JSON.stringify({
                    "content": message.tweet.id,
                    "action": "BLOCK"
                }));
            }
        }

        tweetCardBody.appendChild(blockButton);

        tweet.appendChild(tweetCardBody);
        document.getElementById("tweets").appendChild(tweet)
    }


    document.getElementById("clear").onclick = () => {
        ws.send(JSON.stringify({
            "action": "CLEAR_CONTROL"
        }))
    }

    document.getElementById("clear-board").onclick = () => {
        ws.send(JSON.stringify({
            "action": "CLEAR_BOARD"
        }))
    }

}

handleWs()