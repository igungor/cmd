package bot

import (
	"strings"

	"github.com/igungor/telegram"
)

type Message struct {
	*telegram.Message
}

func (m *Message) Command() string {
	cmd := m.Text
	i := strings.Index(cmd, " ")
	if i >= 0 {
		cmd = cmd[:i]
	}
	return cmd
}

func (m *Message) Args() []string {
	args := strings.TrimSpace(m.Text)
	i := strings.Index(args, " ")
	if i < 0 {
		return nil
	}
	return strings.Fields(args[i+1:])
}
