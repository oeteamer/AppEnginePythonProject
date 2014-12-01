# coding: utf-8
import os

import jinja2

JINJA_ENVIRONMENT = jinja2.Environment(
    loader=jinja2.FileSystemLoader(os.path.dirname(__file__)+'/../templates'),
    extensions=['jinja2.ext.autoescape'],
    autoescape=True)


class AuthorsWriter():
    def __init__(self, template_file):
        self.template = JINJA_ENVIRONMENT.get_template(template_file)

    def write_index_date(self, author, contents):
        template_values = ({
            'author': author,
            'contents': contents
        })

        return self.template.render(template_values)