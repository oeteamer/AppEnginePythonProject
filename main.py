import webapp2

from authors import handlers


app = webapp2.WSGIApplication([
    webapp2.Route('/author/<author_code:\w+>', handlers.Authors),
    webapp2.Route('/', handlers.Index),
    webapp2.Route('/datastore_flush', handlers.DatastoreFlush),
    webapp2.Route('/update-all', handlers.UpdateBooks),
    webapp2.Route('/last-updates', handlers.LastUpdates),
], debug=True)
