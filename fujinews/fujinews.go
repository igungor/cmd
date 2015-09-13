// fujinews is a stupid scraper.
//
// There isn't any **direct** and **no BS** channel to recieve latest firmware
// update news for Fujifilm X series (and possibly new cameras and lenses) as
// far as i know.  Even though scraping a site to fill this gap is not a
// brilliant idea, it does the job for me.
package main

// TODO(ig): Be a good citizen and extract `news` category from http://feeds.feedburner.com/fujilove

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	rss "github.com/jteeuwen/go-pkg-rss"
)

const (
	cacheTimeout = 24 * 60
	webURL       = "http://fujilove.com/category/news/"
	rssURL       = "http://feeds.feedburner.com/fujilove"
)

var c = http.Client{Timeout: 1 * time.Second}

func scrape() {
	resp, err := c.Get(webURL)
	if err != nil {
		log.Fatal(resp)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("#sidebar a").Each(func(i int, s *goquery.Selection) {
		if val, ok := s.Attr("target"); ok {
			if val == "_self" {
				href, _ := s.Attr("href")
				fmt.Printf("- %v\n\t%v\n", s.Text(), href)
			}
		}
	})
}

func itemHandler(f *rss.Feed, ch *rss.Channel, items []*rss.Item) {
	for _, item := range items {
		for _, c := range item.Categories {
			// Firmware and new lens/camera updates are in `News` category.
			if c.Text != "News" {
				continue
			}

			fmt.Println(item.Title)
		}
	}
}

func main() {
	scrape()
}
