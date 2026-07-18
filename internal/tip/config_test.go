package tip

import (
	"os"
	"path/filepath"
	"reflect"
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

func TestLoadAndMergeConfigMergesColorTheme(t *testing.T) {
	dir := t.TempDir()
	globalPath := filepath.Join(dir, "global.toml")
	projectPath := filepath.Join(dir, "project.toml")
	writeConfigFile(t, globalPath, "[theme]\naccent = \"#112233\"\nmuted = \"245\"\n")
	writeConfigFile(t, projectPath, "[theme]\naccent = \"#abcdef\"\n")

	conf, err := loadAndMergeConfig(globalPath, defaultConfig())
	if err != nil {
		t.Fatalf("loadAndMergeConfig(global) error = %v", err)
	}
	conf, err = loadAndMergeConfig(projectPath, conf)
	if err != nil {
		t.Fatalf("loadAndMergeConfig(project) error = %v", err)
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

func writeConfigFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config file: %v", err)
	}
}
