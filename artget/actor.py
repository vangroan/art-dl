
import asyncio
from asyncio import coroutine

class Actor:

    def __init__(self, timeout=10):
        self.inbox = asyncio.Queue()
        self._running = False
        self.timeout = timeout

    @property
    def running(self):
        return self._running

    @coroutine
    def _run(self):
        while self._running:
            msg = yield from asyncio.wait_for(self.inbox.get(), self.timeout)
            result = yield from self.on_message(msg)

    @coroutine
    def on_message(self, message):
        raise NotImplementedError()

    def stop(self):
        self._running = False

    @coroutine
    def start(self):
        self._running = True
        return asyncio.ensure_future(self._run())
