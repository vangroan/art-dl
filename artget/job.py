
class Job(object):

    def __init__(self, key):
        self.key = key
        self.retries = 3

    def retry(self):
        self.retries -= 1

    def __repr__(self):
        return '<{}: {}>'.format(self.__class__.__name__, self.key)

class GetJob(Job):

    def __init__(self, key, url):
        super(GetJob, self).__init__(key)
        self.url = url

class ParseJob(Job):

    def __init__(self, key, xml):
        super(ParseJob, self).__init__(key)
        self.xml = xml

class DownloadJob(Job):

    def __init__(self, key, filename, url):
        super(DownloadJob, self).__init__(key)
        self.filename = filename
        self.url = url
