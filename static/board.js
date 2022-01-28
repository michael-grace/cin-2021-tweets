/**
    URY Tweet Board
    Candidate Interview Night 2022

    Author: Michael Grace, Colin Roitt
    Date: November 2020, January 2022

    github.com/UniversityRadioYork
 */

console.log("Connecting...")

let scheme = window.location.protocol === "https:" ? "wss://" : "ws://"
let ws = new WebSocket(scheme + window.location.host + "/board-ws");

ws.onopen = function() {
    console.log("Connected.");
}

ws.onmessage = function(event) {

    tweetJson = JSON.parse(event.data);
    tweetHtml = atob(tweetJson.html);
    console.log(tweetHtml)

    let newTweet = document.createElement("div");
    newTweet.innerHTML = tweetHtml;
    newTweet.classList = "w-25 p-3";
    document.querySelector("tweets").prepend(newTweet);
    twttr.widgets.load(newTweet);
};

ws.onclose = function() {
    // TODO
}