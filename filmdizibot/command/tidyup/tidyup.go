package tidyup

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/igungor/cmd/filmdizibot/bot"
)

type Tidyup struct{}

func New() *Tidyup { return &Tidyup{} }

func (t *Tidyup) Name() string { return "tidyup" }

func (t *Tidyup) Usage() string { return "foobar" }

func (t *Tidyup) Match(msg string) bool { return strings.HasPrefix(msg, "tidyup") }

func (t *Tidyup) Run(ctx context.Context, b *bot.Bot, msg *bot.Message) {
	chat := msg.Chat.ID
	const root int64 = 0
	files, _, err := b.ListFiles(ctx, root)
	if err != nil {
		b.SendMessage(chat, fmt.Sprintf("tidyup failed: %v", err))
		return
	}
	movieDirID := int64(-1)
	const movieFolder = "movie"
	for _, f := range files {
		if f.Name == movieFolder {
			movieDirID = f.ID
		}
	}

	if movieDirID == -1 {
		b.SendMessage(chat, "could not find 'movie' folder")
		return
	}

	files, _, err = b.ListFiles(ctx, movieDirID)
	if err != nil {
		b.SendMessage(chat, fmt.Sprintf("tidyup failed: %v", err))
		return
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if !strings.HasPrefix(f.ContentType, "video") {
			continue
		}

		dirname := strings.TrimSuffix(f.Name, filepath.Ext(f.Name))
		parent, err := b.Mkdir(ctx, dirname, movieDirID)
		if err != nil {
			b.SendMessage(chat, fmt.Sprintf("tidyup failed: %v", err))
			return
		}

		err = b.MoveFiles(ctx, parent.ID, f.ID)
		if err != nil {
			b.SendMessage(chat, fmt.Sprintf("tidyup failed: %v", err))
			return
		}
	}
	b.SendMessage(chat, "tidyup is done!")
}
