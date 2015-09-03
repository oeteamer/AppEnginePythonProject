package base

import (
	"appengine"
	"appengine/datastore"
	"net/http"
	"time"
)

func authorKey(r *http.Request, code string) *datastore.Key {
	return datastore.NewKey(appengine.NewContext(r), AuthorsKind, code, 0, nil)
}

func bookKey(r *http.Request, bookCode string, authorCode string) *datastore.Key {
	return datastore.NewKey(appengine.NewContext(r), BooksKind, bookCode, 0, authorKey(r, authorCode))
}

func saveAuthor(r *http.Request, Author *Author) error {
	Author.UpdatedAt = time.Now()
	if (Author.CreatedAt == time.Time{}) {
		Author.CreatedAt = time.Now()
	}

	_, err := datastore.Put(appengine.NewContext(r), Author.ID, Author)

	if err != nil {
		return err
	}

	Authors[Author.Code] = Author

	return err
}

func saveBook(r *http.Request, Book *Book) error {
	Book.UpdatedAt = time.Now()
	if (Book.CreatedAt == time.Time{}) {
		Book.CreatedAt = time.Now()
	}

	_, err := datastore.Put(appengine.NewContext(r), Book.ID, Book)

	return err
}

func saveBooks(r *http.Request, books []*Book) error {
	var ids []*datastore.Key
	for _, book := range books {
		book.UpdatedAt = time.Now()
		if (book.CreatedAt == time.Time{}) {
			book.CreatedAt = time.Now()
		}
		ids = append(ids, book.ID)
	}
	_, err := datastore.PutMulti(appengine.NewContext(r), ids, books)

	return err
}

func getAuthors(r *http.Request) {
	var authors []Author

	result := datastore.NewQuery(AuthorsKind)
	keys, _ := result.GetAll(appengine.NewContext(r), &authors)

	count := 0
	Authors = make(map[string]*Author, (len(keys) + 10))
	for _, key := range keys {
		authors[count].Code = key.StringID()
		authors[count].ID = authorKey(r, authors[count].Code)
		Authors[key.StringID()] = &authors[count]
		count++
	}
}

func getBooks(r *http.Request, authorCode string) {
	var books []Book

	result := datastore.NewQuery(BooksKind).Ancestor(authorKey(r, authorCode))
	keys, _ := result.GetAll(appengine.NewContext(r), &books)

	Authors[authorCode].Books = make(map[string]*Book)
	for a, key := range keys {
		books[a].Code = key.StringID()
		books[a].ID = bookKey(r, books[a].Code, authorCode)
		books[a].AuthorCode = authorCode
		Authors[authorCode].Books[key.StringID()] = &books[a]
	}
}

func getAuthor(r *http.Request, code string) {
	author := &Author{Code: code}
	k := authorKey(r, code)
	datastore.Get(appengine.NewContext(r), k, author)
}
