/**
    URY Tweet Board
    Candidate Interview Night 2021

    Author: Michael Grace
    Date: November 2020

    github.com/UniversityRadioYork
 */

// Put's the appropriate hashtag in
var xhttp = new XMLHttpRequest();
xhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
        document.getElementById("hashtag").innerHTML =
            this.responseText;
    }
};
xhttp.open("GET", "/hashtag", true);
xhttp.send();