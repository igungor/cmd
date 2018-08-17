package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/tabwriter"
	"time"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36"

func main() {
	var (
		flagBitbar = flag.Bool("bitbar", false, "Enable bitbar compatible output")
	)
	_ = flagBitbar
	flag.Parse()

	currencies, err := GetCurrencies(flag.Args()...)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(prettyPrint(*flagBitbar, currencies...))
}

func GetCurrencies(codes ...string) ([]Currency, error) {
	c := &http.Client{Timeout: 10 * time.Second}

	const baseurl = "https://www.doviz.com/api/v1/currencies/all/latest"

	resp, err := c.Get(baseurl)
	if err != nil {
		return nil, fmt.Errorf("could not fetch currencies: %v", err)
	}
	defer resp.Body.Close()

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

func prettyPrint(bitbar bool, currencies ...Currency) string {
	if bitbar {
		return printBitbar(currencies...)
	}

	return printLong(currencies...)
}

func printBitbar(currencies ...Currency) string {
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

	format := "\x1b[30;1;8m%v\x1b[0m | ansi=true size=13 href=href=https://tr.investing.com/currencies/%v-try-commentary\n"

	for _, c := range currencies {
		fmt.Fprintf(&buf, format, strings.ToUpper(c.Code), strings.ToLower(c.Code))
		fmt.Fprintf(&buf, "• Fiyat:  %v | color=%v size=11\n", c.Selling, "black")
		fmt.Fprintf(&buf, "• Günlük:  %.2f%% | color=%v size=11\n", c.ChangeRate, color(c.ChangeRate))
		fmt.Fprintf(&buf, "---\n")
	}
	return buf.String()
}

func printLong(currencies ...Currency) string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 2, 2, ' ', 0)

	format := "%v\t%.4f\t%.1f\n"

	fmt.Fprintf(w, "currency\tprice\tdaily\n")
	for _, c := range currencies {
		fmt.Fprintf(w, format, c.Code, c.Selling, c.ChangeRate)
	}
	w.Flush()

	return buf.String()
}
