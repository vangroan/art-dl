
from asyncio import coroutine
import aiohttp

class HttpProvider:

    @coroutine
    def get(self, url, headers={}):
        raise NotImplementedError()

class AsyncHttpProvider(HttpProvider):

    @coroutine
    def get(self, url, headers={}):
        res = yield from aiohttp.request(url, headers=headers)
        return res
