
import os
from urllib.parse import urlencode

from .util import check_or_make_dir

class DeviantartScraper(object):

    def __init__(self, username, host_url, out_directory):

        self.username = username
        self.host_url = host_url
        self.out_directory = os.path.join(out_directory, username)

        check_or_make_dir(self.out_directory)

    def get_query_string(self):
        return 'gallery:{}'.format(self.username)

    def get_rss_url(self):
        url = self.host_url + '?' + urlencode(
            {
                'type' : 'deviation',
                'q' : self.get_query_string()
            }
        )
        return url

    def get_image_filepath(self, filename):
        return os.path.join(self.out_directory, filename)

    def image_exists(self, filename):
        return os.path.exists(self.get_image_filepath(filename))

