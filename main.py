# coding: utf-8

from bs4 import BeautifulSoup

import os

from google.appengine.api import modules
from google.appengine.api import urlfetch
import re

import jinja2
import webapp2

JINJA_ENVIRONMENT = jinja2.Environment(
    loader=jinja2.FileSystemLoader(os.path.dirname(__file__)),
    extensions=['jinja2.ext.autoescape'],
    autoescape=True)


def find_volumes(tag):
    return 'li' == tag.parent.name and tag.name == 'b'


def find_books(tag):
    return 'li' == tag.parent.name and tag.name == 'a'


class MainPage(webapp2.RequestHandler):
    def get(self):
        url = 'http://samlib.ru/a/aleksej_shpik/'

        (author, contents) = self.indexdate_parser(url)

        self.viewer({
            'author': author,
            'contents': contents
        })

    def indexdate_parser(self, url):
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

    def viewer(self, template_values):
        template = JINJA_ENVIRONMENT.get_template('index.html')
        self.response.headers['Content-Type'] = 'text/html'
        self.response.write(template.render(template_values))

app = webapp2.WSGIApplication([
    ('/', MainPage),
], debug=True)