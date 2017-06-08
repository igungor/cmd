package command

import (
	"context"
	"sync"

	"github.com/igungor/cmd/filmdizibot/bot"
	"github.com/igungor/cmd/filmdizibot/command/cd"
	"github.com/igungor/cmd/filmdizibot/command/ls"
	"github.com/igungor/cmd/filmdizibot/command/magnet"
	"github.com/igungor/cmd/filmdizibot/command/mp4"
	"github.com/igungor/cmd/filmdizibot/command/pwd"
	"github.com/igungor/cmd/filmdizibot/command/sub"
)

type Command interface {
	Name() string
	Usage() string
	Match(txt string) bool
	Run(context.Context, *bot.Bot, *bot.Message)
}

// commands
var ()

var (
	mu       sync.Mutex
	commands = map[string]Command{}
)

func init() {
	register(cd.New())
	register(ls.New())
	register(magnet.New())
	register(pwd.New())
	register(sub.New())
	register(mp4.New())
}

func register(cmd Command) {
	mu.Lock()
	defer mu.Unlock()

	name := cmd.Name()
	if name == "" {
		panic("can not register command with an empty name")
	}

	if _, ok := commands[name]; ok {
		panic("command already registered: " + name)
	}
	commands[name] = cmd
}

func Match(msg string) Command {
	mu.Lock()
	defer mu.Unlock()

	if msg == "" {
		return nil
	}

	for _, cmd := range commands {
		if cmd.Match(msg) {
			return cmd
		}
	}
	return nil
}
