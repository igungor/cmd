package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"golang.org/x/net/publicsuffix"
)

var (
	logger = log.New(os.Stdout, "", 0)
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: urlparse <url>\n")
	}
	flag.Parse()

	var arg string
	if flag.NArg() == 1 {
		arg = flag.Arg(0)
		logger.Println(tldplusone(arg))
		return
	}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		line := s.Text()
		logger.Println(tldplusone(line))
	}
}

func tldplusone(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		logger.Fatal(err)
	}

	v, err := publicsuffix.EffectiveTLDPlusOne(u.Hostname())
	if err != nil {
		logger.Fatal(err)
	}
	return v
}
