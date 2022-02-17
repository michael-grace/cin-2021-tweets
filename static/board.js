/**
URY Tweet Board
Copyright (C) 2020, 2021, 2022 
Michael Grace <michael.grace@ury.org.uk>
Colin Roitt

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

console.log("Connecting...")

let scheme = window.location.protocol === "https:" ? "wss://" : "ws://"
let started = false;

const handleWs = () => {

    let ws = new WebSocket(scheme + window.location.host + "/board-ws");

    ws.onopen = () => {
        console.log("Connected.");
        document.getElementById("warning").hidden = true;
        if (!started) {
            ws.send("QUERY")
            started = true
        }
    }

    ws.onmessage = (event) => {

        let message = JSON.parse(event.data);

        if (message.action === "CLEAR") {
            document.getElementById("tweets").innerHTML = "";
            return;
        } else if (message.action === "REMOVE") {
            document.getElementById(message.id).remove()
            return
        }

        tweetHtml = atob(message.tweet.html);

        let newTweet = document.createElement("div");
        newTweet.innerHTML = tweetHtml;
        newTweet.id = message.tweet.id;
        newTweet.classList = "w-25 p-3";
        document.querySelector("tweets").prepend(newTweet);
        twttr.widgets.load(newTweet);
    };

    ws.onclose = () => {
        document.getElementById("warning").hidden = false;
        console.log("reconnecting in 1 sec")
        setTimeout(() => { handleWs() }, 1000)
    }

}

fetch("/info").then(d => d.json()).then(j => {
    document.getElementById("hashtag").innerHTML = j.hashtags.join(", ");
})

handleWs()