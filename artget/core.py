
import asyncio
from asyncio import Queue
import os
import signal

from .scraper import DeviantartScraper
from .worker import Worker
from .job import GetJob
from .util import check_or_make_dir

META_DIR = '.grab'
CLOSED_FILE = 'closed'
DEVIANTART_RSS_URL = 'http://backend.deviantart.com/rss.xml'

NUM_GET_WORKERS = 4
NUM_PARSE_WORKERS = 1

class Application(object):

    def __init__(self, config):

        self.config = config
        self.scrapers = dict()
        self.closed_set = set()

        self.get_queue = Queue()
        self.parse_queue = Queue()
        self.download_queue = Queue()

        self.meta_dir = os.path.join(self.config.output_directory, META_DIR)
        check_or_make_dir(self.meta_dir)
        self.closed_file_path = os.path.join(self.meta_dir, CLOSED_FILE)

        self.workers = []

        self.create_scrapers()

    def seen(self, url):
        return url in self.closed_set

    def add_seen(self, url):
        self.closed_set.add(url)

    def save_closed(self):
        with open(self.closed_file_path, 'w') as fp:
            for url in self.closed_set:
                fp.write(url)

    def load_closed(self):
        if os.path.isfile(self.closed_file_path):
            with open(self.closed_file_path, 'r') as fp:
                for line in fp:
                    self.closed_set.add(line)


    def get_scraper(self, key):
        return self.scrapers[key]

    def build_sites(self):
        for username in self.config.galleries:
            pass

    def seed_url(self):
        pass

    def create_scrapers(self):

        for username in self.config.galleries:
            self.scrapers[username] = DeviantartScraper(username, DEVIANTART_RSS_URL, self.config.output_directory)

    def stop(self):
        for w in self.workers:
            w.stop()

    def run(self):

        # Seed urls
        for key, scraper in self.scrapers.items():
            self.get_queue.put_nowait(GetJob(key, scraper.get_rss_url()))

        tasks = []

        # Get XML workers
        for i in range(NUM_GET_WORKERS):
            w = Worker(self)
            self.workers.append(w)
            tasks.append(asyncio.async(w.run_get()))

        # Download workers
        for i in range(NUM_PARSE_WORKERS):
            w = Worker(self)
            self.workers.append(w)
            tasks.append(asyncio.async(w.run_parse()))

        # Download workers
        for i in range(self.config.workers):
            w = Worker(self)
            self.workers.append(w)
            tasks.append(asyncio.async(w.run_download()))

        loop = asyncio.get_event_loop()

        for signame in ('SIGINT', 'SIGTERM'):
            # NotImplemented on Windows
            #loop.add_signal_handler(getattr(signal, signame),
            #                        lambda: self.stop())
            pass

        try:
            loop.run_until_complete(asyncio.wait(tasks))
        except KeyboardInterrupt:
            print('Shutting down...')
            self.stop()
        finally:
            loop.close()
            pass




