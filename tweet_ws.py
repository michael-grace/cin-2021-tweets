"""
    URY Tweet Board
    Candidate Interview Night 2021

    Author: Michael Grace
    Date: November 2020

    github.com/UniversityRadioYork

"""

import config
import asyncio
import websockets
import json
from twitter import Api

async def recv_tweets(websocket):
    api = Api(
    config.TWITTER_API_KEY,
    config.TWITTER_API_SECRET_KEY,
    config.TWITTER_ACCESS_TOKEN,
    config.TWITTER_ACCESS_TOKEN_SECRET,
    tweet_mode="extended"
    )
    for tweet in api.GetStreamFilter(track=[config.HASHTAG], languages=["en"]):
        if tweet["text"][:2] == "RT":
            if "extended_tweet" in tweet["retweeted_status"].keys():
                body = "{0}: {1}".format(tweet["text"].split(":")[0], tweet["retweeted_status"]["extended_tweet"]["full_text"])
            else:
                body = tweet["text"]
        else:
            if "extended_tweet" in tweet.keys():
                body = tweet["extended_tweet"]["full_text"]
            else:
                body = tweet["text"]
        await websocket.send(json.dumps({
            "title": "{0} - @{1}".format(tweet["user"]["name"], tweet["user"]["screen_name"]),
            "body": body
        }))


async def ws_tweets(websocket, path):
    print("Websocket Connected")
    try:
        await websocket.send("Hello There")
        await asyncio.gather(recv_tweets(websocket))
    except websockets.exceptions.ConnectionClosedError:
        print("RIP Connection")

def ws_server() -> None:
    print("Starting WebSocket Server")
    ws = websockets.serve(ws_tweets, config.HOST, config.WS_PORT)
    asyncio.get_event_loop().run_until_complete(ws)
    asyncio.get_event_loop().run_forever()

if __name__ == "__main__":
    print("Don't do this")