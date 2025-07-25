package tip

import (
	"os"

	"github.com/BurntSushi/toml"
)

const (
	configFileName = "gotip.toml"

	defaultHistoryLimit = 100
	defaultDateFormat   = "2006-01-02 15:04:05"
)

type Config struct {
	Command []string      `toml:"command"`
	History HistoryConfig `toml:"history"`
}

type HistoryConfig struct {
	Limit      int    `toml:"limit"`
	DateFormat string `toml:"date_format"`
}

func defaultConfig() *Config {
	return &Config{
		Command: []string{},
		History: HistoryConfig{
			Limit:      defaultHistoryLimit,
			DateFormat: defaultDateFormat,
		},
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
