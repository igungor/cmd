// fujinews is a stupid scraper.
//
// There is no channel to fetch the latest news from fuji x-series (new cameras
// or firmware updates). Even though a scraper is not a fantastic idea, it does the job.
//
// TODO(ig): Fetch news channel from RSS and push latest news to Twitter maybe?
package main

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
