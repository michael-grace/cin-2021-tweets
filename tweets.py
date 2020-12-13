"""
    URY Tweet Board
    Candidate Interview Night 2021

    github.com/UniversityRadioYork

    Copyright (C) 2020 Michael Grace/University Radio York
    michael.grace@ury.org.uk

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
