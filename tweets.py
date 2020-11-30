"""
    URY Tweet Board
    Candidate Interview Night 2021

    Author: Michael Grace
    Date: November 2020

    github.com/UniversityRadioYork

"""

from multiprocessing import Process
from tweet_http import http_server
from tweet_ws import ws_server

# HTTP Server
http_server: Process = Process(target=http_server)
http_server.start()
http_server.join()

# WebSocket Server
ws_server: Process = Process(target=ws_server)
ws_server.start()
ws_server.join()