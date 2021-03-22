package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	userAgent   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36"
	timeLayout  = "02/01/2006"
	storagePath = ".local/share"
)

var (
	ErrDisabled = fmt.Errorf("disabled on weekends")
)

func main() {
	flag.Parse()

	funds, err := GetFunds(flag.Args()...)
	if err == ErrDisabled {
		fmt.Println(err)
		os.Exit(0)
	}

	fmt.Print(prettyPrint(funds...))
}

func GetFunds(codes ...string) ([]Fund, error) {
	const baseurl = "https://www.tefas.gov.tr/FonAnaliz.aspx?FonKod=%v"

	c := &http.Client{Timeout: time.Minute}

	today := time.Now()

	switch today.Weekday() {
	case 6, 0: // saturday and sunday
		return nil, ErrDisabled
	}

	var funds []Fund
	for _, code := range codes {
		code = strings.ToUpper(code)

		u := fmt.Sprintf(baseurl, code)

		req, _ := http.NewRequest("GET", u, nil)
		req.Header.Set("User-Agent", userAgent)

		resp, err := c.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Printf("unexpected status code: %v", resp.StatusCode)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return nil, err
		}

		fundName := strings.TrimSpace(doc.Find(".main-indicators h2").Text())
		fund := Fund{
			Code: code,
			Name: fundName,
		}

		fund.Price = atof(doc.Find(".top-list > li:nth-child(1) span").Text())
		fund.Daily = atof(doc.Find(".top-list > li:nth-child(2) span").Text())

		doc.Find(".price-indicators li span").Each(func(i int, sel *goquery.Selection) {
			changePercent := atof(sel.Text())

			switch i {
			case 0:
				fund.Monthly = changePercent
			case 1:
				fund.Quarterly = changePercent
			case 2:
				fund.Biannual = changePercent
			case 3:
				fund.Annual = changePercent
			}

		})
		funds = append(funds, fund)
	}

	return funds, nil
}

type Fund struct {
	Type      FundType
	Code      string
	Name      string
	Price     float64
	Daily     float64
	Monthly   float64
	Quarterly float64
	Biannual  float64
	Annual    float64
}

type FundType uint8

const (
	CommodityFunds     FundType = 3
	FixedIncomeFunds   FundType = 6
	ForeignEquityFunds FundType = 111
)

func prettyPrint(funds ...Fund) string {
	color := func(f float64) string {
		switch {
		case f == 0:
			return "grey"
		case f < 0:
			return "red"
		case f > 0:
			return "green"
		}
		return "black"
	}

	var buf bytes.Buffer

	header := "%v (%v) \x1b[31;1;8m%v\x1b[0m | size=13 href=https://www.tefas.gov.tr/FonAnaliz.aspx?FonKod=%v\n"

	fmt.Fprintln(&buf, "Refresh | refresh=true")
	fmt.Fprintf(&buf, "---\n")

	for _, f := range funds {
		sethop(f.Code, time.Now(), f.Daily)

		// calculate the fall
		hop := gethop(f.Code)
		var freefall string
		for _, h := range hop {
			if h.DailyChange < 0 {
				freefall += "◉ "
			} else {
				freefall = ""
			}
		}

		fmt.Fprintf(&buf, header, f.Code, f.Name, freefall, f.Code)
		fmt.Fprintf(&buf, "• Fiyat:  %v | size=11\n", f.Price)
		fmt.Fprintf(&buf, "• Günlük:  %v%% | color=%v size=11\n", f.Daily, color(f.Daily))
		fmt.Fprintf(&buf, "• Aylık: %v%% | color=%v size=11\n", f.Monthly, color(f.Monthly))
		fmt.Fprintf(&buf, "• 3 Aylık: %v%% | color=%v size=11\n", f.Quarterly, color(f.Quarterly))
		fmt.Fprintf(&buf, "• 6 Aylık: %v%% | color=%v size=11\n", f.Biannual, color(f.Biannual))
		fmt.Fprintf(&buf, "• Yıllık:  %v%% | color=%v size=11\n", f.Annual, color(f.Annual))
		fmt.Fprintf(&buf, "---\n")
	}
	return buf.String()
}

func sethop(code string, date time.Time, change float64) {
	// change is not reflected on the site yet, hence the zero value.
	if change == 0 {
		return
	}
	home := os.Getenv("HOME")
	path := filepath.Join(home, storagePath)
	os.MkdirAll(path, 0755)
	fpath := filepath.Join(path, "funds.json")

	_, err := os.Stat(fpath)
	if os.IsNotExist(err) {
		ioutil.WriteFile(fpath, []byte("{}"), 0644)
	}

	allhops := gethops()

	datestr := date.Format(timeLayout)
	hop := Hop{Date: datestr, DailyChange: change}

	hops, ok := allhops[code]
	if !ok {
		hops = []Hop{hop}
	} else {
		var found bool
		for _, h := range hops {
			if h.Date == datestr {
				found = true
				break
			}
		}
		if !found {
			hops = append(hops, hop)
		}
	}

	// limit the size of the slice to 3
	if len(hops) >= 3 {
		hops = hops[len(hops)-3:]
	}

	allhops[code] = hops

	b, _ := json.Marshal(allhops)
	_ = ioutil.WriteFile(fpath, b, 0644)
}

func gethops() map[string][]Hop {
	home := os.Getenv("HOME")
	path := filepath.Join(home, storagePath)
	os.MkdirAll(path, 0755)
	fpath := filepath.Join(path, "funds.json")

	m := make(map[string][]Hop)

	content, _ := ioutil.ReadFile(fpath)
	_ = json.Unmarshal(content, &m)

	return m
}

func gethop(code string) []Hop {
	return gethops()[code]
}

type Hop struct {
	Date        string
	DailyChange float64
}

func atof(s string) float64 {
	s = strings.TrimPrefix(s, "%")
	s = strings.ReplaceAll(s, ",", ".")
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
