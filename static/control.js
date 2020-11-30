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
        document.getElementById("hashtag").innerHTML = server_data.hashtag;

        // WebSocket Connection
        var ws = new WebSocket(server_data.ws_conn);
        var alert = document.getElementById("server");
        ws.onopen = function() {
            alert.innerText = "Connected to Server";
            alert.classList.remove("alert-info");
            alert.classList.remove("alert-danger");
            alert.classList.add("alert-success");
        };
        ws.onclose = function() {
            alert.innerText = "Disconnected from Server";
            alert.classList.remove("alert-info");
            alert.classList.remove("alert-success");
            alert.classList.add("alert-danger");
        };
        ws.onmessage = function(event) {
            console.log(event.data);
        }

    }
};
xhttp.open("GET", "/info", true);
xhttp.send();