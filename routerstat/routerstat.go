// routerstat is a simple CLI tool that extracts router connections from my
// router's ugly and unintuitive web interface.
//
// it reads username and password from ROUTER_USERNAME and ROUTER_PASSWORD
// environment variables.
//
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf"
)

var (
	routerURL       = "http://192.168.100.1/"
	routerAccessURL = "http://192.168.100.1/wlanAccess.asp"
	routerUsername  = os.Getenv("ROUTER_USERNAME")
	routerPassword  = os.Getenv("ROUTER_PASSWORD")
)

var showStale = flag.Bool("a", false, "show all (including stale) connections")

type row struct {
	mac   string
	age   time.Duration
	dbm   string
	ip    string
	host  string
	mode  string
	speed string
}

type ByAge []row

func (r ByAge) Len() int           { return len(r) }
func (r ByAge) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ByAge) Less(i, j int) bool { return r[i].age < r[j].age }

func extract() ([]row, error) {
	b := surf.NewBrowser()
	if err := b.Open(routerURL); err != nil {
		return nil, err
	}

	form, _ := b.Form("form")
	form.Input("loginUsername", routerUsername)
	form.Input("loginPassword", routerPassword)
	form.Submit()

	if err := b.Open(routerAccessURL); err != nil {
		return nil, err
	}

	var table []string
	b.Find("table").Each(func(i int, s *goquery.Selection) {
		if i != 3 { // wlan connections table. has no id or class
			return
		}
		s.Find("td").Each(func(_ int, td *goquery.Selection) {
			table = append(table, strings.TrimSpace(td.Text()))
		})
	})

	const c = 7 // column count
	rowcount := len(table) / c

	t := table
	var rows []row
	for i := 1; i < rowcount; i++ {
		newrow := row{
			mac:   t[0+i*c],
			dbm:   t[2+i*c],
			ip:    t[3+i*c],
			host:  t[4+i*c],
			mode:  t[5+i*c],
			speed: t[6+i*c],
		}
		age, _ := time.ParseDuration(t[1+i*c] + "s")
		newrow.age = age

		// trim long hostnames
		if len(newrow.host) > 14 {
			newrow.host = newrow.host[:14] + "…"
		}

		rows = append(rows, newrow)
	}

	sort.Sort(ByAge(rows))

	return rows, nil
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("routerstat: ")

	flag.Parse()

	rows, err := extract()
	if err != nil {
		log.Fatal(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	defer w.Flush()

	fmt.Fprintln(w, "Hostname\tIP Address\tSpeed(kbps)\tAge")
	fmt.Fprintln(w, "--------\t----------\t-----------\t---")
	for _, r := range rows {
		if !*showStale {
			if r.age.Hours() > 1 {
				continue
			}
		}
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", r.host, r.ip, r.speed, r.age)
	}
}
