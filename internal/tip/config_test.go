package tip

import (
	"reflect"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/lusingander/gotip/internal/theme"
)

func TestDefaultConfigUsesDefaultColorTheme(t *testing.T) {
	got := defaultConfig().Theme
	want := theme.DefaultColorTheme()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("default theme = %#v, want %#v", got, want)
	}
}

func TestDefaultConfigUsesGotipFuzzyMatcher(t *testing.T) {
	if got := defaultConfig().Filter.FuzzyMatcher; got != "gotip" {
		t.Fatalf("default fuzzy matcher = %q, want %q", got, "gotip")
	}
}

func TestDecodeAndMergeConfigMergesFuzzyMatcher(t *testing.T) {
	globalConfig := strings.NewReader("[filter]\nfuzzy_matcher = \"legacy\"\n")
	projectConfig := strings.NewReader("[filter]\nfuzzy_matcher = \"gotip\"\n")

	conf, err := decodeAndMergeConfig(globalConfig, defaultConfig())
	if err != nil {
		t.Fatalf("decodeAndMergeConfig(global) error = %v", err)
	}
	if conf.Filter.FuzzyMatcher != "legacy" {
		t.Errorf("global fuzzy matcher = %q, want %q", conf.Filter.FuzzyMatcher, "legacy")
	}

	conf, err = decodeAndMergeConfig(projectConfig, conf)
	if err != nil {
		t.Fatalf("decodeAndMergeConfig(project) error = %v", err)
	}
	if conf.Filter.FuzzyMatcher != "gotip" {
		t.Errorf("project fuzzy matcher = %q, want %q", conf.Filter.FuzzyMatcher, "gotip")
	}
}

func TestConfigValidateRejectsUnknownFuzzyMatcher(t *testing.T) {
	conf := defaultConfig()
	conf.Filter.FuzzyMatcher = "unknown"

	if err := conf.validate(); err == nil {
		t.Fatal("validate() error = nil, want an error")
	}
}

func TestDecodeAndMergeConfigMergesColorTheme(t *testing.T) {
	globalConfig := strings.NewReader("[theme]\naccent = \"#112233\"\nmuted = \"245\"\n")
	projectConfig := strings.NewReader("[theme]\naccent = \"#abcdef\"\n")

	conf, err := decodeAndMergeConfig(globalConfig, defaultConfig())
	if err != nil {
		t.Fatalf("decodeAndMergeConfig(global) error = %v", err)
	}
	conf, err = decodeAndMergeConfig(projectConfig, conf)
	if err != nil {
		t.Fatalf("decodeAndMergeConfig(project) error = %v", err)
	}

	if conf.Theme.Accent != lipgloss.Color("#abcdef") {
		t.Errorf("theme accent = %q, want %q", conf.Theme.Accent, "#abcdef")
	}
	if conf.Theme.Muted != lipgloss.Color("245") {
		t.Errorf("theme muted = %q, want %q", conf.Theme.Muted, "245")
	}
	if conf.Theme.Text != theme.DefaultColorTheme().Text {
		t.Errorf("theme text = %q, want default %q", conf.Theme.Text, theme.DefaultColorTheme().Text)
	}
}
