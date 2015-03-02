import webapp2
import models

from core import core_parser
from core import viewer
from datetime import datetime
from google.appengine.api import taskqueue


class Index(webapp2.RequestHandler):
    def get(self):
        authors_model = models.Authors.query().fetch()
        authors = []
        for model in authors_model:
            authors.append({
                'href': '/author/'+model.key.id(),
                'href_samlib': core_parser.AuthorsParser.get_author_link(model.key.id()),
                'name': model.name
            })

        self.response.headers['Content-Type'] = 'text/html'
        self.response.write(viewer.AuthorsWriter('index.html').write(authors))


class Authors(webapp2.RequestHandler):
    def get(self, author_code):
        author_model = models.Authors.get_by_id(author_code)
        if not author_model:
            author_model = models.Authors.create({
                'code': author_code,
                'name': core_parser.AuthorsParser.get_author_name_from_url(author_code)}
            ).put()

        author = {'code': author_model.key.id(), 'name': author_model.name}

        contents = models.AuthorsBooks.query(ancestor=models.authors_key(author['code']))\
            .order(models.AuthorsBooks.updated_at)\
            .fetch()

        parse_contents = core_parser.AuthorsParser.indexdate_parser(author['code'])

        for parse_item in parse_contents:
            book_exsist = False
            for item in contents:
                if parse_item['id'] == item.key.id():
                    book_exsist = True
                    if parse_item['volume'] != item.volume:
                        parse_item['updated'] = item.volume+'->'+parse_item['volume']
                        Books.update_book(parse_item, author)
                    continue
            if not book_exsist:
                parse_item['updated'] = '0->'+parse_item['volume']
                Books.update_book(parse_item, author)

        self.response.headers['Content-Type'] = 'text/html'
        self.response.write(viewer.AuthorsWriter('books.html').write(parse_contents, author['name']))


class Books():
    def __init__(self):
        pass

    @staticmethod
    def update_book(item, author):
        book = models.AuthorsBooks(parent=models.authors_key(author['code']), id=item['id'])
        book.book = item['book']
        book.href = item['href']
        book.volume = item['volume']
        book.put()


class DatastoreFlush(webapp2.RequestHandler):
    def get(self):
        authors = models.Authors.query().fetch()

        for model in authors:
            contents = models.AuthorsBooks.query(ancestor=models.authors_key(model.key.id())).fetch()
            for item in contents:
                item.key.delete()


class UpdateBooks(webapp2.RequestHandler):
    def get(self):
        authors_model = models.Authors.query().fetch()

        authors = []
        for model in authors_model:
            author_url = '/author/'+model.key.id()
            taskqueue.add(url=author_url, method='GET')

            authors.append({
                'href': author_url,
                'href_samlib': core_parser.AuthorsParser.get_author_link(model.key.id()),
                'name': model.name
            })

        self.response.headers['Content-Type'] = 'text/html'
        self.response.write(viewer.AuthorsWriter('index.html').write(authors))


class LastUpdates(webapp2.RedirectHandler):
    def get(self):
        authors = models.Authors.query().fetch()
        date = datetime.today()
        date = datetime(date.year, date.month, date.day)
        total_contents = []
        for model in authors:
            author_code = model.key.id()

            author = {'code': author_code, 'name': model.name}

            contents = models.AuthorsBooks.query(ancestor=models.authors_key(author['code']))\
                .filter(models.AuthorsBooks._properties["updated_at"] > date)\
                .order(models.AuthorsBooks.updated_at)\
                .fetch()

            for item in contents:
                print item
                total_contents.append({
                    'author_id': author['code'],
                    'author': author['name'],
                    'book': item.book,
                    'id': item.key.id(),
                    'href': item.href,
                    'volume': item.volume
                })

        self.response.headers['Content-Type'] = 'text/html'
        self.response.write(viewer.AuthorsWriter('last-updates.html').write(total_contents, 'Updates!'))
