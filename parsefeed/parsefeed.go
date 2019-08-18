package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"unicode"

	"github.com/godexsoft/gofeed"
	"github.com/godexsoft/gofeed/rss"
)

const (
	defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36"
)

func main() {
	var (
		flagURL  = flag.String("u", "", "Download URL and parse")
		flagSort = flag.Bool("s", false, "Sort feed items by pubdate")
	)
	flag.Parse()

	var r io.Reader = os.Stdin

	if *flagURL != "" {
		fmt.Println(flag.Arg(0))
		req, _ := http.NewRequest("GET", *flagURL, nil)
		req.Header.Set("User-Agent", defaultUserAgent)
		req.Header.Set("Content-Type", "application/atom+xml")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		r = resp.Body
	}

	printOnly := func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}

	body, _ := ioutil.ReadAll(r)
	body = bytes.Map(printOnly, body)

	p := gofeed.NewParser()
	p.RSSTranslator = newRssTranslator()
	rf, err := p.Parse(bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	if *flagSort {
		sort.Slice(rf.Items, func(i, j int) bool {
			pubdate0 := rf.Items[i].PublishedParsed
			pubdate1 := rf.Items[j].PublishedParsed
			if pubdate0 == nil {
				return false
			}
			if pubdate1 == nil {
				return false
			}
			p0 := *pubdate0
			p1 := *pubdate1
			return p0.After(p1)
		})
	}

	for _, item := range rf.Items {
		fmt.Printf("pubdate: %v - %v\n", item.PublishedParsed, item.Title)
	}
}

type rssTranslator struct {
	d *gofeed.DefaultRSSTranslator
}

func newRssTranslator() *rssTranslator {
	return &rssTranslator{
		d: &gofeed.DefaultRSSTranslator{},
	}
}

func (rt *rssTranslator) Translate(src interface{}) (*gofeed.Feed, error) {
	rf, ok := src.(*rss.Feed)
	if !ok {
		return nil, fmt.Errorf("feed is not RSS")
	}

	feed, err := rt.d.Translate(rf)
	if err != nil {
		return nil, err
	}

	for i, item := range rf.Items {
		for _, fn := range []func(*rss.Item) string{
			findFromEnclosure,
			findFromMediaGroup,
			findFromMediaContent,
			findFromMediaThumbnail,
		} {
			thumb := fn(item)
			if thumb != "" {
				feed.Items[i].Image = &gofeed.Image{URL: thumb}
				break
			}
		}
	}

	return feed, nil
}

func findFromMediaThumbnail(item *rss.Item) string {
	ext := item.Extensions
	media, ok := ext["media"]
	if !ok {
		return ""
	}

	if media == nil {
		return ""
	}

	thumbnail, ok := media["thumbnail"]
	if !ok {
		return ""
	}

	if len(thumbnail) == 0 {
		return ""
	}

	attrs := thumbnail[0].Attrs
	thumb, ok := attrs["url"]
	if !ok {
		return ""
	}

	return thumb
}

func findFromMediaGroup(item *rss.Item) string {
	ext := item.Extensions
	media, ok := ext["media"]
	if !ok {
		return ""
	}

	if media == nil {
		return ""
	}

	group, ok := media["group"]
	if !ok {
		return ""
	}

	if len(group) == 0 {
		return ""
	}

	children := group[0].Children
	if children == nil {
		return ""
	}

	content, ok := children["content"]
	if !ok {
		return ""
	}

	if len(content) == 0 {
		return ""
	}

	attrs := content[0].Attrs
	thumb, ok := attrs["url"]
	if !ok {
		return ""
	}

	return thumb
}

func findFromMediaContent(item *rss.Item) string {
	ext := item.Extensions
	media, ok := ext["media"]
	if !ok {
		return ""
	}

	if media == nil {
		return ""
	}

	contents, ok := media["content"]
	if !ok {
		return ""
	}

	if len(contents) == 0 {
		return ""
	}

	attrs := contents[0].Attrs
	thumb, ok := attrs["url"]
	if !ok {
		return ""
	}

	return thumb
}

func findFromEnclosure(item *rss.Item) string {
	enclosure := item.Enclosure
	if enclosure == nil {
		return ""
	}
	return enclosure.URL

}
