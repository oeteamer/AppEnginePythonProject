package base

import (
	"appengine/datastore"
)

func authorKey(code string) *datastore.Key {
	return datastore.NewKey(Context, AuthorsKind, code, 0, nil)
}

//func query(c appengine.Context) *datastore.Query {
//	return datastore.NewQuery("Greeting").Ancestor(guestbookKey(c)).Order("-Date").Limit(10)
//}

func putObject(Author Author, code string) (*datastore.Key, error) {
	key := authorKey(code)

	return datastore.Put(Context, key, &Author)
}

var (
	Authors []Author
)

func getAuthors() {
	result := datastore.NewQuery(AuthorsKind)
	result.GetAll(Context, &Authors)
}

func getAuthorByCode(code string) (*Author, error) {
	k := authorKey(code)
	e := &Author{Code: code}
	err := datastore.Get(Context, k, e)

	return e, err
}
