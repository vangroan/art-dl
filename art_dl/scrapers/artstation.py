
from asyncio import coroutine, sleep
from collections import namedtuple
import os
from sys import maxsize

from art_dl.scraper import Scraper, ScrapingException
from art_dl.util import check_or_make_dir, filename_from_url


class ArtstationScraper(Scraper):

    def __init__(self, http_client, logger, username, out_dir, overwrite):
        super().__init__(http_client, logger, overwrite)
        self.username = username
        self.out_dir = out_dir

    @classmethod
    def create_scraper(cls, ctx, username):
        return cls(ctx['http_client'], ctx['logger'], username,
                   ctx['output_directory'], ctx['overwrite'])

    @property
    def artist_dir(self):
        return os.path.join(self.out_dir, self.username)

    def projects_url(self, page_num):
        return 'https://www.artstation.com/users/{username}/projects.json?page={page}'.format(
            username=self.username,
            page=page_num,
        )

    def project_url(self, project):
        return 'https://www.artstation.com/projects/{hash_id}.json'.format(hash_id=project.hash_id)

    @coroutine
    def fetch_projects(self):

        last_length = maxsize
        page = 1
        projects = []

        while last_length > 0:
            response = yield from self.get(self.projects_url(page))

            data = yield from response.json()
            for project in data['data']:
                projects.append(ArtstationScraper.Project(
                    project['hash_id'],
                    project['title'].strip()
                ))

            response.close()
            page += 1
            last_length = len(data['data'])

        return projects

    @coroutine
    def fetch_project_image_url(self, project):

        image_urls = []

        url = self.project_url(project)
        response = yield from self.get(url)
        data = yield from response.json()

        for asset in data['assets']:
            image_urls.append(asset['image_url'])

        response.close()
        return image_urls

    @coroutine
    def run(self):

        check_or_make_dir(self.artist_dir)
        projects = yield from self.fetch_projects()

        for project in projects:
            self.info(project.title)
            for image_url in (yield from self.fetch_project_image_url(project)):
                filename = filename_from_url(image_url)
                filepath = os.path.join(self.artist_dir, filename)
                yield from self.download(image_url, filepath, self.overwrite)

        sleep(0.1)

    Project = namedtuple('Project', ['hash_id', 'title'])
