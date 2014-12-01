import webapp2

from authors import handlers


app = webapp2.WSGIApplication([
    ('/', handlers.Authors),
], debug=True)