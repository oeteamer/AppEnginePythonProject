package base

import (
	"appengine"
	"appengine/urlfetch"
	"bytes"
	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var (
	AuthorMatch, _       = regexp.Compile("/author/([a-z_]+)")
	AuthorUpdateMatch, _ = regexp.Compile("/author/([a-z_]+)/update")
	authorNameRegexp, _  = regexp.Compile(`<title>.*?\.\s{0,}(.*?)\s{0,}\..*?</title>`)
	booksRegexp, _       = regexp.Compile(`<DL><DT>(.*?)</DL>`)
	bookRegexp, _        = regexp.Compile(`<A HREF=(.*)><b>(.*)</b></A>`)
	bookVolumeRegexp, _  = regexp.Compile(`<b>([0-9]+k)</b>`)
)

func parseAuthorPage(r *http.Request, code string) (string, []Book, error) {
	var (
		authorName string
		books      []Book
		authorLink = getAuthorLink(code)
	)
	//	client := urlfetch.Client(appengine.NewContext(r))
	//	response, err := client.Get(fmt.Sprint(authorLink, "/indexdate.shtml"))

	tr := &urlfetch.Transport{Context: appengine.NewContext(r), Deadline: time.Duration(30) * time.Second}

	req, err := http.NewRequest("GET", fmt.Sprint(authorLink, "/indexdate.shtml"), strings.NewReader(""))

	response, err := tr.RoundTrip(req)
	if err != nil {
		return authorName, books, err
	} else if response.StatusCode < 200 || response.StatusCode >= 300 {
		return authorName, books, fmt.Errorf(fmt.Sprint("StatusCode = ", response.StatusCode))
	}

	defer response.Body.Close()

	body := toUTF(response.Body)

	authorName, err = checkMatch(authorNameRegexp.FindStringSubmatch(body))

	booksMatches := booksRegexp.FindAllStringSubmatch(body, 1000)
	for _, val := range booksMatches {
		volume := bookVolumeRegexp.FindStringSubmatch(val[0])
		if len(volume) > 1 {
			book := bookRegexp.FindStringSubmatch(val[0])
			books = append(books, Book{
				Code:      book[1],
				Name:      book[2],
				Href:      authorLink + "/" + book[1],
				Volume:    volume[1],
				UpdatedAt: time.Now(),
				CreatedAt: time.Now(),
			})
		}
	}

	return authorName, books, err
}

func getAuthorByUrl(url string) (string, error) {
	return checkMatch(AuthorMatch.FindStringSubmatch(url))
}

func getAuthorLink(code string) string {
	return fmt.Sprint("http://samlib.ru/", string(code[0]), "/", code)
}

func isAuthorUpdate(url string) bool {
	_, err := checkMatch(AuthorUpdateMatch.FindStringSubmatch(url))

	if err != nil {
		return false
	} else {
		return true
	}
}

func toUTF(win1251 io.Reader) string {
	reader, _ := charset.NewReader("windows-1251", win1251)
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	return buf.String()
}

func checkMatch(match []string) (string, error) {
	if len(match) == 2 {
		return match[1], nil
	} else {
		return "", fmt.Errorf("Matches not found")
	}
}
