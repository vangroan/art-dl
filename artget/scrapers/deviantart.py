
from asyncio import coroutine, sleep
from collections import namedtuple
import os
from urllib.parse import urlencode
from xml.etree import ElementTree

from bs4 import BeautifulSoup

from artget.scraper import Scraper, ScrapingException
from artget.util import check_or_make_dir, filename_from_url


class DeviantartScraper(Scraper):

    _config_dir_name = '.artdl'
    _host_url = 'http://www.deviantart.com'
    _rss_url = 'http://backend.deviantart.com/rss.xml'
    _rss_namespaces = {
        'media': 'http://search.yahoo.com/mrss/',
        'atom': 'http://www.w3.org/2005/Atom'
    }

    def __init__(self, http_client, username, out_dir):
        super().__init__(http_client)
        self.username = username
        self.out_dir = out_dir

    @staticmethod
    def create_scraper(ctx, username):
        return DeviantartScraper(ctx['http_client'], username, ctx['output_directory'])

    @staticmethod
    def create_scraper_for_gallery(ctx, username, gallery):
        # TODO: Set up scraper so that it downloads only a specified gallery
        pass

    def get_query_string(self):
        return 'gallery:{}'.format(self.username)

    @property
    def cache_dir(self):
        return os.path.join(self.out_dir, self._config_dir_name, 'cache')

    @property
    def ignore_file_name(self):
        return os.path.join(self.cache_dir, self.username)

    @property
    def deviant_dir(self):
        return os.path.join(self.out_dir, self.username)

    @property
    def rss_url(self):
        return self._rss_url + '?' + urlencode(
                {
                    'type': 'deviation',
                    'q': self.get_query_string()
                }
        )

    @coroutine
    def fetch_rss(self):
        status, body = yield from self.get_body(self.rss_url)
        # TODO: Handle failure status
        return body

    def scrape_deviations_list(self, rss_xml):
        """Using the xml from the rss feed, return lists
        of deviations"""

        tree = ElementTree.fromstring(rss_xml)
        item_nodes = tree.find('channel').findall('item')
        deviations = []

        for item_node in item_nodes:
            dev = DeviantartScraper.DeviationPage(
                    self.username,
                    item_node.find('guid').text,
                    item_node.find('link').text,
            )
            deviations.append(dev)

        return deviations

    @coroutine
    def fetch_deviation_page(self, url):
        status, body = yield from self.get_body(url)
        # TODO: Handle failure status
        return body

    def scrape_deviation_image_url(self, deviation_guid, dev_page_html):
        soup = BeautifulSoup(dev_page_html, 'html.parser')
        img_nodes = soup.select('img .dev-content-full')
        if not img_nodes:
            raise ScrapingException('Could not find image url on deviant page [%s]' % deviation_guid)
        return img_nodes[0]['src']

    @coroutine
    def download_deviation(self, image_url, image_filename):
        filepath = os.path.join(self.deviant_dir, image_filename)
        yield from self.download(image_url, filepath)

    @coroutine
    def run(self):

        check_or_make_dir(self.deviant_dir)

        # First get the rss feed which lists the deviantions
        rss_xml = yield from self.fetch_rss()

        # Visit each deviation serially and get the page html
        for dev in self.scrape_deviations_list(rss_xml):
            dev_page_html = yield from self.fetch_deviation_page(dev.url)

            image_url = self.scrape_deviation_image_url(dev.guid, dev_page_html)
            image_filename = filename_from_url(image_url)

            yield from self.download_deviation(image_url, image_filename)

        yield from sleep(0.001)

    DeviationPage = namedtuple('DeviationPage', ['username', 'guid', 'url'])