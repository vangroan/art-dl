
import os
from xml.etree import ElementTree as ET
import logging
logging.basicConfig(level=logging.DEBUG)

import asyncio
from asyncio import Queue, QueueEmpty
import aiohttp

from .job import GetJob, ParseJob, DownloadJob
from .util import filename_from_url

def dump_tree(el, level=0):
    print('{}{} "{}"'.format(''.join(' '*level), el.tag, el.text))
    for child in el:
        dump_tree(child, level=level+1)

class Worker(object):

    _id = 0

    def __init__(self, app):

        self.id = Worker.newid()
        self.name = 'Worker-{}'.format(self.id)
        self.running = False
        self.app = app

    @classmethod
    def newid(cls):
        _id = cls._id
        cls._id += 1
        return _id

    def stop(self):
        self.running = False

    @asyncio.coroutine
    def run_get(self):
        logging.info('Starting %s as Getter' % self)
        self.running = True

        while self.running:

            try:

                job = self.app.get_queue.get_nowait()

                logging.debug(job)

                logging.info('{}: GET {}'.format(self, job.url))
                response = yield from aiohttp.request('GET', job.url)
                logging.info('{}: {} {}'.format(self, response.status, job.url))

                if response.status >= 400:
                    if job.retries > 0:
                        job.retry()
                        self.app.get_queue.put_nowait(job)
                        continue
                    else:
                        # Job failed
                        logging.warning('{}: {} job failed'.format(self, job))
                        continue

                body = yield from response.read()

                self.app.parse_queue.put_nowait(ParseJob(job.key, body))

            except QueueEmpty:
                yield from asyncio.sleep(self.app.config.sleep)

    @asyncio.coroutine
    def run_parse(self):
        logging.info('Starting %s as Parser' % self)
        self.running = True

        while self.running:
            try:

                job = self.app.parse_queue.get_nowait()

                logging.debug('{}: {}'.format(self, job))

                tree = ET.fromstring(job.xml)
                images = tree.findall("channel/item/{http://search.yahoo.com/mrss/}content[@medium='image']")

                scraper = self.app.get_scraper(job.key)

                found = False
                for image in images:
                    url = image.get('url')
                    filename = filename_from_url(url)
                    if scraper.image_exists(filename):
                        logging.debug('{} seen'.format(filename))
                        continue
                    if self.app.seen(url):
                        continue
                    self.app.add_seen(url)
                    yield from self.app.download_queue.put(DownloadJob(job.key, filename, url))
                    found = True

                if not found:
                    logging.warning('{}: No images found for {}'.format(self, job))

            except QueueEmpty:

                yield from asyncio.sleep(self.app.config.sleep)

    @asyncio.coroutine
    def run_download(self):
        logging.info('Starting %s as Downloader' % self)
        running = True

        while running:
            try:

                job = self.app.download_queue.get_nowait()

                r = yield from aiohttp.request('GET', job.url)

                if r.status >= 400:
                    if job.retries > 0:
                        job.retry()
                        self.app.download_queue.put_nowait(job)
                        continue
                    else:
                        # Job failed
                        logging.warning('{}: {} job failed'.format(self, job))
                        continue

                scraper = self.app.get_scraper(job.key)
                filepath = scraper.get_image_filepath(job.filename)
                temp_filepath = filepath + '.part'
                with open(temp_filepath, 'wb') as fp:
                    while True:
                        chunk = yield from r.content.read(4 * 1024)
                        if not chunk:
                            break
                        fp.write(chunk)

                os.rename(temp_filepath, filepath)

            except QueueEmpty:
                yield from asyncio.sleep(self.app.config.sleep)

    def __repr__(self):
        return '<{}>'.format(self.name)