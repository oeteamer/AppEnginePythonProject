# coding: utf-8

from bs4 import BeautifulSoup

from google.appengine.api import urlfetch
import re


def find_volumes(tag):
    return 'li' == tag.parent.name and tag.name == 'b'


def find_books(tag):
    return 'li' == tag.parent.name and tag.name == 'a'


class AuthorsParser():
    def __init__(self):
        pass

    @staticmethod
    def indexdate_parser(url):
        response = urlfetch.fetch(url+'indexdate.shtml').content
        response = response.decode('cp1251')
        soup = BeautifulSoup(response)

        author = soup.title
        author = re.search('\.(.*?)\.', author.encode('utf8')).group(0, 1)
        author = author[1].decode('utf8')

        contents_soup = soup.dl

        volumes = contents_soup.find_all(find_volumes)
        books = contents_soup.find_all(find_books)

        contents = []
        for (i, link) in enumerate(books):
            contents.append({
                'book': link.get_text(),
                'href': url+link.get('href'),
                'volume': volumes[i].get_text()
            })

        return [author, contents]