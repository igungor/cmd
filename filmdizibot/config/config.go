package config

import (
	"fmt"

	"github.com/burntsushi/toml"
)

var Config config

type config struct {
	Host  string
	Port  int
	Putio struct {
		Token string
	}
	Telegram struct {
		Webhook string
		Token   string
	}
	Addic7ed struct {
		Username string
		Password string
	}
}

func Load(path string) error {
	_, err := toml.DecodeFile(path, &Config)
	if Config.Putio.Token == "" {
		return fmt.Errorf("putio.token must be specified")
	}
	return err
}
