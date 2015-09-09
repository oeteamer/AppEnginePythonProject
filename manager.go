package base

import (
	"appengine"
	"net/http"
	"time"
)

var (
	started     = false
	Authors     map[string]Author
	AuthorsKind = "Authors"
	BooksKind   = "AuthorsBooks"
	AllAuthors  = []string{
		"aleksandrin_a_a",
		"arestowich_a_w",
		"atamanow_m_a",
		"bashun_w_m",
		"carenko_t_p",
		"chitajushij",
		"christyan_b",
		"demchenko_aw",
		"dmitrij_c",
		"drawin_i",
		"ekzar",
		"feliks_o",
		"fox",
		//		"glushanowskij_a_a",
		//		"hrabryh_k",
		//		"hwan_a",
		//		"karpow_a_n",
		//		"kim_sergej_aleksandrowich",
		//		"konard",
		//		"kowtun_d",
		//		"liwidus_a_a",
		//		"lomanowa_e_a",
		//		"magazinnikow_i_w",
		//		"metelxskij_n_a",
		//		"mirotworcew_p",
		//		"muhin_d_w",
		//		"nikitin_m_a",
		//		"odinow_d_j",
		//		"orlow_i",
		//		"plotnikow_sergej_aleksandrowich",
		//		"raavasta",
		//		"rajkar",
		//		"raznicin_w_s",
		//		"sadow_m_w",
		//		"samuilow_a_w",
		//		"sergej_strazhnyj",
		//		"shtajn_f",
		//		"sirius_m",
		//		"skyd",
		//		"starh_a",
		//		"starolisow_a_e",
		//		"sworm",
		//		"tagern",
		//		"tign",
		//		"tkach",
		//		"wechnyj_a_a",
		//		"wiktor_dobryj",
		//		"zajcew_aleskandr",
		//		"zajcew_p_a",
		//		"zloj",
	}
)

func toInit(r *http.Request) appengine.Context {
	if started == false {
		c := appengine.NewContext(r)
		c.Infof("%s", "up")
		module := appengine.ModuleName(c)
		if module != "default" {
			AuthorsKind = AuthorsKind + "_" + module
			BooksKind = BooksKind + "_" + module
		}
		getAuthors(c)
		for code, _ := range Authors {
			getBooks(c, code)
		}
		started = true
	}

	return appengine.NewContext(r)
}

func createNewAuthor(c appengine.Context, code string) (Author, error) {
	newAuthor := Author{Code: code}
	name, books, err := parseAuthorPage(c, code)
	if err != nil {
		return newAuthor, err
	} else {
		newAuthor.Name = name
		newAuthor.ID = *authorKey(c, code)
		newAuthor.UpdatedAt = time.Now()
		newAuthor.CreatedAt = time.Now()

		err = saveAuthor(c, newAuthor)
		if err != nil {
			return newAuthor, err
		}

		var (
			newBooks []Book
			count    = 0
		)
		newAuthor.Books = make(map[string]Book)
		for _, book := range books {
			count++
			book.ID = *bookKey(c, book.Code, newAuthor.Code)
			book.UpdateInfo = "Added " + formatTime(time.Now())
			newAuthor.Books[book.Code] = book
			newBooks = append(newBooks, book)
			if count == 499 {
				err = saveBooks(c, newBooks)
				newBooks = newBooks[:0]
			}
		}

		Authors[newAuthor.Code] = newAuthor

		err = saveBooks(c, newBooks)
	}

	return newAuthor, err
}

func updateAuthorBooks(c appengine.Context, Author Author) error {
	var (
		err          error
		books        []Book
		updatedBooks []Book
		count        = 0
	)
	_, books, err = parseAuthorPage(c, Author.Code)

	if err != nil {
		return err
	}
	for _, book := range books {
		if _, b := Author.Books[book.Code]; b {
			if book.Volume != Author.Books[book.Code].Volume {
				count++
				authorBook := Author.Books[book.Code]

				authorBook.Name = book.Name
				authorBook.UpdateInfo = "Updated " + Author.Books[book.Code].Volume + "->" + book.Volume
				authorBook.Volume = book.Volume
				authorBook.UpdatedAt = book.UpdatedAt

				Author.Books[book.Code] = authorBook
				updatedBooks = append(updatedBooks, Author.Books[book.Code])
				c.Infof("Update book %s", Author.Books[book.Code].Name)
			}
		} else {
			count++
			book.ID = *bookKey(c, book.Code, Author.Code)
			book.UpdateInfo = "Added " + formatTime(time.Now())
			Author.Books[book.Code] = book
			updatedBooks = append(updatedBooks, Author.Books[book.Code])
			c.Infof("Add book %s", Author.Books[book.Code].Name)
		}

		if count == 499 {
			err = saveBooks(c, updatedBooks)
			updatedBooks = updatedBooks[:0]
		}
	}

	if len(updatedBooks) > 0 {
		err = saveBooks(c, updatedBooks)
	} else {
		c.Infof("Not new books updates %s", "")
	}

	return err
}

func handleError(w http.ResponseWriter, err error, c appengine.Context) {
	c.Errorf("%v", err)
	http.Error(w, err.Error(), 500)
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
