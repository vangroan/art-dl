
import asyncio
from asyncio import Semaphore, coroutine

from aiohttp import ClientSession


class ThrottledClient(ClientSession):
    def __init__(self, loop, concurrent_requests):
        super().__init__(loop=loop)
        self._semaphore = Semaphore(concurrent_requests)

    @coroutine
    def get_throttled(self, url, timeout, headers=None):
        with (yield from self._semaphore):
            # compress=True breaks DrawCrowd image download
            return (yield from asyncio.wait_for(self.get(url, 
                                                headers=headers), timeout))

    @coroutine
    def throttled_content_read(self, response, queue):
        with (yield from self._semaphore):
            while True:
                chunk = yield from response.content.read(4 * 1024)
                if not chunk:
                    yield from queue.put(None)
                    break
                else:
                    yield from queue.put(chunk)
