package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type config struct {
	LogLevel string `yaml:"log-level"`
	Plex     struct {
		Addr  string `yaml:"addr"`
		Token string `yaml:"token"`
	} `yaml:"plex"`

	Lists []string `yaml:"lists"`
}

func decodeConfig(path string) (config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return config{}, err
	}

	var c config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return c, err
	}

	if err := validate(c); err != nil {
		return c, err
	}

	return c, nil
}

func validate(cfg config) error {
	if cfg.Plex.Addr == "" {
		return fmt.Errorf("Plex address must be provided")
	}

	if cfg.Plex.Token == "" {
		return fmt.Errorf("Plex API token must be provided")
	}

	if len(cfg.Lists) == 0 {
		return fmt.Errorf("'lists' field is empty")
	}

	return nil
}
