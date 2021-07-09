package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/peak/picolo"
)

var dictionaries = []dictionary{
	{
		name: "Büyük Türkçe Sözlük",
		url:  "https://sozluk.gov.tr/gts",
	},
	// {name: "Atasözleri ve Deyimler Sözlüğü", url: "https://sozluk.gov.tr/atasozu"},
	// {name: "bati", url: ""},
	// {name: "terim", url: ""},
	// {name: "hemsirelik", url: ""},
	// {name: "eczacilik", url: ""},
	// {name: "metroloji", url: ""},
	// {name: "yazim", url: ""},
}

var client = http.Client{Timeout: 10 * time.Second}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	term := flag.Arg(0)

	logger := picolo.New()

	var wg sync.WaitGroup
	wg.Add(len(dictionaries))

	resultch := make(chan result, len(dictionaries))
	for _, dict := range dictionaries {
		dict := dict
		go func() {
			defer wg.Done()

			result, err := query(term, dict)
			if err != nil {
				logger.Errorf("%v: %v", dict, err)
				return
			}

			resultch <- result
		}()
	}

	wg.Wait()
	close(resultch)

	var buf bytes.Buffer
	for result := range resultch {
		fmt.Fprintln(&buf, result.String())
	}

	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(120),
	)

	out, err := r.Render(buf.String())
	if err != nil {
		logger.Errorf("render: %v", err)
		return
	}

	fmt.Print(out)
}

func query(term string, dict dictionary) (result, error) {
	result := result{dict: dict}

	u, err := url.Parse(dict.url)
	if err != nil {
		return result, err
	}

	params := u.Query()
	params.Set("ara", term)
	u.RawQuery = params.Encode()

	resp, err := client.Get(u.String())
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	var r tdkResponse
	if err := json.Unmarshal(body, &r); err != nil {
		var errResponse struct {
			Error string
		}
		if eerr := json.Unmarshal(body, &errResponse); eerr != nil {
			return result, err
		} else {
			return result, fmt.Errorf(errResponse.Error)
		}

		return result, err
	}

	result.meanings = r
	return result, nil
}

type result struct {
	dict     dictionary
	meanings tdkResponse
}

func (r result) String() string {
	if len(r.meanings) == 0 {
		return ""
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "# %v\n", r.dict.name)

	for _, meaning := range r.meanings[0].AnlamlarListe {
		fmt.Fprintf(&buf, "%v. %v\n", meaning.AnlamSira, meaning.Anlam)
	}

	return buf.String()
}

type dictionary struct {
	name string
	url  string
}

type tdkResponse []struct {
	AnlamlarListe []struct {
		Anlam     string `json:"anlam"`
		AnlamSira string `json:"anlam_sira"`
	} `json:"anlamlarListe,omitempty"`
}
