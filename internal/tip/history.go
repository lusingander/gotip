package tip

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"time"
)

type Histories struct {
	ProjectDir string
	Histories  []*History
}

func newHistories(projectDir string) (*Histories, error) {
	absDir, err := filepath.Abs(projectDir)
	if err != nil {
		return nil, err
	}
	return &Histories{
		ProjectDir: absDir,
		Histories:  []*History{},
	}, nil
}

func (h *Histories) Add(target *Target, limit int) {
	history := &History{
		Path:            target.Path,
		PackageName:     target.PackageName,
		TestNamePattern: target.TestNamePattern,
		IsPrefix:        target.IsPrefix,
		RunAt:           time.Now(),
	}

	// Remove existing history if it refers to the same test to avoid duplicates
	if i := slices.IndexFunc(h.Histories, history.referToSameHistory); i >= 0 {
		h.Histories = append(h.Histories[:i], h.Histories[i+1:]...)
	}

	h.Histories = append([]*History{history}, h.Histories...)
	if limit >= 0 && len(h.Histories) > limit {
		h.Histories = h.Histories[:limit]
	}
}

type History struct {
	Path            string
	PackageName     string
	TestNamePattern string
	IsPrefix        bool
	RunAt           time.Time
}

func (h *History) referToSameHistory(other *History) bool {
	return h.Path == other.Path &&
		h.PackageName == other.PackageName &&
		h.TestNamePattern == other.TestNamePattern &&
		h.IsPrefix == other.IsPrefix
}

func (h *History) ToTarget() *Target {
	return &Target{
		Path:            h.Path,
		PackageName:     h.PackageName,
		TestNamePattern: h.TestNamePattern,
		IsPrefix:        h.IsPrefix,
	}
}

func LoadHistories(projectDir string) (*Histories, error) {
	filePath, err := projectHistoriesFilePath(projectDir)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(filePath); err != nil {
		return newHistories(projectDir)
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var histories Histories
	if err = json.Unmarshal(bytes, &histories); err != nil {
		return nil, err
	}
	return &histories, nil
}

func SaveHistories(projectDir string, histories *Histories) error {
	filePath, err := projectHistoriesFilePath(projectDir)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(filepath.Dir(filePath), 0o700); err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(histories, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, bytes, 0o600)
}

func projectHistoriesFilePath(projectDir string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fileName, err := projectHistoriesFileName(projectDir)
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "state", "gotip", "history", fileName), nil
}

func projectHistoriesFileName(projectDir string) (string, error) {
	absDir, err := filepath.Abs(projectDir)
	if err != nil {
		return "", err
	}
	dir := filepath.ToSlash(absDir)
	hash := md5.Sum([]byte(dir))
	return hex.EncodeToString(hash[:]) + ".json", nil
}
