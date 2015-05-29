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
            models.Authors.create({
                'code': author_code,
                'name': core_parser.AuthorsParser.get_author_name_from_url(author_code)}
            ).put()
            author_model = models.Authors.get_by_id(author_code)

        author = {'code': author_model.key.id(), 'name': author_model.name}

        contents = models.AuthorsBooks.query(ancestor=models.authors_key(author['code']))\
            .order(models.AuthorsBooks.updated_at)\
            .fetch()

        parse_contents = core_parser.AuthorsParser.indexdate_parser(author['code'])

        updated_items = []
        for parse_item in parse_contents:
            book_exsist = False
            for item in contents:
                if parse_item['id'] == item.key.id():
                    parse_item['updated_at'] = datetime_to_string(item.updated_at)
                    book_exsist = True
                    if parse_item['volume'] != item.volume:
                        parse_item['update_info'] = 'Update '+item.volume+'->'+parse_item['volume']
                        book = models.AuthorsBooks.create_book_entity(parse_item, author)
                        updated_items.append(book)
                    continue
            if not book_exsist:
                parse_item['updated_at'] = datetime_to_string(datetime.today())
                parse_item['update_info'] = 'Added '+datetime_to_string(datetime.today())
                book = models.AuthorsBooks.create_book_entity(parse_item, author)
                updated_items.append(book)

        if updated_items:
            models.AuthorsBooks.update_books(updated_items)

        self.response.headers['Content-Type'] = 'text/html'
        self.response.write(viewer.AuthorsWriter('books.html').write(parse_contents, author['name']))


class LastUpdates(webapp2.RedirectHandler):
    def get(self):
        authors = models.Authors.query().fetch()
        date = datetime
        date_from = datetime(date.today().year, date.today().month, date.today().day)
        date_to = datetime(
            date.today().year,
            date.today().month,
            date.today().day,
            date.today().hour,
            date.today().minute,
            date.today().second
        )
        total_contents = []
        author = {}
        for model in authors:
            author_code = model.key.id()
            author[author_code] = model.name

        contents = models.AuthorsBooks.query()\
            .filter(models.AuthorsBooks._properties["updated_at"] >= date_from)\
            .filter(models.AuthorsBooks._properties["updated_at"] <= date_to)\
            .order(-models.AuthorsBooks.updated_at)\
            .fetch()\

        for item in contents:
            if item.key.parent().id() in author:
                total_contents.append({
                    'author_id': item.key.parent().id(),
                    'author': author[item.key.parent().id()],
                    'book': item.book,
                    'id': item.key.id(),
                    'href': item.href,
                    'volume': item.volume,
                    'update_info': item.update_info,
                    'updated_at': datetime_to_string(item.updated_at)
                })
            else:
                item.key.delete()

        self.response.headers['Content-Type'] = 'text/html'
        self.response.write(viewer.AuthorsWriter('last-updates.html').write(total_contents, 'Updates!'))


class UpdateBooks(webapp2.RequestHandler):
    def get(self):
        authors_model = models.Authors.query().fetch()

        for model in authors_model:
            author_url = '/author/'+model.key.id()
            taskqueue.add(url=author_url, method='GET')

        self.redirect('/last-updates', True)


class DatastoreFlush(webapp2.RequestHandler):
    @staticmethod
    def get():
        authors = models.Authors.query().fetch()

        for model in authors:
            contents = models.AuthorsBooks.query(ancestor=models.authors_key(model.key.id())).fetch()
            for item in contents:
                item.key.delete()


class TaskQueueStats(webapp2.RequestHandler):
    def get(self):
        tasks = taskqueue.QueueStatistics.fetch([taskqueue.Queue('default')])

        self.response.headers['Content-Type'] = 'text/html'
        self.response.write(tasks[0].tasks)