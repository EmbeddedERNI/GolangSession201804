package main

import (
	"errors"
	"golang.org/x/net/html"
	// "log"
	"net/http"
	urltools "net/url"
)

type CrawlerStep struct {
	baseUrl  string
	childUrl string
	depth    int
}

func findHref(t *html.Token) (string, error) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			return a.Val, nil
		}
	}
	return "", errors.New("Href not found")
}

func FindUrlsIn(url string, depth int, urls chan<- CrawlerStep, errors chan<- error) {

	response, err := http.Get(url)
	if err != nil {
		errors <- err
		return
	}

	z := html.NewTokenizer(response.Body)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			errors <- nil
			return
		case tt == html.StartTagToken:
			t := z.Token()
			isAnchor := t.Data == "a"
			if isAnchor {
				// log.Println("We found a <a>!")
				foundurl, finderr := findHref(&t)
				if finderr != nil {
					// log.Println("Not found href field")
					continue
				}
				urlformated, formaterr := urltools.ParseRequestURI(foundurl)
				if formaterr != nil {
					// log.Println(formaterr)
					continue
				}
				if urlformated.Hostname() == "" {
					// log.Println("Not valid host")
					continue
				}
				urls <- CrawlerStep{url, urlformated.String(), depth}
			}
		}
	}

}
