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
	"strings"
	"time"
)

const (
	userAgent   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36"
	timeLayout  = "02/01/2006"
	storagePath = ".local/share"
)

func main() {
	flag.Parse()

	currencies, err := GetCurrencies(flag.Args()...)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(prettyPrint(currencies...))
}

func GetCurrencies(codes ...string) ([]Currency, error) {
	c := &http.Client{Timeout: 10 * time.Second}

	const baseurl = "https://www.doviz.com/api/v1/currencies/all/latest"

	req, err := http.NewRequest("GET", baseurl, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %v", err)
	}
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not fetch currencies: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected HTTP status code: %v", resp.StatusCode)
	}

	var currencies []Currency
	if err := json.NewDecoder(resp.Body).Decode(&currencies); err != nil {
		return nil, fmt.Errorf("could not decode json: %v", err)
	}

	if len(codes) == 0 {
		return currencies, nil
	}

	var selection []Currency
	for _, curr := range currencies {
		currcode := strings.ToLower(curr.Code)
		for _, code := range codes {
			code = strings.ToLower(code)
			if currcode == code {
				selection = append(selection, curr)
			}
		}

	}
	return selection, nil
}

type Currency struct {
	Selling    float64 `json:"selling"`
	UpdateDate float64 `json:"update_date"`
	Currency   int     `json:"currency"`
	Buying     float64 `json:"buying"`
	ChangeRate float64 `json:"change_rate"`
	Name       string  `json:"name"`
	FullName   string  `json:"full_name"`
	Code       string  `json:"code"`
}

func prettyPrint(currencies ...Currency) string {
	const (
		boldblack = "\x1b[30m"
		black     = "\x1b[1;30m"
		red       = "\x1b[35m"
		green     = "\x1b[32m"
		reset     = "\x1b[0m"
	)
	color := func(f float64) string {
		switch {
		case f < 0:
			return "red"
		case f > 0:
			return "green"
		default:
			return "black"
		}
	}

	var buf bytes.Buffer

	for _, c := range currencies {
		code := strings.ToUpper(c.Code)

		sethop(code, time.Now(), c.ChangeRate)

		// calculate the fall
		hop := gethop(code)
		var freefall string
		for _, h := range hop {
			if h.DailyChange < 0 {
				freefall += "◉ "
			} else {
				freefall = ""
			}
		}

		fmt.Fprintf(&buf, "%v <%.2f%%> \x1b[31;1;8m%v\x1b[0m | size=13 color=%v href=https://tr.investing.com/currencies/%v-try-commentary\n", code, c.ChangeRate, freefall, color(c.ChangeRate), strings.ToLower(code))
		fmt.Fprintf(&buf, "• Al: %.4f | size=11 color=black\n", c.Buying)
		fmt.Fprintf(&buf, "• Sat: %.4f | size=11 color=black\n", c.Selling)
		fmt.Fprintf(&buf, "---\n")
	}
	return buf.String()
}

func sethop(code string, date time.Time, change float64) {
	home := os.Getenv("HOME")
	path := filepath.Join(home, storagePath)
	os.MkdirAll(path, 0755)
	fpath := filepath.Join(path, "currency.json")

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
		hops = hops[len(hops)-3 : len(hops)]
	}

	allhops[code] = hops

	b, _ := json.Marshal(allhops)
	_ = ioutil.WriteFile(fpath, b, 0644)
}

func gethops() map[string][]Hop {
	home := os.Getenv("HOME")
	path := filepath.Join(home, storagePath)
	os.MkdirAll(path, 0755)
	fpath := filepath.Join(path, "currency.json")

	m := make(map[string][]Hop)

	content, _ := ioutil.ReadFile(fpath)
	_ = json.Unmarshal(content, &m)

	return m
}

func gethop(code string) []Hop {
	all := gethops()
	v, _ := all[code]
	return v
}

type Hop struct {
	Date        string
	DailyChange float64
}
