
from asyncio import coroutine
from collections import namedtuple

from artget.scraper import Scraper


class DrawcrowdScraper(Scraper):

    def __init__(self, http_client, username, out_dir):
        super().__init__(http_client)

        self.username = username
        self.out_dir = out_dir

    @classmethod
    def create_scraper(cls, ctx, username):
        return cls(
                ctx['http_client'],
                username,
                ctx['output_directory'])

    def projects_url(self, offset, limit):
        return 'http://drawcrowd.com/{username}/projects?offset={offset}&sort=newest&limit={limit}'.format(
                username=self.username,
                offset=offset,
                limit=limit,
            )

    @coroutine
    def fetch_projects(self):
        offset = 0
        limit = 50
        url = self.projects_url(offset, limit)
        headers = {
            'Accept' : 'application/json'
        }
        response = yield from self.get(url, headers=headers)
        # TODO: Handle failure status code
        print((yield from response.read()))
        data = yield from response.json()
        response.close()

        print(data)

    @coroutine
    def run(self):
        yield from self.fetch_projects()

    DrawcrowdProject = namedtuple('DrawcrowdProject', ['slug', 'title', 'original_image'])