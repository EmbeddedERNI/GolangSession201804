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

	// Part1 OMIT
	// Get the webpage
	response, err := http.Get(url) // HL1
	if err != nil {
		errors <- err
		return
	}
	// Parse html code
	z := html.NewTokenizer(response.Body) // HL2

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken: // HL3
			// End of the document, we're done OMIT
			errors <- nil
			return
		case tt == html.StartTagToken: // HL4
			// find <a> tags
			t := z.Token()
			isAnchor := t.Data == "a"
			if isAnchor { //HL5
				//Extract href
				// END OMIT
				// log.Println("We found a <a>!") OMIT
				foundurl, finderr := findHref(&t)
				if finderr != nil {
					// log.Println("Not found href field") OMIT
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
				// Part2 OMIT
				urls <- CrawlerStep{url, urlformated.String(), depth}
				// END OMIT
			}
		}
	}

}
