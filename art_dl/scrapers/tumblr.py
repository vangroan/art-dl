
from asyncio import coroutine
from collections import namedtuple
import os
import re

from art_dl.scraper import Scraper, ScrapingException
from art_dl.util import check_or_make_dir, filename_from_url


class TumblrScraper(Scraper):

    _userurl = 'http://{username}.tumblr.com'
    _pageurl = 'http://{username}.tumblr.com/page/{pagenum}/'
    _postpattern = r'(http://{username}\.tumblr\.com/post/\d+/[a-zA-Z0-9-]+)'
    _imgpattern = r'(http://\d+\.media\.tumblr\.com/[a-zA-Z0-9]+/[a-zA-Z0-9_]+\.[a-zA-Z]+)'

    def __init__(self, http_client, logger, username, out_dir, overwrite):
        super().__init__(http_client, logger, overwrite)
        self.username = username
        self.out_dir = out_dir

        self.debug('Initialized')
        self.debug('Out directory: ' + self.tumblr_dir)

    @staticmethod
    def create_scraper(ctx, username):
        return TumblrScraper(ctx['http_client'], ctx['logger'], username, ctx['output_directory'], ctx['overwrite'])

    @property
    def tumblr_dir(self):
        return os.path.join(self.out_dir, self.username)

    @property
    def postpattern(self):
        return self._postpattern.format(username=self.username)

    @coroutine
    def fetch_page(self, pagenum):
        url = self._pageurl.format(username=self.username, pagenum=pagenum)
        status, body = yield from self.get_body(url)
        return body.decode('utf-8')

    @coroutine
    def fetch_post(self, url):
        status, body = yield from self.get_body(url)
        return body.decode('utf-8')

    def scrape_posts(self, page_html):
        posts = []
        for match in re.finditer(self.postpattern, page_html):
            posts.append(TumblrScraper.TumblrPost(
                self.username,
                match.group(0)
            ))
        return posts

    def scrape_images(self, post_html):
        for match in re.finditer(self._imgpattern, post_html):
            yield match.group(0)

    @coroutine
    def download_image(self, url, image_filename):
        file_path = os.path.join(self.tumblr_dir, image_filename)
        yield from self.download(url, file_path, self.overwrite)

    @coroutine
    def run(self):

        check_or_make_dir(self.tumblr_dir)

        pagenum = 1
        while True:
            page_html = yield from self.fetch_page(pagenum)

            posts = self.scrape_posts(page_html)
            if not posts:
                # No posts were found. We've probably reached an out of range page
                break

            for post in posts:
                post_html = yield from self.fetch_post(post.url)
                for image_url in self.scrape_images(post_html):
                    file_name = filename_from_url(image_url)
                    yield from self.download_image(image_url, file_name)

            pagenum += 1

    TumblrPost = namedtuple('TumblrPost', ['username', 'url'])
