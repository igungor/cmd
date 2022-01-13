package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	userAgent   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36"
	timeLayout  = "02/01/2006"
	storagePath = ".local/share"
	baseURL     = "https://www.tefas.gov.tr/FonAnaliz.aspx?FonKod=%v"
)

var (
	ErrDisabled = fmt.Errorf("disabled on weekends")
)

func GetFunds(ctx context.Context, codes ...string) ([]Fund, error) {

	c := &http.Client{Timeout: time.Minute}

	today := time.Now()

	switch today.Weekday() {
	case 6, 0: // saturday and sunday
		return nil, ErrDisabled
	}

	var funds []Fund
	for _, code := range codes {
		code = strings.ToUpper(code)

		u := fmt.Sprintf(baseURL, code)

		req, _ := http.NewRequest("GET", u, nil)
		req = req.WithContext(ctx)
		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

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

func atof(s string) float64 {
	s = strings.TrimPrefix(s, "%")
	s = strings.ReplaceAll(s, ",", ".")
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
