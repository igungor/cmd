package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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
	const (
		baseurl = "http://www.akportfoy.com.tr/ajax/getfundreturns?fundsubtypeId=%v&enddate=%v&lang=tr"
	)

	c := &http.Client{Timeout: time.Minute}

	today := time.Now()

	switch today.Weekday() {
	case 6, 0: // saturday and sunday
		return nil, ErrDisabled
	}

	const fund = YabanciHisseSenedi
	u := fmt.Sprintf(baseurl, fund, today.Format(timeLayout))
	req, _ := http.NewRequest("POST", u, nil)
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response struct {
		Title string
		Table string
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.Table))
	if err != nil {
		return nil, err
	}

	atof := func(s string) float64 {
		f, _ := strconv.ParseFloat(s, 64)
		return f
	}

	var funds []Fund
	doc.Find("tr").Each(func(i int, sel *goquery.Selection) {
		if i == 0 {
			return
		}

		code := sel.Find(".fundcode").Text()
		name := sel.Find("th").Text()
		name = strings.TrimPrefix(name, code)
		name = strings.TrimSpace(name)

		if len(codes) != 0 {
			var found bool
			for _, c := range codes {
				if strings.ToLower(c) == strings.ToLower(code) {
					found = true
					break
				}
			}
			if !found {
				return
			}
		}

		fund := Fund{
			Code: code,
			Name: name,
		}

		sel.Children().Each(func(n int, sel *goquery.Selection) {
			switch n {
			case 1:
				fund.Price = atof(sel.Text())
			case 2:
				fund.Daily = atof(sel.Text())
			case 3:
				fund.Weekly = atof(sel.Text())
			case 4:
				fund.Monthly = atof(sel.Text())
			case 5:
				fund.Annual = atof(sel.Text())
			}
		})

		funds = append(funds, fund)

	})

	return funds, nil
}

type Fund struct {
	Type    FundType
	Code    string
	Name    string
	Price   float64
	Daily   float64
	Weekly  float64
	Monthly float64
	Annual  float64
}

type FundType uint8

const (
	ParaPiyasasi            FundType = 5
	BorclanmaAraclari       FundType = 6
	Katilim                 FundType = 4
	Degisken                FundType = 2
	HisseSenedi             FundType = 1
	YabanciHisseSenedi      FundType = 7
	DegerliMaden            FundType = 3
	FonSepeti               FundType = 9
	GayrimenkulYatirim      FundType = 32
	GirisimSermayesiYatirim FundType = 30
	DovizSerbest            FundType = 71
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

		name := strings.TrimPrefix(f.Name, "Ak Portföy ")
		name = strings.TrimSuffix(name, "Yabancı Hisse Senedi Fonu")
		name = strings.TrimSpace(name)

		fmt.Fprintf(&buf, header, f.Code, name, freefall, f.Code)
		fmt.Fprintf(&buf, "• Fiyat:  %v | color=%v size=11\n", f.Price, "black")
		fmt.Fprintf(&buf, "• Günlük:  %v%% | color=%v size=11\n", f.Daily, color(f.Daily))
		fmt.Fprintf(&buf, "• Haftalık:  %v%% | color=%v size=11\n", f.Weekly, color(f.Weekly))
		fmt.Fprintf(&buf, "• Aylık: %v%% | color=%v size=11\n", f.Monthly, color(f.Monthly))
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
	fpath := filepath.Join(path, "funds.json")

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
