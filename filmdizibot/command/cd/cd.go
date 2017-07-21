package cd

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/igungor/cmd/filmdizibot/bot"
	"github.com/igungor/go-putio/putio"
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

func chooseFolder(parent int64, children []putio.File, args ...string) int64 {
	const (
		rootFolder        int64 = 0
		nonexistingFolder int64 = -1
	)

	// go to root
	if len(args) == 0 {
		return rootFolder
	}

	// go to parent dir
	dirname := strings.Join(args, " ")
	if dirname == ".." {
		return parent
	}

	// lookup for foldername
	for _, f := range children {
		if !f.IsDir() {
			continue
		}

		if f.Name == dirname {
			return f.ID
		}
	}

	// lookup for indices that we generate
	idx, err := strconv.ParseInt(dirname, 10, 64)
	if err != nil || idx <= 0 {
		return nonexistingFolder
	}

	if len(children) < int(idx) {
		return nonexistingFolder
	}

	// idx starts from 1
	return children[idx-1].ID
}

func (c *CD) Run(ctx context.Context, b *bot.Bot, msg *bot.Message) {
	chat := msg.Chat.ID
	args := msg.Args()

	folderId := chooseFolder(b.CWD().ParentID, b.Children(), args...)
	if folderId == -1 {
		b.SendMessage(chat, fmt.Sprintf("%q klasorunu bulamadim", folderId))
		return
	}

	children, parent, err := b.ListFiles(ctx, folderId)
	if err != nil {
		log.Printf("Error listing files for [%v] %q: %v\n", b.CWD(), folderId, err)
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
