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

async def ws_tweets(websocket, path):
    print("Websocket Connected")
    try:
        for _ in range(5):
            await websocket.send(json.dumps({"title": "Tweet Title - @tweet", "body": "Hello, this is a tweet, #CIN21"}))
        async for message in websocket:
            await websocket.send("Hello There")
    except websocket.WebSocketException.ClosedConnection:
        print("RIP Connection")

def ws_server() -> None:
    print("Starting WebSocket Server")
    ws = websockets.serve(ws_tweets, config.HOST, config.WS_PORT)
    asyncio.get_event_loop().run_until_complete(ws)
    asyncio.get_event_loop().run_forever()

if __name__ == "__main__":
    print("Don't do this")