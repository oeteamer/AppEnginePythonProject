# coding: utf-8
import sys

from bs4 import BeautifulSoup

from google.appengine.api import urlfetch
import re


def find_volumes(tag):
    volume = False
    if 'li' == tag.parent.name and tag.name == 'b':
        volume = True
        if not re.search('(.*k)$', tag.get_text().encode('utf8')):
            volume = False

    return volume


def find_books(tag):
    return 'li' == tag.parent.name and tag.name == 'a'


def find_book(tag):
    return tag.name == 'a' and tag.get('href')


class AuthorsParser():
    def __init__(self):
        pass

    @staticmethod
    def indexdate_parser(author_code):
        urlfetch.set_default_fetch_deadline(60)
        response = urlfetch.fetch(AuthorsParser.get_author_link(author_code)+'/indexdate.shtml').content
        response = response.decode('cp1251')
        dls = re.findall('<DL>(.*)</DL>', response)

        contents = []
        for (i, dl) in enumerate(dls):
            book = re.search('<A HREF=(.*)><b>(.*)</b></A>', dl)
            if book:
                contents.append({
                    'book': book.group(2),
                    'id': book.group(1),
                    'href': AuthorsParser.get_author_link(author_code)+'/'+book.group(1),
                    'volume': re.search('<b>([0-9]+k)</b>', dl).group(1)
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