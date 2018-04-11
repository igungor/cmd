package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	humanize "github.com/dustin/go-humanize"
	"github.com/gorilla/schema"
	"github.com/igungor/cmd/filmdizibot/bot"
	"github.com/igungor/cmd/filmdizibot/command"
	"github.com/igungor/cmd/filmdizibot/config"
	"github.com/igungor/telegram"
)

const groupWhatsup = -230439016

func main() {
	var (
		flagConfig = flag.String("c", "filmdizibot.toml", "Path to configuration file")
	)
	flag.Parse()

	log.SetFlags(0)
	log.SetPrefix("filmdizibot: ")

	if err := config.Load(*flagConfig); err != nil {
		log.Fatalf("could not load configuration: %v", err)
	}

	ctx := context.Background()
	bot, err := bot.New(ctx)
	if err != nil {
		log.Fatalf("Error creating the bot: %v\n", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", bot.Handler())
	mux.HandleFunc("/cb", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			bot.SendMessage(groupWhatsup, fmt.Sprintf("ParseForm failed: %v", err))
			return
		}

		dec := schema.NewDecoder()
		dec.IgnoreUnknownKeys(true)
		dec.SetAliasTag("json")

		var t transfer
		err = dec.Decode(&t, r.PostForm)
		if err != nil {
			bot.SendMessage(groupWhatsup, fmt.Sprintf("Decode failed: %v", err))
			return
		}

		// ignore spam requests
		if t.Name == "" {
			return
		}

		md := telegram.WithParseMode(telegram.ModeMarkdown)
		txt := fmt.Sprintf("🗣 New file downloaded!\n\n*%v*\n\nSize: %v", t.Name, humanize.Bytes(uint64(t.Size)))
		bot.SendMessage(groupWhatsup, txt, md)
	})

	go func() {
		addr := net.JoinHostPort(config.Config.Host, strconv.Itoa(config.Config.Port))
		log.Fatal(http.ListenAndServe(addr, mux))
	}()

	for msg := range bot.Messages() {
		if msg.IsService() {
			log.Printf("Skipping service message...\n")
			continue
		}

		cmdname := msg.Command()
		cmd := command.Match(cmdname)
		if cmd == nil {
			continue
		}

		log.Printf("New request: %v\n", msg.Text)
		go cmd.Run(ctx, bot, msg)
	}
}

type transfer struct {
	Name       string `json:"name"`
	Size       int    `json:"size"`
	FileID     int64  `json:"file_id"`
	DownloadID int64  `json:"download_id"`
	ParentID   int64  `json:"save_parent_id"`
}

func (t transfer) String() string {
	return fmt.Sprintf("%q\n\n indirilmeye baslandi.\n Boyut: **%v**\n", t.Name, humanize.Bytes(uint64(t.Size)))
}
