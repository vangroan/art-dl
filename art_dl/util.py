
import os


def check_or_make_dir(path):
    if not os.path.isdir(path):
        os.makedirs(path)


def filename_from_url(url):
    url = url.split('?')[0]
    return url.split('/')[-1]
