# Restructuring scraper
#
# Notes:
#   Coroutines for requesting an html page start with 'fetch_*'
#   Coroutines for streaming larger files start with 'download_*'
#   Methods for scraping pages start with 'scrape_*'
#   The concrete scraper is responsible for choosing it's parsing tools

import asyncio
from asyncio import coroutine
from collections import namedtuple
import os
from shutil import move


class ScrapingException(Exception): pass


# TODO: Save urls to file for faster skipping of seen pages
# TODO: Use the imghdr module to guess image filetype for files with no extension
class Scraper:

    def __init__(self, http_client):
        self.client = http_client

    @coroutine
    def get(self, url, timeout=120, headers=None):
        return (yield from self.client.get_throttled(
                                url, 
                                timeout=timeout, 
                                headers=headers))

    @coroutine
    def get_body(self, url):
        response = yield from self.get(url)
        status = response.status
        body = yield from response.read()
        response.close()
        return status, body

    @coroutine
    def download(self, url, target_file, overwrite=False):

        if os.path.exists(target_file) and not overwrite:
            return

        response = yield from self.get(url)
        partial_file = target_file + '.part'

        with open(partial_file, 'wb') as fp:
            chunk_queue = asyncio.Queue()
            yield from self.client.throttled_content_read(response, chunk_queue)
            while True:
                chunk = yield from chunk_queue.get()
                if not chunk:
                    break
                fp.write(chunk)

        move(partial_file, target_file)
        response.close()

    @coroutine
    def run(self):
        raise NotImplementedError('run() is not implemented for %s' % self)

    def __repr__(self):
        return '<%s>' % type(self).__name__ 

ScraperResponse = namedtuple('ScraperResponse', ['scraper', 'task'])
