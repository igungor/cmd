package pwd

import (
	"context"
	"fmt"
	"strings"

	"github.com/igungor/cmd/filmdizibot/bot"
)

type PWD struct{}

func New() *PWD { return &PWD{} }

func (p *PWD) Name() string { return "pwd" }

func (p *PWD) Usage() string { return "Print working directory" }

func (p *PWD) Match(msg string) bool { return strings.HasPrefix(msg, "pwd") }

func (p *PWD) Run(ctx context.Context, b *bot.Bot, msg *bot.Message) {
	chat := msg.Chat.ID
	cwd := b.CWD()
	b.SendMessage(chat, fmt.Sprintf("[%v] %v", cwd.ID, cwd.Name))
}
