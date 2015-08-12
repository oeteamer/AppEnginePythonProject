package base

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", showUpdatedBooks)
	http.HandleFunc("/author/", author)
	http.HandleFunc("/update-all", updateBooks)
	http.HandleFunc("/authors", showAuthors)
}

func showAuthors(w http.ResponseWriter, r *http.Request) {
	setContext(r)

	getAuthors()

	for a := 0; a < len(Authors); a++ {
		fmt.Fprint(w, Authors[a].Name)
		fmt.Fprint(w, "\n")
	}
}

func author(w http.ResponseWriter, r *http.Request) {
	var (
		AuthorEntity *Author
		code         string
		err          error
	)

	setContext(r)
	code, err = getAuthorByUrl(r.URL.Path)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	AuthorEntity, err = getAuthorByCode(code)
	if err != nil && err.Error() == "datastore: no such entity" {
		AuthorEntity, err = createNewAuthor(AuthorEntity, AuthorEntity.Code)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if isAuthorUpdate(r.URL.Path) {
		err = updateAuthorBooks(AuthorEntity)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprint(w, "done")

		return
	}

	fmt.Fprint(w, AuthorEntity.Name)
}

func updateBooks(w http.ResponseWriter, r *http.Request) {
	setContext(r)

	fmt.Fprint(w, r.Body)
}

func updateAuthorBooks(Author *Author) error {
	var err error

	return err
}

func showUpdatedBooks(w http.ResponseWriter, r *http.Request) {
	setContext(r)

	fmt.Fprint(w, "showUpdatedBooks")
}
