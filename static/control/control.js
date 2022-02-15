/**
URY Tweet Board
Copyright (C) 2020, 2022 Michael Grace <michael.grace@ury.org.uk>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

let scheme = window.location.protocol === "https:" ? "wss://" : "ws://"

let alert = document.getElementById("server");
let authAlert = document.getElementById("authenticated")
let started = false;

const handleWs = () => {

    let ws = new WebSocket(scheme + window.location.host + "/control-ws");

    const createBlockedUser = (user) => {
        let userCard = document.createElement("DIV");
        userCard.classList.add("card", "m-2", "border-secondary");
        userCard.id = user

        let userCardBody = document.createElement("DIV")
        userCardBody.classList.add("card-body")

        let userName = document.createElement("LABEL")
        userName.innerText = user;
        userName.classList.add("card-text")

        userCardBody.appendChild(userName)

        let unblockButton = document.createElement("BUTTON");
        unblockButton.classList.add("btn", "btn-warning", "btn-sm", "ml-2");
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

        console.log("disconnected...retrying in 3 sec")
        setTimeout(() => { handleWs() }, 3000)
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

                    if (!started) {
                        ws.send(JSON.stringify({
                            "action": "QUERY"
                        }))
                        started = true
                    }

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
                tweet.classList.add("card", "m-2", "border-secondary");
                tweet.id = (message.action == "CONSIDER" ? "consider-" : "recent-") + message.tweet.id.toString()

                let tweetCardBody = document.createElement("DIV");
                tweetCardBody.classList.add("card-body");

                let tweetTitle = document.createElement("B");
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