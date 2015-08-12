package base

import (
	"appengine"
	"net/http"
	"time"
)

var (
	Context     appengine.Context
	AuthorsKind = "Authors"
)

func createNewAuthor(newAuthor *Author, code string) (*Author, error) {
	name, err := parseAuthorName(code)
	if err != nil {
		return newAuthor, err
	} else {
		newAuthor.Name = name
		newAuthor.CreatedAt = time.Now()
		newAuthor.UpdatedAt = time.Now()

		_, err := putObject(Author{
			Name:      newAuthor.Name,
			CreatedAt: newAuthor.CreatedAt,
			UpdatedAt: newAuthor.UpdatedAt,
		}, code)

		if err != nil {
			return newAuthor, err
		}
	}

	return newAuthor, err
}

func setContext(r *http.Request) {
	if Context == nil {
		Context = appengine.NewContext(r)
		module := appengine.ModuleName(Context)
		if module != "default" {
			AuthorsKind = AuthorsKind + "_" + module
		}
	}
}
