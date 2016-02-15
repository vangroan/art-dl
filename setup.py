from setuptools import setup, find_packages

setup(
    name='art-dl',
    version='2016.02.15',
    packages=['art_dl', 'art_dl.scrapers'],
    install_requires=[
        'aiohttp',
        'beautifulsoup4'
    ],
    entry_points={
        'console_scripts': ['art-dl = art_dl.cli:main']
    }
)
