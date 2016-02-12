
from asyncio import coroutine
from collections import namedtuple
import os
from sys import maxsize

from art_dl.scraper import Scraper
from art_dl.util import check_or_make_dir, filename_from_url


class DrawcrowdScraper(Scraper):

    def __init__(self, http_client, username, out_dir):
        super().__init__(http_client)
        print('Drawcrowd username: %s' % username)
        self.username = username
        self.out_dir = out_dir

    @classmethod
    def create_scraper(cls, ctx, username):
        return cls(
                ctx['http_client'],
                username,
                ctx['output_directory'])

    @property
    def project_dir(self):
        return os.path.join(self.out_dir, self.username)

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
        headers = {
            'Accept' : 'application/json'
        }
        projects = []
        last_length = maxsize

        while last_length > 0:
            print('Getting project list offset %d limit %d last_length %d' % (offset, limit, last_length))
            url = self.projects_url(offset, limit)
            response = yield from self.get(url, headers=headers)

            # TODO: Handle failure status code
            data = yield from response.json()
            for project in data['projects']:
                projects.append(DrawcrowdScraper.DrawcrowdProject(
                    project['slug'],
                    project['title'],
                    project['original_image']
                ))

            response.close()

            offset += 50

            last_length = int(data['meta']['length'])

        return projects

    @coroutine
    def run(self):

        check_or_make_dir(self.project_dir)
        projects = yield from self.fetch_projects()

        for project in projects:
            image_url = project.original_image
            filename = filename_from_url(image_url)
            file_path = os.path.join(self.project_dir, filename)
            os.path.join(self.project_dir, file_path)
            yield from self.download(image_url, file_path)

    DrawcrowdProject = namedtuple('DrawcrowdProject', ['slug', 'title', 'original_image'])