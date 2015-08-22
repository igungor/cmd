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

	"github.com/PuerkitoBio/goquery"
)

var baseURL = "http://fujilove.com/category/news/"

func main() {
	doc, err := goquery.NewDocument(baseURL)
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
