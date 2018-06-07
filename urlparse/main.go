package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: urlparse <url>\n")
	}
	flag.Parse()

	var r io.Reader = os.Stdin
	if arg := flag.Arg(0); arg != "" {
		r = strings.NewReader(arg)
	}

	var buf bytes.Buffer
	io.Copy(&buf, r)

	u, err := url.Parse(buf.String())
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

	w := tabwriter.NewWriter(&buf, 0, 2, 2, ' ', 0)
	for _, key := range keys {
		fmt.Fprintf(w, "  %v:\t%v\n", key, u.Query().Get(key))
	}
	w.Flush()
	fmt.Fprintf(&buf, "Fragment: %v\n", u.Fragment)

	fmt.Println(buf.String())
}
