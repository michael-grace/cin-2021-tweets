"""
    URY Tweet Board
    Candidate Interview Night 2021

    Author: Michael Grace
    Date: November, December 2020

    github.com/UniversityRadioYork

"""

from multiprocessing import Process

from tweet_http import http_server
from tweet_twitter import start_recv_tweets
from tweet_ws import ws_server


def tweet_board() -> None:
    # HTTP Server
    http_process: Process = Process(target=http_server)
    http_process.start()

    # WebSocket Server
    ws_process: Process = Process(target=ws_server)
    ws_process.start()

    # Twitter Caller
    twitter_process: Process = Process(target=start_recv_tweets)
    twitter_process.start()

    http_process.join()
    ws_process.join()
    twitter_process.join()


if __name__ == "__main__":
    tweet_board()
