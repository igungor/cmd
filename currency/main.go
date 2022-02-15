package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	yahooFinanceURL = "https://query1.finance.yahoo.com/v8/finance/chart/"
)

func main() {
	if err := realmain(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func realmain() error {
	fmt.Println(request("EUR"))

	return nil
}

func request(code string) (float64, error) {
	const to = "TRY=X"

	u, _ := url.Parse(yahooFinanceURL)
	u.Path += fmt.Sprintf("%v%v", code, to)
	params := u.Query()
	params.Set("range", "1d")
	u.RawQuery = params.Encode()

	c := http.Client{Timeout: time.Minute}

	resp, err := c.Get(u.String())
	if err != nil {
		return 0, fmt.Errorf("yahoo: could not fetch response: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("yahoo: unexpected status code %v", resp.StatusCode)
	}

	var response yahooFinanceResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return 0, fmt.Errorf("yahoo: parse json: %v", err)
	}

	result := response.Chart.Result
	quote := result[len(result)-1].Indicators.Quote
	closed := quote[len(quote)-1].Close
	prevClose := result[len(result)-1].Meta.PreviousClose

	if len(closed) == 0 {
		// some currencies are not available in 'close data', such as BGN.
		// dont let people down.
		if prevClose != 0 {
			return prevClose, nil
		}
		return 0, fmt.Errorf("yahoo: no value found for code %q", code)
	}

	var rates []float64
	for _, v := range closed {
		rate, ok := v.(float64)
		// skip unrecognized values to a list for later use
		if !ok {
			continue
		}
		rates = append(rates, rate)
	}
	return rates[len(rates)-1], nil
}

type yahooFinanceResponse struct {
	Chart struct {
		Error  interface{} `json:"error"`
		Result []struct {
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Close  []interface{} `json:"close"`
					High   []interface{} `json:"high"`
					Low    []interface{} `json:"low"`
					Open   []interface{} `json:"open"`
					Volume []interface{} `json:"volume"`
				} `json:"quote"`
			} `json:"indicators"`
			Meta struct {
				PreviousClose float64 `json:"previousClose"`
			} `json:"meta"`
		} `json:"result"`
	} `json:"chart"`
}
