package ls

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/igungor/cmd/filmdizibot/bot"
	"github.com/igungor/telegram"
)

type List struct{}

func New() *List { return &List{} }

func (l *List) Name() string { return "ls" }

func (l *List) Usage() string { return "foobar" }

func (l *List) Match(msg string) bool { return strings.HasPrefix(msg, "ls") }

func (l *List) Run(ctx context.Context, b *bot.Bot, msg *bot.Message) {
	chat := msg.Chat.ID

	cwd := b.CWD()
	files, _, err := b.ListFiles(ctx, cwd.ID)
	if err != nil {
		txt := fmt.Sprintf("error listing %q: %v", cwd.Name, err)
		b.SendMessage(chat, txt)
		return
	}
	var buf bytes.Buffer
	for i, f := range files {
		fname := f.Name
		if f.IsDir() {
			fname = fname + "/"
		}
		fname = strings.Replace(fname, "[", "［", -1)
		fname = strings.Replace(fname, "]", "］", -1)
		buf.WriteString(fmt.Sprintf("%02d  [%v](https://put.io/files/%v)\n", i+1, fname, f.ID))
	}
	md := telegram.WithParseMode(telegram.ModeMarkdown)
	b.SendMessage(chat, buf.String(), md)
}
