
from argparse import ArgumentParser
import os, sys

from artget.core import Application


def default_output_directory():
    return os.path.join(os.environ['USERPROFILE'], 'Downloads', 'test')


def parse_args():
    parser = ArgumentParser(description='A deviantart webscraper')

    parser.add_argument(
            'galleries',
            nargs='+',
            help='Usernames of deviantart galleries'
        )

    parser.add_argument(
            '--concurrent', '-x',
            type=int,
            dest='concurrent',
            default=4
        )

    parser.add_argument(
            '--directory', '-d',
            dest='output_directory',
            default=default_output_directory(),
            help='Output directory'
        )

    parser.add_argument(
            '--include', '-i',
            dest='include',
            help='Text file containing list of galleries in include'
        )

    # TODO: Implement file download overwrite
    parser.add_argument(
            '--overwrite',
            action='store_true',
            default=False,
            help='Overwrite existing files'
        )

    parser.add_argument(
            '--sleep', '-s',
            dest='sleep',
            default=1
        )

    parser.add_argument(
            '--timeout',
            type=int,
            default=5,
            help='General timeout for requests'
        )

    parser.add_argument(
            '--debug',
            action='store_true',
            default=False,
        )

    return parser.parse_args()


def main():
    args = parse_args()

    app = Application(args)
    app.run()


if __name__ == '__main__':
    main()
