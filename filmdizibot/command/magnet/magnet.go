package magnet

import (
	"context"
	"fmt"
	"strings"

	"github.com/igungor/cmd/filmdizibot/bot"
	"github.com/igungor/telegram"
)

type Magnet struct{}

func New() *Magnet { return &Magnet{} }

func (m *Magnet) Name() string { return "magnet" }

func (m *Magnet) Usage() string { return "foobar" }

func (m *Magnet) Match(msg string) bool { return strings.HasPrefix(msg, "magnet:?") }

func (m *Magnet) Run(ctx context.Context, b *bot.Bot, msg *bot.Message) {
	chat := msg.Chat.ID
	t, err := b.NewTransfer(ctx, msg.Text)
	if err != nil {
		txt := fmt.Sprintf("error creating new transfer: %v", err)
		b.SendMessage(chat, txt)
		return
	}
	md := telegram.WithParseMode(telegram.ModeMarkdown)
	b.SendMessage(chat, fmt.Sprintf("*%v* is added to transfer list", t.Name), md)
}
