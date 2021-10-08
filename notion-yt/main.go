package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dstotijn/go-notion"
	"github.com/pelletier/go-toml"
)

func main() {
	if err := realmain(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Config struct {
	Currencies map[string]string
	Funds      map[string]string
}

func realmain() error {
	var (
		flagConfig = flag.String("c", "config.toml", "TOML configuration path")
	)
	flag.Parse()

	logger := log.New(os.Stdout, "", 0)

	bytes, err := ioutil.ReadFile(*flagConfig)
	if err != nil {
		return err
	}

	var cfg Config
	if err := toml.Unmarshal(bytes, &cfg); err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	secret := os.Getenv("NOTION_YT_SECRET")
	if secret == "" {
		return fmt.Errorf("NOTION_YET_SECRET environment variable is not set")
	}

	funds := cfg.Funds
	currencies := cfg.Currencies

	var fundcodes []string
	for code := range funds {
		fundcodes = append(fundcodes, code)
	}

	fundresults, err := GetFunds(ctx, fundcodes...)
	if err != nil {
		if err == ErrDisabled {
			return nil
		}
		return err
	}

	c := notion.NewClient(secret)

	for _, r := range fundresults {
		pageID := funds[r.Code]

		if r.Price == 0 {
			continue
		}

		price := fmt.Sprintf("%.6f", r.Price)
		params := notion.UpdatePageParams{
			DatabasePageProperties: &notion.DatabasePageProperties{
				"Current Price ₺": notion.DatabasePageProperty{
					Type: notion.DBPropTypeRichText,
					RichText: []notion.RichText{
						{
							Type: notion.RichTextTypeText,
							Text: &notion.Text{Content: price},
						},
					},
				},
			},
		}

		logger.Printf("Updating %q...", r.Code)
		_, err := c.UpdatePageProps(ctx, pageID, params)
		if err != nil {
			return err
		}
	}

	for code, pageID := range currencies {
		v, err := request(code)
		if err != nil {
			return err
		}

		price := fmt.Sprintf("%.6f", v)
		params := notion.UpdatePageParams{
			DatabasePageProperties: &notion.DatabasePageProperties{
				"Current Price ₺": notion.DatabasePageProperty{
					Type: notion.DBPropTypeRichText,
					RichText: []notion.RichText{
						{
							Type: notion.RichTextTypeText,
							Text: &notion.Text{Content: price},
						},
					},
				},
			},
		}

		logger.Printf("Updating %q...", code)
		_, err = c.UpdatePageProps(ctx, pageID, params)
		if err != nil {
			return err
		}
	}

	return nil
}
