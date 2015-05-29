from google.appengine.ext import ndb


def authors_key(author):
    return ndb.Key('Authors', author)


class Authors(ndb.Model):
    name = ndb.StringProperty(indexed=True)
    created_at = ndb.DateTimeProperty(auto_now_add=True, indexed=True)
    updated_at = ndb.DateTimeProperty(auto_now=True, indexed=True)

    @staticmethod
    def create(author):
        author_entity = Authors(id=author['code'])
        author_entity.name = author['name']
        return author_entity


class AuthorsBooks(ndb.Model):
    book = ndb.StringProperty(indexed=False)
    href = ndb.StringProperty(indexed=False)
    volume = ndb.StringProperty(indexed=True)
    update_info = ndb.StringProperty(indexed=False)
    created_at = ndb.DateTimeProperty(auto_now_add=True, indexed=True)
    updated_at = ndb.DateTimeProperty(auto_now=True, indexed=True)

    @staticmethod
    def create_book_entity(item, author):
        book = AuthorsBooks(parent=authors_key(author['code']), id=item['id'])
        book.book = item['book']
        book.href = item['href']
        book.volume = item['volume']
        book.update_info = item['update_info']

        return book

    @staticmethod
    def update_book(item, author):
        book = AuthorsBooks.create_book_entity(item, author)
        book.put()

    @staticmethod
    def update_books(items):
        ndb.put_multi(items)
