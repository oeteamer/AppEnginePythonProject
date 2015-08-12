package base

import (
	"appengine/urlfetch"
	"bytes"
	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
	"fmt"
	"regexp"
)

var (
	AuthorMatch, _       = regexp.Compile("/author/([a-z_]+)")
	AuthorUpdateMatch, _ = regexp.Compile("/author/([a-z_]+)/update")
	authorNameRegexp, _  = regexp.Compile(`<title>.*?\.\s{0,}(.*?)\s{0,}\..*?</title>`)
)

func parseAuthorName(code string) (string, error) {
	authorLink := getAuthorLink(code)

	client := urlfetch.Client(Context)
	response, err := client.Get(authorLink)
	defer response.Body.Close()

	if err != nil {
		return "", err
	} else if response.StatusCode < 200 || response.StatusCode >= 300 {
		return "", fmt.Errorf(fmt.Sprint("StatusCode = ", response.StatusCode))
	}

	responseBodyReader, _ := charset.NewReader("windows-1251", response.Body)
	buf := new(bytes.Buffer)
	buf.ReadFrom(responseBodyReader)
	responseBody := buf.String()

	authorName := authorNameRegexp.FindStringSubmatch(responseBody)[1]

	return authorName, err
}

func getAuthorByUrl(url string) (string, error) {
	match := AuthorMatch.FindStringSubmatch(url)

	if len(match) == 2 {
		return match[1], nil
	} else {
		return "", fmt.Errorf("Author code is invalid")
	}
}

func getAuthorLink(code string) string {
	return fmt.Sprint("http://samlib.ru/", string(code[0]), "/", code)
}

func isAuthorUpdate(url string) bool {
	match := AuthorUpdateMatch.FindStringSubmatch(url)
	if len(match) == 2 {
		return true
	} else {
		return false
	}
}
