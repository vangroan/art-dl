
from asyncio import coroutine
import os

from art_dl.scraper import Scraper


class TumblrScraper(Scraper):

    _userurl = 'characterdesigninspiration.tumblr.com'

    def __init__(self, http_client, logger, username, out_dir, overwrite):
        super().__init__(http_client, logger, overwrite)
        self.username = username
        self.out_dir = out_dir

        self.debug('Initialized')
        self.debug('Out directory: ' + self.tumblr_dir)

    @staticmethod
    def create_scraper(ctx, username):
        return TumblrScraper(ctx['http_client'], ctx['logger'], username, ctx['output_directory'], ctx['overwrite'])

    def tumblr_dir(self):
        return os.path.join(self.out_dir, self.username)

    @coroutine
    def run(self):
        pass
