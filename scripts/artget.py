
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
        '--directory', '-d',
        dest='output_directory',
        default=default_output_directory(),
        help='Output directory'
    )

    parser.add_argument(
        '--workers', '-w',
        dest='workers',
        default=15
    )

    parser.add_argument(
        '--sleep', '-s',
        dest='sleep',
        default=0.5
    )

    return parser.parse_args()

def main():
    args = parse_args()

    app = Application(args)
    app.run()

if __name__ == '__main__':

    main()
