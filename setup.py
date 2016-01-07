
from setuptools import setup, find_packages

setup(
    name = 'artget',
    version = '0.0.2b',
    packages = ['artget', 'artget.scrapers'],
    install_requires = [
        'aiohttp',
        'beautifulsoup4'
    ],
    entry_points = {
        'console_scripts' : ['artget = artget.cli:main']
    }
)
