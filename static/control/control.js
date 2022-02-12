/**
    URY Tweet Board
    Candidate Interview Night 2021

    Author: Michael Grace
    Date: November 2020, January 2022

    github.com/UniversityRadioYork
 */

let scheme = window.location.protocol === "https:" ? "wss://" : "ws://"

let alert = document.getElementById("server");
let authAlert = document.getElementById("authenticated")

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

        fetch("/ws-auth").then(d => d.json()).then(j => {
            ws.send(JSON.stringify({
                "action": "AUTH",
                "content": j.wsAuthToken
            }))
        })
    };

    ws.onclose = () => {
        alert.innerText = "Disconnected from Server";
        alert.classList.remove("alert-info", "alert-success");
        alert.classList.add("alert-danger");

        authAlert.innerText = "Waiting for Connection to Authenticate";
        authAlert.classList.remove("alert-danger", "alert-success");
        authAlert.classList.add("alert-info");

        console.log("disconnected...retrying in 1 sec")
        setTimeout(() => { handleWs() }, 1000)
    };

    ws.onmessage = (event) => {

        let message = JSON.parse(event.data);

        switch (message.action) {
            case "CLEAR_CONTROL":
                document.getElementById("tweets").innerHTML = "";
                break;
            case "REMOVE":
                document.getElementById("consider-" + message.id).remove();
                break;
            case "UNBLOCK":
                document.getElementById(message.user).remove()
                break;
            case "BLOCK":
                createBlockedUser(message.user)
                break;
            case "UNRECENT":
                document.getElementById("recent-" + message.id).remove();
                break;
            case "CLEAR_BOARD":
                document.getElementById("recents").innerHTML = "";
                break;
            case "AUTH":
                if (message.ok) {
                    authAlert.innerText = "Authenticated";
                    authAlert.classList.remove("alert-info", "alert-danger");
                    authAlert.classList.add("alert-success");
                } else {
                    authAlert.innerText = "Not Authenticated";
                    authAlert.classList.remove("alert-info", "alert-success");
                    authAlert.classList.add("alert-danger");
                }
                break;
            default:

                // CONSIDER tweet or RECENT tweet

                // Now lets put the tweet on the control screen
                let tweet = document.createElement("DIV");
                tweet.classList.add("card");
                tweet.id = (message.action == "CONSIDER" ? "consider-" : "recent-") + message.tweet.id.toString()

                let tweetCardBody = document.createElement("DIV");
                tweetCardBody.classList.add("card-body");

                let tweetTitle = document.createElement("H4");
                tweetTitle.innerText = message.tweet.name + " - @" + message.tweet.user;
                tweetTitle.classList.add("card-title");
                tweetCardBody.appendChild(tweetTitle);

                let tweetBody = document.createElement("P");
                tweetBody.innerText = message.tweet.tweet;
                tweetBody.classList.add("card-text");
                tweetCardBody.appendChild(tweetBody);

                if (message.action == "CONSIDER") {
                    let acceptButton = document.createElement("BUTTON");
                    acceptButton.classList.add("btn", "btn-primary", "btn-sm");
                    acceptButton.innerText = "Accept Tweet";

                    acceptButton.onclick = () => {
                        ws.send(JSON.stringify({
                            "content": message.tweet.id,
                            "action": "ACCEPT"
                        }));
                    }

                    tweetCardBody.appendChild(acceptButton);

                    let rejectButton = document.createElement("BUTTON");
                    rejectButton.classList.add("btn", "btn-danger", "btn-sm");
                    rejectButton.innerText = "Reject Tweet";

                    rejectButton.onclick = () => {
                        ws.send(JSON.stringify({
                            "content": message.tweet.id,
                            "action": "REJECT"
                        }));
                    }

                    tweetCardBody.appendChild(rejectButton);

                    let blockButton = document.createElement("BUTTON");
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

                } else if (message.action == "RECENT") {
                    let removeButton = document.createElement("BUTTON");
                    removeButton.classList.add("btn", "btn-warning", "btn-sm");
                    removeButton.innerText = "Remove from Wall";

                    removeButton.onclick = () => {
                        if (confirm("Are you sure you want to remove this tweet from the wall?")) {
                            ws.send(JSON.stringify({
                                "action": "BOARD_REMOVE",
                                "content": message.tweet.id
                            }))
                        }
                    }

                    tweetCardBody.appendChild(removeButton)
                }

                tweet.appendChild(tweetCardBody);

                document.getElementById(message.action == "CONSIDER" ? "tweets" : "recents").appendChild(tweet)
        }
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