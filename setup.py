
from setuptools import setup, find_packages

setup(
    name = 'artget',
    version = '0.0.1',
    packages = 'artget',
    package_dir = { 'artget' : 'src'},
    scripts = ['scripts/artget.py']
)
