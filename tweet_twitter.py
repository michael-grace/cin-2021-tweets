"""
    URY Tweet Board
    Candidate Interview Night 2021

    Author: Michael Grace
    Date: December 2020

    github.com/UniversityRadioYork

"""


import asyncio
import json
import random
from typing import Any, Dict

import websockets
from twitter import Api, error

import config


async def recv_tweets():
    api: Api = Api(
        config.TWITTER_API_KEY,
        config.TWITTER_API_SECRET_KEY,
        config.TWITTER_ACCESS_TOKEN,
        config.TWITTER_ACCESS_TOKEN_SECRET,
        tweet_mode="extended"
    )
    while True:
        async with websockets.connect(f"ws://{config.HOST}:{config.WS_PORT}/internal") as websocket:
            try:
                stream: Any[Dict[str, Any]] = api.GetStreamFilter(
                    track=[config.HASHTAG], languages=["en"])

                for tweet in stream:
                    if tweet["text"][:2] == "RT":
                        if "extended_tweet" in tweet["retweeted_status"].keys():
                            body = "{0}: {1}".format(
                                tweet["text"].split(":")[0],
                                tweet["retweeted_status"]["extended_tweet"]["full_text"]
                            )
                        else:
                            body = tweet["text"]

                    else:
                        if "extended_tweet" in tweet.keys():
                            body = tweet["extended_tweet"]["full_text"]
                        else:
                            body = tweet["text"]

                    tweet_info = {
                        "id": random.randint(0, 10000000000),
                        "title": "{0} - @{1}".format(tweet["user"]["name"], tweet["user"]["screen_name"]),
                        "body": body
                    }

                    await websocket.send(json.dumps(tweet_info))

            except error.TwitterError as e:
                print(f"Twitter Error: Probably Rate Limiting: {e}")


def start_recv_tweets():
    asyncio.get_event_loop().run_until_complete(recv_tweets())


if __name__ == "__main__":
    print("Don't do this")
