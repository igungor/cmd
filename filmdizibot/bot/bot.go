package bot

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"golang.org/x/oauth2"

	"github.com/igungor/cmd/filmdizibot/command/sub/provider"
	"github.com/igungor/cmd/filmdizibot/config"
	"github.com/igungor/go-putio/putio"
	"github.com/igungor/telegram"
)

const (
	groupWhatsup = -230439016
)

type Bot struct {
	*telegram.Bot
	p *putio.Client

	// mu guards below
	mu sync.Mutex

	// current working directory
	cwd putio.File

	// files in cwd
	children []putio.File

	// current subtitles
	cursubs []*provider.Subtitle

	// reference to last selected subtitle
	subsel putio.File
}

func New(ctx context.Context) (*Bot, error) {
	var (
		putioToken      = config.Config.Putio.Token
		telegramToken   = config.Config.Telegram.Token
		telegramWebhook = config.Config.Telegram.Webhook
	)
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: putioToken})
	oc := oauth2.NewClient(ctx, tokenSource)
	oc.Timeout = 10 * time.Second
	pc := putio.NewClient(oc)
	children, cwd, err := pc.Files.List(ctx, 0)
	if err != nil {
		return nil, fmt.Errorf("Error listing Putio root folder: %v", err)
	}

	bot := &Bot{
		Bot:      telegram.New(telegramToken),
		p:        pc,
		cwd:      cwd,
		children: children,
	}
	err = bot.SetWebhook(telegramWebhook)
	if err != nil {
		return nil, fmt.Errorf("Error setting webhook: %v", err)
	}
	return bot, nil
}

func (b *Bot) Messages() <-chan *Message {
	ch := make(chan *Message)
	go func() {
		for msg := range b.Bot.Messages() {
			ch <- &Message{msg}
		}
		close(ch)
	}()
	return ch
}

func (b *Bot) NewTransfer(ctx context.Context, url string) (putio.Transfer, error) {
	return b.p.Transfers.Add(ctx, url, -1, "")
}

func (b *Bot) ListFiles(ctx context.Context, parent int64) ([]putio.File, putio.File, error) {
	return b.p.Files.List(ctx, parent)
}

func (b *Bot) UploadFile(ctx context.Context, rc io.ReadCloser, filename string, parentid int64) error {
	_, err := b.p.Files.Upload(ctx, rc, filename, parentid)
	return err
}

func (b *Bot) Mkdir(ctx context.Context, dirname string, parentid int64) (*putio.File, error) {
	f, err := b.p.Files.CreateFolder(ctx, dirname, parentid)
	return &f, err
}

func (b *Bot) MoveFiles(ctx context.Context, parentid int64, files ...int64) error {
	return b.p.Files.Move(ctx, parentid, files...)
}

func (b *Bot) ConvertToMP4(ctx context.Context, id int64) error {
	return b.p.Files.Convert(ctx, id)
}

func (b *Bot) CWD() putio.File {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.cwd
}

func (b *Bot) SetCWD(cwd putio.File) {
	b.mu.Lock()
	b.cwd = cwd
	b.mu.Unlock()
}

func (b *Bot) Children() []putio.File {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.children
}
func (b *Bot) SetChildren(c []putio.File) {
	b.mu.Lock()
	b.children = c
	b.mu.Unlock()
}

func (b *Bot) SelectedFile() putio.File {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.subsel
}

func (b *Bot) SelectFile(f putio.File) {
	b.mu.Lock()
	b.subsel = f
	b.mu.Unlock()
}

func (b *Bot) Subtitles() []*provider.Subtitle {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.cursubs
}

func (b *Bot) SetSubtitles(subs []*provider.Subtitle) {
	b.mu.Lock()
	b.cursubs = subs
	b.mu.Unlock()
}
