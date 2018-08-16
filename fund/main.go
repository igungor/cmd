package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36"

func main() {
	var (
		flagBitbar = flag.Bool("bitbar", false, "Enable bitbar compatible output")
	)
	flag.Parse()

	funds, err := GetFunds(flag.Args()...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(prettyPrint(*flagBitbar, funds...))
}

func GetFund(code string) (Fund, error) {
	code = strings.ToUpper(code)

	funds, err := GetFunds()
	if err != nil {
		return Fund{}, err
	}

	for _, fund := range funds {
		if fund.Code == code {
			return fund, nil
		}
	}
	return Fund{}, fmt.Errorf("fund not found")
}

func GetFunds(codes ...string) ([]Fund, error) {
	const (
		baseurl    = "http://www.akportfoy.com.tr/ajax/getfundreturns?fundsubtypeId=%v&enddate=%v&lang=tr"
		timelayout = "02/01/2006"
	)

	c := &http.Client{Timeout: time.Minute}

	const fund = YabanciHisseSenedi
	today := time.Now().Format(timelayout)

	u := fmt.Sprintf(baseurl, fund, today)
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
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
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

func prettyPrint(bitbar bool, funds ...Fund) string {
	if bitbar {
		return printBitbar(funds...)
	}

	return printLong(funds...)
}

func printBitbar(funds ...Fund) string {
	color := func(f float64) string {
		if f < 0 {
			return "red"
		}
		return "green"
	}

	var buf bytes.Buffer

	format := "%v %v (%v%%) | color=%v href=http://www.tefas.gov.tr/FonAnaliz.aspx?FonKod=%v\n"

	for _, f := range funds {
		fmt.Fprintf(&buf, format, f.Code, f.Price, f.Daily, color(f.Daily), f.Code)
		fmt.Fprintf(&buf, "----\n")
		fmt.Fprintf(&buf, "-- %v | color=black\n", f.Name)
		fmt.Fprintf(&buf, "-----\n")
		fmt.Fprintf(&buf, "-- Weekly:  %v%% | color=%v\n", f.Weekly, color(f.Weekly))
		fmt.Fprintf(&buf, "-- Monthly: %v%% | color=%v\n", f.Monthly, color(f.Monthly))
		fmt.Fprintf(&buf, "-- Annual:  %v%% | color=%v\n", f.Annual, color(f.Annual))
	}
	return buf.String()
}

func printLong(funds ...Fund) string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 2, 2, ' ', 0)

	format := "%v\t%v\t%v\t%v\t%v\t%v\t%v\n"

	fmt.Fprintf(w, format, "code", "name", "price", "daily", "weekly", "monthly", "annual")
	for _, f := range funds {
		fmt.Fprintf(w, format, f.Code, f.Name, f.Price, f.Daily, f.Weekly, f.Monthly, f.Annual)
	}
	w.Flush()

	return buf.String()
}
