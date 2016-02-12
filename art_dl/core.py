
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
        self.logger = self.create_logger('art-dl')

        configure_rules(self.rules)

    @staticmethod
    def _create_context_processor(http_client, logger, output_directory):
        def context_processor(ctx):
            ctx['http_client'] = http_client
            ctx['output_directory'] = output_directory
            ctx['logger'] = logger

        return context_processor

    @staticmethod
    def create_logger(name):
        logger = logging.getLogger(name)
        logger.setLevel(logging.DEBUG)

        ch = logging.StreamHandler()
        ch.setLevel(logging.DEBUG)

        fmt = logging.Formatter('%(asctime)s (%(name)s) [%(levelname)s] %(message)s')
        ch.setFormatter(fmt)

        logger.addHandler(ch)

        return logger

    @staticmethod
    def log_config(logger, config):
        s = ['Config:\n']
        for k, v in vars(config).items():
            s.append('\t')
            s.append(k)
            s.append(': ')
            s.append(str(v))
            s.append('\n')
        logger.debug(''.join(s))

    def run(self):

        # TODO: Better logging
        if self.config.debug:
            logging.basicConfig(level=logging.DEBUG)

        loop = get_event_loop()
        self.log_config(self.logger, self.config)
        loop.set_debug(self.config.debug)
        client = ThrottledClient(loop, self.config.concurrent)
        tasks = None

        try:
            # TODO: Needs readability
            scrapers = [self.rules.dispatch(gallery, context_processor=self._create_context_processor(
                client, self.create_logger(gallery), self.config.output_directory))
                for gallery in self.config.galleries]

            tasks = asyncio.gather(*(s.run() for s in scrapers))
            loop.run_until_complete(tasks)

        except KeyboardInterrupt:
            self.logger.info('Shutting down...')

            self.logger.info('Cancelling Tasks')
            all_tasks = asyncio.gather(*asyncio.Task.all_tasks())
            all_tasks.cancel()

            self.logger.info('Restarting loop')
            loop.run_forever()

            if tasks:
                tasks.exception()  # Avoid warning for not fetching Exceptions
            all_tasks.exception()

            self.logger.info('Done')
        finally:
            client.close()
            loop.close()
