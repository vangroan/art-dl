
from setuptools import setup, find_packages

setup(
    name = 'artget',
    version = '2016.01.09',
    packages = ['artget', 'artget.scrapers'],
    install_requires = [
        'aiohttp',
        'beautifulsoup4'
    ],
    entry_points = {
        'console_scripts' : ['artget = artget.cli:main']
    }
)
