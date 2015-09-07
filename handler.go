package base

import (
	"appengine"
	"appengine/taskqueue"
	"fmt"
	"net/http"
	"time"
)

func init() {
	http.HandleFunc("/", showUpdatedBooks)
	http.HandleFunc("/author/", getAuthorPage)
	http.HandleFunc("/update-all", updateBooks)
	http.HandleFunc("/authors", showAuthors)
	http.HandleFunc("/task-stats", showTaskStats)
	http.HandleFunc("/add-old-authors", addOldAuthors)

	http.HandleFunc("/_ah/start", start)
	http.HandleFunc("/_ah/stop", start)
}

func showAuthors(w http.ResponseWriter, r *http.Request) {
	toInit(r)

	AuthorsTemplate.Execute(w, Authors)
}

func getAuthorPage(w http.ResponseWriter, r *http.Request) {
	toInit(r)

	var (
		AuthorEntity Author
		code         string
		err          error
	)

	code, err = getAuthorByUrl(r.URL.Path)
	if err != nil {
		handleError(w, err, r)
		return
	}

	if _, b := Authors[code]; b {
		AuthorEntity = Authors[code]

		if isAuthorUpdate(r.URL.Path) {
			err = updateAuthorBooks(r, AuthorEntity)
			if err != nil {
				handleError(w, err, r)
				return
			}
			fmt.Fprint(w, "done")

			return
		}
	} else {
		AuthorEntity, err = createNewAuthor(r, code)
		if err != nil {
			handleError(w, err, r)
			return
		}
	}

	AuthorTemplate.Execute(w, AuthorEntity)
}

func updateBooks(w http.ResponseWriter, r *http.Request) {
	toInit(r)

	for code := range Authors {
		//		channelParse <- Channel{c: appengine.NewContext(r), authorCode: code}
		t := taskqueue.NewPOSTTask(fmt.Sprint("/author/", code, "/update"), map[string][]string{})
		if _, err := taskqueue.Add(appengine.NewContext(r), t, "default"); err != nil {
			handleError(w, err, r)
			return
		}
	}

	fmt.Fprint(w, "done")
}

func showUpdatedBooks(w http.ResponseWriter, r *http.Request) {
	toInit(r)

	var updatedAuthors = make(map[string]Author)
	for code, author := range Authors {
		updatedAuthor := author
		updatedBooks := make(map[string]Book)
		for key, book := range author.Books {
			if time.Now().Sub(book.UpdatedAt) <= (time.Hour * 24) {
				updatedBooks[key] = book
			}
		}
		updatedAuthor.Books = updatedBooks
		updatedAuthors[code] = updatedAuthor
	}

	UpdatedBooksTemplate.Execute(w, updatedAuthors)
}

func start(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "")
}

func showTaskStats(w http.ResponseWriter, r *http.Request) {
	queues, _ := taskqueue.QueueStats(appengine.NewContext(r), []string{"default"}, 0)

	taskAmount := 0
	for _, queue := range queues {
		taskAmount += queue.Tasks
	}

	fmt.Fprint(w, taskAmount)
}

func addOldAuthors(w http.ResponseWriter, r *http.Request) {
	toInit(r)

	for _, code := range AllAuthors {
		//		channelParse <- Channel{c: appengine.NewContext(r), authorCode: code}
		t := taskqueue.NewPOSTTask(fmt.Sprint("/author/", code), map[string][]string{})
		if _, err := taskqueue.Add(appengine.NewContext(r), t, "default"); err != nil {
			handleError(w, err, r)
			return
		}
	}

	fmt.Fprint(w, "done")
}
