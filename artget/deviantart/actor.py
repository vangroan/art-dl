# src/deviantart.py

from actor import Actor
from http_client import AsyncHttpProvider

from asyncio import coroutine

class PageGetActor(Actor):

    def __init__(self):
        super(PageGetActor, self).__init__()

    @coroutine
    def on_message(self, msg)
        pass
