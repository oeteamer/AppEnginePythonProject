import webapp2

from core import core_parser
from core import viewer


class Authors(webapp2.RequestHandler):
    def get(self):
        url = 'http://samlib.ru/a/aleksej_shpik/'

        (author, contents) = core_parser.AuthorsParser.indexdate_parser(url)

        writer = viewer.AuthorsWriter('index.html')
        template = writer.write_index_date(author, contents)

        self.response.headers['Content-Type'] = 'text/html'
        self.response.write(template)