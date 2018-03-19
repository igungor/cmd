package sub

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/igungor/cmd/filmdizibot/bot"
	"github.com/igungor/cmd/filmdizibot/command/sub/provider"
	"github.com/igungor/cmd/filmdizibot/config"
	"github.com/igungor/telegram"
)

type Sub struct {
	providers []provider.Provider
}

func New() *Sub {
	return &Sub{
		providers: []provider.Provider{
			provider.NewAddic7ed(
				config.Config.Addic7ed.Username,
				config.Config.Addic7ed.Password,
			),
		},
	}
}

func (s *Sub) Name() string { return "sub" }

func (s *Sub) Usage() string { return "sub <idx>\ndl <idx>" }

func (s *Sub) Match(msg string) bool { return strings.HasPrefix(msg, "sub") }

func (s *Sub) Run(ctx context.Context, b *bot.Bot, msg *bot.Message) {
	chat := msg.Chat.ID
	args := msg.Args()

	if len(args) == 0 {
		b.SendMessage(chat, "choose something")
		return
	}

	if args[0] == "dl" {
		err := s.handleDownload(ctx, b, msg, args[1:])
		if err != nil {
			b.SendMessage(chat, err.Error())
			return
		}
		b.SendMessage(chat, "subtitle is downloaded!")
		return
	}

	if len(args) > 1 {
		b.SendMessage(chat, "choose only one thing")
		return
	}

	subs, err := s.handleSelect(ctx, b, msg, args)
	if err != nil {
		b.SendMessage(chat, err.Error())
		return
	}

	if len(subs) == 0 {
		b.SendMessage(chat, "no subtitle found")
		return
	}

	sub := subs[0]
	title := fmt.Sprintf("Subtitles for '%v S%02dE%02d'", sub.Title, sub.Season, sub.Episode)
	separator := strings.Repeat("-", len(title))
	header := "idx | ep | title | lang | release"

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "```\n")
	fmt.Fprintln(&buf, title)
	fmt.Fprintln(&buf, separator)
	fmt.Fprintln(&buf, header)
	fmt.Fprintf(&buf, "```\n")
	for i, sub := range subs {
		fmt.Fprintf(&buf, "%02d %v\n", i+1, sub)
	}

	withMarkdown := telegram.WithParseMode(telegram.ModeMarkdown)
	b.SendMessage(chat, buf.String(), withMarkdown)
}

func (s *Sub) handleSelect(ctx context.Context, b *bot.Bot, msg *bot.Message, args []string) ([]*provider.Subtitle, error) {
	sel := args[0]
	idx, err := strconv.ParseInt(sel, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("wrong selection: %v", err)
	}

	children := b.Children()
	if len(children) < int(idx) {
		return nil, fmt.Errorf("wrong selection: no such idx")
	}

	f := children[idx-1]
	if f.IsDir() {
		return nil, fmt.Errorf("wrong selection: choose a file")
	}

	b.SelectFile(f)

	// FIXME: make it more generic
	subs, err := s.providers[0].Query(ctx, f.Name)
	if err != nil {
		return nil, fmt.Errorf("error listing subtitles: %v", err)
	}

	b.SetSubtitles(subs)
	return subs, nil
}

func (s *Sub) handleDownload(ctx context.Context, b *bot.Bot, msg *bot.Message, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("choose something")
	}

	idx, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		return fmt.Errorf("wrong selection: %v", err)
	}

	subs := b.Subtitles()
	if len(subs) < int(idx) {
		return fmt.Errorf("wrong selection: no such idx")
	}

	sub := subs[idx-1]
	rc, err := s.providers[0].Download(ctx, sub)
	if err != nil {
		return fmt.Errorf("error downloading subtitle: %v", err)
	}
	defer rc.Close()

	f := b.SelectedFile()
	ext := filepath.Ext(f.Name)
	fname := strings.TrimSuffix(f.Name, ext)
	fname = fmt.Sprintf("%v.srt", fname)

	err = b.UploadFile(ctx, rc, fname, f.ParentID)
	if err != nil {
		return fmt.Errorf("error uploading subtitle: %v", err)
	}
	return nil
}
