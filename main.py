import webapp2

app = webapp2.WSGIApplication([
    ('/', 'app.views.HelloWorld'),
], debug=True)
