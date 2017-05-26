package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/tabwriter"
)

var (
	verbose = flag.Bool("v", false, "enable verbose mode")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("subs: ")

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
	}

	q := flag.Arg(0)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	provider := NewAddic7ed()

	if !provider.Available() {
		log.Fatalf("Provider %v is not available\n", provider)
	}

	subs, err := provider.Query(ctx, q)
	if err != nil {
		log.Fatalf("query failed for provider %q: %v\n", provider, err)
	}

	if len(subs) == 0 {
		log.Fatalf("No subtitle found for %q\n", q)
	}

	tw := tabwriter.NewWriter(os.Stdout, 4, 4, 1, ' ', tabwriter.Debug)
	tw.Write([]byte("S\tE\tTitle\tLanguage\tRelease\tHI\n"))
	for _, sub := range subs {
		io.Copy(tw, strings.NewReader(sub.String()))
	}

	err = tw.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "subs <argument>\n")
	fmt.Fprintf(os.Stderr, "\targument: a video file\n")
	fmt.Fprintf(os.Stderr, "\targument: an episode title\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func debugf(format string, v ...interface{}) {
	if *verbose {
		log.Printf(format, v...)
	}
}
