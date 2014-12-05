# coding: utf-8
import sys

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
    def indexdate_parser(author_code):
        response = urlfetch.fetch(AuthorsParser.get_author_link(author_code)+'/indexdate.shtml').content
        response = response.decode('cp1251')
        soup = BeautifulSoup(response)
        soup = soup.dl

        volumes = soup.find_all(find_volumes)
        books = soup.find_all(find_books)

        contents = []
        for (i, link) in enumerate(books):
            contents.append({
                'book': link.get_text(),
                'id': link.get('href'),
                'href': AuthorsParser.get_author_link(author_code)+'/'+link.get('href'),
                'volume': volumes[i].get_text()
            })

        return contents

    @staticmethod
    def get_author_name_from_url(author_code):
        response = urlfetch.fetch(AuthorsParser.get_author_link(author_code)+'/indexdate.shtml')

        if response.status_code != 200:
            sys.exit('error with code: '+str(response.status_code))

        response = response.content
        response = response.decode('cp1251')
        soup = BeautifulSoup(response)
        author_string = soup.title.get_text()
        author_string = re.search('\.(.*?)\.', author_string.encode('utf8')).group(0, 1)
        name = author_string[1].decode('utf8')

        return name

    @staticmethod
    def get_author_link(author_code):
        return 'http://samlib.ru/'+author_code[0]+'/'+author_code