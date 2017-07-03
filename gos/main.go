package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const godocURL = "https://api.godoc.org/search?q="

func main() {
	log.SetFlags(0)

	// flags
	var (
		flagCount = flag.Int("n", 5, "number of results")
	)

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
	}

	client := &http.Client{Timeout: 20 * time.Second}

	query := flag.Arg(0)
	resp, err := client.Get(godocURL + url.QueryEscape(query))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("not ok")
	}

	var r result
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		log.Fatal(err)
	}

	if len(r.Results) == 0 {
		log.Fatal("no result found")
	}

	n := min(*flagCount, len(r.Results))
	for _, result := range r.Results[:n] {
		fmt.Printf("\n\033[1m%v\033[m\n", result.Path)
		syn := result.Synopsis
		if syn != "" {
			fmt.Printf("  %v\n", syn)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type result struct {
	Results []struct {
		Path     string
		Synopsis string
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  gos [query]\n")
	flag.PrintDefaults()
	os.Exit(1)
}
