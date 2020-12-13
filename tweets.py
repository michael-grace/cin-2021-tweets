"""
    URY Tweet Board
    Candidate Interview Night 2021

    Author: Michael Grace
    Date: November, December 2020

    github.com/UniversityRadioYork

"""

from multiprocessing import Process
from tweet_http import http_server
from tweet_ws import ws_server
from tweet_twitter import start_recv_tweets

# HTTP Server
http_server: Process = Process(target=http_server)
http_server.start()

# WebSocket Server
ws_server: Process = Process(target=ws_server)
ws_server.start()

# Twitter Caller
twitter_caller: Process = Process(target=start_recv_tweets)
twitter_caller.start()

http_server.join()
ws_server.join()
twitter_caller.join()
