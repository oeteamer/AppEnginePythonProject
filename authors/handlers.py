import webapp2
import models

from core import core_parser
from core import viewer
from datetime import datetime
from google.appengine.api import taskqueue


def datetime_to_string(date_time):
    year = str(date_time.year)
    month = date_string_plus_null(date_time.month)
    day = date_string_plus_null(date_time.day)
    hour = date_string_plus_null(date_time.hour)
    minute = date_string_plus_null(date_time.minute)
    second = date_string_plus_null(date_time.second)

    return year+'-'+month+'-'+day+' '+hour+':'+minute+':'+second


def date_string_plus_null(value):
    return '0'+str(value) if value < 10 else str(value)


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
                    parse_item['updated_at'] = datetime_to_string(item.updated_at)
                    book_exsist = True
                    if parse_item['volume'] != item.volume:
                        parse_item['updated'] = 'Update '+item.volume+'->'+parse_item['volume']
                        Books.update_book(parse_item, author)
                    continue
            if not book_exsist:
                parse_item['updated_at'] = datetime_to_string(datetime.today())
                parse_item['updated'] = 'Added '+datetime_to_string(datetime.today())
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
        book.update_info = item['updated']
        book.put()


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
                    'volume': item.volume,
                    'updated_at': datetime_to_string(item.updated_at)
                })

        self.response.headers['Content-Type'] = 'text/html'
        self.response.write(viewer.AuthorsWriter('last-updates.html').write(total_contents, 'Updates!'))


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


class DatastoreFlush(webapp2.RequestHandler):
    def get(self):
        authors = models.Authors.query().fetch()

        for model in authors:
            contents = models.AuthorsBooks.query(ancestor=models.authors_key(model.key.id())).fetch()
            for item in contents:
                item.key.delete()

