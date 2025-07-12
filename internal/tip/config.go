package tip

import (
	"os"

	"github.com/BurntSushi/toml"
)

const (
	configFileName = "gotip.toml"
)

type Config struct {
	Command []string `toml:"command"`
}

func defaultConfig() *Config {
	return &Config{
		Command: []string{},
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
