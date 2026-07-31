package tip

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/lusingander/gotip/internal/theme"
)

const (
	configFileName = "gotip.toml"

	defaultHistoryLimit = 100
	defaultDateFormat   = "2006-01-02 15:04:05"
	defaultFuzzyMatcher = "gotip"
)

type Config struct {
	Command     []string          `toml:"command"`
	Ignore      []string          `toml:"ignore"`
	Filter      FilterConfig      `toml:"filter"`
	History     HistoryConfig     `toml:"history"`
	Keybindings KeybindingsConfig `toml:"keybindings"`
	Theme       theme.ColorTheme  `toml:"theme"`
}

type FilterConfig struct {
	FuzzyMatcher string `toml:"fuzzy_matcher"`
}

type HistoryConfig struct {
	Limit      int    `toml:"limit"`
	DateFormat string `toml:"date_format"`
}

func defaultConfig() *Config {
	return &Config{
		Command: []string{},
		Ignore:  []string{},
		Filter: FilterConfig{
			FuzzyMatcher: defaultFuzzyMatcher,
		},
		History: HistoryConfig{
			Limit:      defaultHistoryLimit,
			DateFormat: defaultDateFormat,
		},
		Keybindings: DefaultKeybindingsConfig(),
		Theme:       theme.DefaultColorTheme(),
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
	if err := conf.validate(); err != nil {
		return nil, err
	}

	return conf, nil
}

func (c Config) validate() error {
	switch c.Filter.FuzzyMatcher {
	case "gotip", "legacy":
	default:
		return fmt.Errorf(
			"invalid filter.fuzzy_matcher %q: must be \"gotip\" or \"legacy\"",
			c.Filter.FuzzyMatcher,
		)
	}
	if err := c.Keybindings.validate(); err != nil {
		return err
	}
	return c.Theme.Validate()
}

func loadAndMergeConfig(filePath string, base *Config) (*Config, error) {
	file, err := os.Open(filePath)
	if errors.Is(err, os.ErrNotExist) {
		return base, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return decodeAndMergeConfig(file, base)
}

func decodeAndMergeConfig(reader io.Reader, base *Config) (*Config, error) {
	if _, err := toml.NewDecoder(reader).Decode(base); err != nil {
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
