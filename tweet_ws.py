"""
    URY Tweet Board
    Candidate Interview Night 2021

    Author: Michael Grace
    Date: November, December 2020

    github.com/UniversityRadioYork

"""

import asyncio
import json
from typing import Any, Dict, Optional, Set

import websockets

import config

controller_connection: Optional[websockets.server.WebSocketServerProtocol] = None

pending_tweets: Dict[int, Dict[str, Any]] = {}
client_connections: Set[websockets.server.WebSocketServerProtocol] = set()


async def send_to_client(tweet: Dict[str, Any]) -> None:
    try:
        await asyncio.wait([conn.send(json.dumps(tweet)) for conn in client_connections])
    except ValueError:
        # No clients connected - Not a problem
        pass


async def recv_tweets(websocket: websockets.server.WebSocketServerProtocol) -> None:
    async for message in websocket:
        data = json.loads(message)
        pending_tweets[data["id"]] = data
        if controller_connection:
            await controller_connection.send(message)


async def recv_decisions(websocket: websockets.server.WebSocketServerProtocol) -> None:
    async for message in websocket:
        data = json.loads(message)
        if data["decision"] == "ACCEPT":
            await send_to_client(pending_tweets[data["id"]])
        try:
            del pending_tweets[data["id"]]
        except KeyError:
            pass


async def keep_client_alive(websocket: websockets.server.WebSocketServerProtocol) -> None:
    # Keep Connection Alive
    async for _ in websocket:
        pass


async def ws_tweets(websocket: websockets.server.WebSocketServerProtocol, path: str) -> None:
    global controller_connection
    print("Websocket Connected")

    try:
        await websocket.send("Hello There")

        if path == "/control":
            print("Controller Connection")
            controller_connection = websocket
            await recv_decisions(websocket)

        elif path == "/client":
            print("Client Connection")
            client_connections.add(websocket)
            await keep_client_alive(websocket)

        elif path == "/internal":
            print("Internal Connection")
            await recv_tweets(websocket)

    except websockets.exceptions.ConnectionClosedError:
        print("RIP Connection")
    finally:
        if websocket in client_connections:
            client_connections.remove(websocket)


def ws_server() -> None:
    print("Starting WebSocket Server")

    ws: websockets.server.Serve = websockets.serve(
        ws_tweets, config.HOST, config.WS_PORT)
    asyncio.get_event_loop().run_until_complete(ws)
    asyncio.get_event_loop().run_forever()


if __name__ == "__main__":
    print("Don't do this")
