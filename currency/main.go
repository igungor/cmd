package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36"

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

		fmt.Fprintf(&buf, "%v <%.1f%%> | size=13 color=%v href=href=https://tr.investing.com/currencies/%v-try-commentary\n", code, c.ChangeRate, color(c.ChangeRate), code)
		fmt.Fprintf(&buf, "• Sell: %.4f | size=11 color=black\n", c.Selling)
		fmt.Fprintf(&buf, "• Buy: %.4f | size=11 color=black\n", c.Buying)
		fmt.Fprintf(&buf, "---\n")
	}
	return buf.String()
}
