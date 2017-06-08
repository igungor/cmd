package cd

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/igungor/cmd/filmdizibot/bot"
	"github.com/igungor/go-putio/putio"
)

type MP4 struct{}

func New() *MP4 { return &MP4{} }

func (m *MP4) Name() string { return "mp4" }

func (m *MP4) Usage() string {
	var buf bytes.Buffer
	buf.WriteString("Convert video to mp4\n")
	buf.WriteString("mp4: idx0 idx1 idx2...\n")
	return buf.String()
}

func (m *MP4) Match(msg string) bool {
	return strings.HasPrefix(msg, "mp4")
}

func (m *MP4) Run(ctx context.Context, b *bot.Bot, msg *bot.Message) {
	chat := msg.Chat.ID
	args := msg.Args()

	if len(args) == 0 {
		b.SendMessage(chat, "choose something")
		return
	}

	children := b.Children()

	var files []putio.File
	for _, arg := range args {
		idx, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			b.SendMessage(chat, fmt.Sprintf("wrong selection %q: %v", arg, err))
			return
		}

		if len(children) < int(idx) {
			b.SendMessage(chat, fmt.Sprintf("wrong selection: no such idx"))
			return
		}

		//  idx starts from 1
		f := children[idx-1]
		if f.IsDir() {
			b.SendMessage(chat, "wrong selection: choose a file")
			return
		}
		files = append(files, f)
	}

	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, f := range files {
		go func(f putio.File) {
			defer wg.Done()

			err := b.ConvertToMP4(ctx, f.ID)
			if err != nil {
				b.SendMessage(chat, fmt.Sprintf("error converting %v: %v", f.Name, err))
				return
			}
		}(f)
	}

	wg.Wait()

	b.SendMessage(chat, "all convert request are sent")
}
