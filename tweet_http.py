"""
    URY Tweet Board
    Candidate Interview Night 2021

    Author: Michael Grace
    Date: November 2020

    github.com/UniversityRadioYork

"""

from flask import Flask, url_for, redirect
import config

app = Flask(__name__)

@app.route("/")
def base():
    return "", 400

@app.route("/board")
def board():
    return redirect(url_for("static", filename="board.html"))

@app.route("/graphic")
def graphic():
    return redirect(url_for("static", filename="graphic.html"))

@app.route("/control")
def control():
    return redirect(url_for("static", filename="control.html"))

def http_server() -> None:
    print("Starting HTTP Server")
    app.run(host=config.HOST, port=config.PORT, debug=True)

if __name__ == "__main__":
    print("Don't do this")