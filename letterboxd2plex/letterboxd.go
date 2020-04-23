package main

import (
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36"
	baseURL   = "https://letterboxd.com"
)

type List struct {
	Title   string
	Summary string
	Movies  []string
}

func FetchMovieList(url string) (List, error) {
	return fetchMovieList(url, nil)
}

func fetchMovieList(url string, list *List) (List, error) {
	client := &http.Client{Timeout: time.Minute}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return List{}, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return List{}, err
	}

	if list == nil {
		list = &List{}
	}

	doc.Find("div .film-poster img").Each(func(i int, sel *goquery.Selection) {
		title, ok := sel.Attr("alt")
		if !ok {
			return
		}

		list.Movies = append(list.Movies, title)
	})

	if href := hasNext(doc); href != "" {
		url = baseURL + href
		return fetchMovieList(url, list)
	}

	list.Title = title(doc)
	list.Summary = summary(doc)

	return *list, nil
}

func hasNext(doc *goquery.Document) string {
	var href string
	doc.Find("a.next").Each(func(i int, sel *goquery.Selection) {
		href, _ = sel.Attr("href")
	})
	return href
}

func title(doc *goquery.Document) string {
	var title string
	doc.Find("div.list-title-intro > h1").Each(func(i int, sel *goquery.Selection) {
		title = sel.Text()
	})
	return title
}

func summary(doc *goquery.Document) string {
	var summary string
	doc.Find("div.list-title-intro > .body-text").Each(func(i int, sel *goquery.Selection) {
		summary = sel.Text()
	})
	return summary
}
