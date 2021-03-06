/**
    URY Tweet Board
    Candidate Interview Night 2021

    Author: Michael Grace
    Date: November 2020

    github.com/UniversityRadioYork
 */

var xhttp = new XMLHttpRequest();
xhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
        server_data = JSON.parse(this.responseText);
        console.log("Connecting...")
        console.log(server_data.ws_conn);
        var ws = new WebSocket(server_data.ws_conn + "/client");
        ws.onopen = function() {
            console.log("Connected.");
        }
        ws.onmessage = function(event) {
            console.log(event.data);
        };
        ws.onclose = function() {
            console.log("Random Screaming!")
        }
    }
};
xhttp.open("GET", "/info", true);
xhttp.send();