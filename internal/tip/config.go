package tip

import (
	"os"

	"github.com/BurntSushi/toml"
)

const (
	configFileName = "gotip.toml"

	defaultHistoryLimit = 100
)

type Config struct {
	Command      []string `toml:"command"`
	HistoryLimit int      `toml:"history_limit"`
}

func defaultConfig() *Config {
	return &Config{
		Command:      []string{},
		HistoryLimit: defaultHistoryLimit,
	}
}

func LoadConfig() (*Config, error) {
	conf := defaultConfig()

	if _, err := os.Stat(configFileName); err != nil {
		return conf, nil
	}

	bytes, err := os.ReadFile(configFileName)
	if err != nil {
		return nil, err
	}

	if _, err = toml.Decode(string(bytes), &conf); err != nil {
		return nil, err
	}
	return conf, nil
}
