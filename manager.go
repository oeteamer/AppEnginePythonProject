package base

import (
	"appengine"
	"log"
	"net/http"
	"time"
)

var (
	started     = false
	Authors     map[string]*Author
	AuthorsKind = "Authors"
	BooksKind   = "AuthorsBooks"
)

func toInit(r *http.Request) {
	if started == false {
		log.Print("up")
		//		c := appengine.NewContext(r)
		//		module := appengine.ModuleName(c)
		//		if module != "default" {
		//			AuthorsKind = AuthorsKind + "_" + module
		//			BooksKind = BooksKind + "_" + module
		//		}
		getAuthors(r)
		for code, _ := range Authors {
			getBooks(r, code)
		}
		started = true
	}
}

func createNewAuthor(r *http.Request, code string) (*Author, error) {
	newAuthor := &Author{Code: code}
	name, books, err := parseAuthorPage(r, code)
	if err != nil {
		return newAuthor, err
	} else {
		newAuthor.Name = name
		newAuthor.ID = authorKey(r, code)
		Authors[newAuthor.Code] = newAuthor
		err = saveAuthor(r, newAuthor)
		if err != nil {
			return newAuthor, err
		}

		Authors[newAuthor.Code].Books = make(map[string]*Book, (len(books) + 10))
		var (
			newBooks []*Book
			count    = 0
		)
		for _, book := range books {
			count++
			book.ID = bookKey(r, book.Code, newAuthor.Code)
			book.UpdateInfo = "Added " + formatTime(time.Now())
			Authors[newAuthor.Code].Books[book.Code] = book
			newBooks = append(newBooks, book)
			if count == 499 {
				err = saveBooks(r, newBooks)
				newBooks = newBooks[:0]
			}
		}

		err = saveBooks(r, newBooks)
	}

	return newAuthor, err
}

func updateAuthorBooks(r *http.Request, Author *Author) error {
	var (
		err          error
		books        []*Book
		updatedBooks []*Book
		count        = 0
	)
	_, books, err = parseAuthorPage(r, Author.Code)

	for _, book := range books {
		if _, b := Author.Books[book.Code]; b {
			if book.Volume != Author.Books[book.Code].Volume {
				count++
				Author.Books[book.Code].UpdateInfo = "Updated " + Author.Books[book.Code].Volume + "->" + book.Volume
				Author.Books[book.Code].Volume = book.Volume
				updatedBooks = append(updatedBooks, Author.Books[book.Code])
				c := appengine.NewContext(r)
				c.Infof("Update book %s", Author.Books[book.Code].Name)
			}
		} else {
			count++
			book.ID = bookKey(r, book.Code, Author.Code)
			book.UpdateInfo = "Added " + formatTime(time.Now())
			Author.Books[book.Code] = book
			updatedBooks = append(updatedBooks, Author.Books[book.Code])
			c := appengine.NewContext(r)
			c.Infof("Add book %s", Author.Books[book.Code].Name)
		}

		if count == 499 {
			err = saveBooks(r, updatedBooks)
			updatedBooks = updatedBooks[:0]
		} else {
			c := appengine.NewContext(r)
			c.Infof("Not new books updates", "")
		}
	}

	if len(updatedBooks) > 0 {
		err = saveBooks(r, updatedBooks)
	}

	return err
}

func handleError(w http.ResponseWriter, err error, r *http.Request) {
	log.Print(err.Error())
	c := appengine.NewContext(r)
	c.Errorf("%v", err)
	http.Error(w, err.Error(), 500)
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}