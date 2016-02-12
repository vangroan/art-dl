
import asyncio
from asyncio import get_event_loop
import logging

from art_dl.client import ThrottledClient
from art_dl.rulematch import PatternRules
from art_dl.rules import configure_rules


class Application(object):

    def __init__(self, config):

        self.config = config
        self.rules = PatternRules()

        configure_rules(self.rules)

    @staticmethod
    def _create_context_processor(http_client, output_directory):
        def context_processor(ctx):
            ctx['http_client'] = http_client
            ctx['output_directory'] = output_directory
        return context_processor

    def run(self):

        # TODO: Better logging
        if self.config.debug:
            logging.basicConfig(level=logging.DEBUG)

        loop = get_event_loop()
        print(self.config)
        loop.set_debug(self.config.debug)
        client = ThrottledClient(loop, self.config.concurrent)
        tasks = None

        try:
            context_processor = self._create_context_processor(client, self.config.output_directory)
            scrapers = [self.rules.dispatch(gallery, context_processor=context_processor) 
                        for gallery in self.config.galleries]

            tasks = asyncio.gather(*(s.run() for s in scrapers))
            loop.run_until_complete(tasks)

        except KeyboardInterrupt:
            print('Shutting down...')

            print('Cancelling Tasks')
            all_tasks = asyncio.gather(*asyncio.Task.all_tasks())
            all_tasks.cancel()

            print('Restarting loop')
            loop.run_forever()

            if tasks:
                tasks.exception()  # Avoid warning for not fetching Exceptions
            all_tasks.exception()

            print('Done')
        finally:
            client.close()
            loop.close()
