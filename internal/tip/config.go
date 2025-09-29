package tip

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const (
	configFileName = "gotip.toml"

	defaultHistoryLimit = 100
	defaultDateFormat   = "2006-01-02 15:04:05"
)

type Config struct {
	Command []string      `toml:"command"`
	Ignore  []string      `toml:"ignore"`
	History HistoryConfig `toml:"history"`
}

type HistoryConfig struct {
	Limit      int    `toml:"limit"`
	DateFormat string `toml:"date_format"`
}

func defaultConfig() *Config {
	return &Config{
		Command: []string{},
		Ignore:  []string{},
		History: HistoryConfig{
			Limit:      defaultHistoryLimit,
			DateFormat: defaultDateFormat,
		},
	}
}

func LoadConfig(projectDir string) (*Config, error) {
	conf := defaultConfig()

	globalConfigPath, err := globalConfigFilePath()
	if err != nil {
		return nil, err
	}
	if conf, err = loadAndMergeConfig(globalConfigPath, conf); err != nil {
		return nil, err
	}

	projectConfigPath, err := projectConfigFilePath(projectDir)
	if err != nil {
		return nil, err
	}
	if conf, err = loadAndMergeConfig(projectConfigPath, conf); err != nil {
		return nil, err
	}

	return conf, nil
}

func loadAndMergeConfig(filePath string, base *Config) (*Config, error) {
	if _, err := os.Stat(filePath); err != nil {
		return base, nil
	}
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	if _, err = toml.Decode(string(bytes), base); err != nil {
		return nil, err
	}
	return base, nil
}

func projectConfigFilePath(projectDir string) (string, error) {
	absDir, err := filepath.Abs(projectDir)
	if err != nil {
		return "", err
	}
	return filepath.Join(absDir, configFileName), nil
}

func globalConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "gotip", configFileName), nil
}
