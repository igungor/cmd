package cd

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/igungor/cmd/filmdizibot/bot"
)

type CD struct{}

func New() *CD { return &CD{} }

func (c *CD) Name() string { return "cd" }

func (c *CD) Usage() string {
	var buf bytes.Buffer
	buf.WriteString("Change the working directory\n")
	buf.WriteString("cd: go to /\n")
	buf.WriteString("cd <idx>: go to the directory with the index idx\n")
	buf.WriteString("cd <folder>: go to the directory with the name 'folder'\n")
	return buf.String()
}

func (c *CD) Match(msg string) bool {
	return strings.HasPrefix(msg, "cd")
}

func (c *CD) Run(ctx context.Context, b *bot.Bot, msg *bot.Message) {
	chat := msg.Chat.ID
	args := msg.Args()
	dirname := strings.Join(args, " ")
	children := b.Children()
	var cwd int64
	if len(args) == 0 { // go to home
		cwd = 0
	} else if dirname == ".." { // parent
		cwd = b.CWD().ParentID
	} else if len(dirname) == 2 { // indice based navigation
		idx, err := strconv.ParseInt(dirname, 10, 64)
		if err != nil || idx <= 0 {
			cwd = -1
		} else {
			if len(children) < int(idx) {
				cwd = -1
			} else {
				cwd = children[idx-1].ID
			}
		}
	} else {
		cwd = -1
		for _, f := range children {
			if !f.IsDir() {
				continue
			}
			if f.Name == dirname {
				cwd = f.ID
				break
			}
		}
	}
	if cwd == -1 {
		b.SendMessage(chat, fmt.Sprintf("%q klasorunu bulamadim", dirname))
		return
	}

	children, parent, err := b.ListFiles(ctx, cwd)
	if err != nil {
		log.Printf("Error listing files for [%v] %q: %v\n", cwd, dirname, err)
		b.SendMessage(chat, fmt.Sprintf("error listing files: %v", err))
		return
	}

	if !parent.IsDir() {
		b.SendMessage(chat, fmt.Sprintf("%q bir klasor degil", parent.Name))
		return
	}
	b.SetCWD(parent)
	b.SetChildren(children)
}
