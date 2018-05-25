package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"sort"
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	urlstr := flag.Arg(0)

	u, err := url.Parse(urlstr)
	if err != nil {
		log.Fatal(err)
	}

	pprint(u)
}

func pprint(u *url.URL) {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "Scheme: %v\n", u.Scheme)
	fmt.Fprintf(&buf, "Host: %v\n", u.Hostname())
	fmt.Fprintf(&buf, "Port: %v\n", u.Port())
	fmt.Fprintf(&buf, "User: %v\n", u.User)
	fmt.Fprintf(&buf, "Path: %v\n", u.Path)

	fmt.Fprintf(&buf, "Query:\n")
	var keys []string
	for key := range u.Query() {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		fmt.Fprintf(&buf, "  %v: %v\n", key, u.Query().Get(key))
	}
	fmt.Fprintf(&buf, "Fragment: %v\n", u.Fragment)

	fmt.Println(buf.String())
}
