import webapp2
import models

from core import core_parser
from core import viewer


class Index(webapp2.RequestHandler):
    def get(self):
        authors_model = models.Authors.query()
        authors_model = authors_model.fetch()
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
        author = {'code': author_code, 'name': core_parser.AuthorsParser.get_author_name_from_url(author_code)}

        author_model = models.Authors.get_by_id(author['code'])
        if not author_model:
            models.Authors.create(author).put()

        contents = models.AuthorsBooks.query(ancestor=models.authors_key(author['code'])).order(models.AuthorsBooks.updated_at)
        contents = contents.fetch()
        if not contents:
            contents = core_parser.AuthorsParser.indexdate_parser(author['code'])
            for item in contents:
                book = models.AuthorsBooks(parent=models.authors_key(author['code']), id=item['id'])
                book.book = item['book']
                book.href = item['href']
                book.volume = item['volume']
                book.put()

        # save author
        # models.Authors.create(author).put()

        # save books from indexdate
        # contents = core_parser.AuthorsParser.indexdate_parser(author['code'])
        # for item in contents:
        #     author_entity = models.AuthorsBooks(parent=models.authors_key(author['code']), id=item['href'])
        #     author_entity.book = item['book']
        #     author_entity.volume = item['volume']
        #     author_entity.put()

        self.response.headers['Content-Type'] = 'text/html'
        self.response.write(viewer.AuthorsWriter('books.html').write(contents, author['name']))

    @staticmethod
    def update_books():
        authors = models.Authors.query()
        authors = authors.fetch()

        for model in authors:
            contents = models.AuthorsBooks.query(ancestor=models.authors_key(model.key.id()))


class DatastoreFlush(webapp2.RequestHandler):
    def get(self):
        authors = models.Authors.query()
        authors = authors.fetch()

        for model in authors:
            contents = models.AuthorsBooks.query(ancestor=models.authors_key(model.key.id()))
            contents = contents.fetch()
            for item in contents:
                item.key.delete()
