/**
    URY Tweet Board
    Candidate Interview Night 2022

    Author: Michael Grace, Colin Roitt
    Date: November 2020, January 2022

    github.com/UniversityRadioYork
 */

console.log("Connecting...")

let scheme = window.location.protocol === "https:" ? "wss://" : "ws://"

const handleWs = () => {

    let ws = new WebSocket(scheme + window.location.host + "/board-ws");

    ws.onopen = () => {
        console.log("Connected.");
        document.getElementById("warning").hidden = true;
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