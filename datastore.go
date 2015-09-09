package base

import (
	"appengine"
	"appengine/datastore"
	//	"time"
)

func authorKey(c appengine.Context, code string) *datastore.Key {
	return datastore.NewKey(c, AuthorsKind, code, 0, nil)
}

func bookKey(c appengine.Context, bookCode string, authorCode string) *datastore.Key {
	return datastore.NewKey(c, BooksKind, bookCode, 0, authorKey(c, authorCode))
}

func saveAuthor(c appengine.Context, Author Author) error {
	_, err := datastore.Put(c, &Author.ID, &Author)

	return err
}

//func saveBook(c appengine.Context, Book Book) error {
//	Book.UpdatedAt = time.Now()
//	if (Book.CreatedAt == time.Time{}) {
//		Book.CreatedAt = time.Now()
//	}
//
//	_, err := datastore.Put(c, &Book.ID, Book)
//
//	return err
//}

func saveBooks(c appengine.Context, books []Book) error {
	var ids []*datastore.Key
	for a, _ := range books {
		ids = append(ids, &books[a].ID)
	}

	_, err := datastore.PutMulti(c, ids, books)

	return err
}

func getAuthors(c appengine.Context) {
	var authors []Author

	result := datastore.NewQuery(AuthorsKind)
	keys, _ := result.GetAll(c, &authors)

	count := 0
	Authors = make(map[string]Author, (len(keys) + 10))
	for _, key := range keys {
		authors[count].Code = key.StringID()
		authors[count].ID = *authorKey(c, authors[count].Code)
		Authors[key.StringID()] = authors[count]
		count++
	}
}

func getBooks(c appengine.Context, authorCode string) {
	var books []Book

	result := datastore.NewQuery(BooksKind).Ancestor(authorKey(c, authorCode))
	keys, _ := result.GetAll(c, &books)

	author := Authors[authorCode]
	author.Books = make(map[string]Book)
	for a, key := range keys {
		books[a].Code = key.StringID()
		books[a].ID = *bookKey(c, books[a].Code, authorCode)
		author.Books[key.StringID()] = books[a]
	}
	Authors[authorCode] = author
}

//func getAuthor(c appengine.Context, code string) {
//	author := &Author{Code: code}
//	k := authorKey(c, code)
//	datastore.Get(c, k, author)
//}
