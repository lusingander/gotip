package tip

import (
	"testing"
	"time"
)

func TestHistoriesAdd(t *testing.T) {
	sut := &Histories{
		ProjectDir: "/path/to/project",
		Histories:  []*History{},
	}

	sut.Add(NewTarget("./foo/foo_test.go", "TestA", false), 10)
	sut.Add(NewTarget("./foo/foo_test.go", "TestB", false), 10)
	sut.Add(NewTarget("./bar/bar_test.go", "TestC", false), 10)
	sut.Add(NewTarget("./bar/bar_test.go", "TestD", false), 10)

	assertHistoriesCount(t, sut, 4)
	assertHistoryTestName(t, sut.Histories[0], "TestD")
	assertHistoryTestName(t, sut.Histories[1], "TestC")
	assertHistoryTestName(t, sut.Histories[2], "TestB")
	assertHistoryTestName(t, sut.Histories[3], "TestA")

	sut.Add(NewTarget("./foo/foo_test.go", "TestE", false), 3)

	assertHistoriesCount(t, sut, 3)
	assertHistoryTestName(t, sut.Histories[0], "TestE")
	assertHistoryTestName(t, sut.Histories[1], "TestD")
	assertHistoryTestName(t, sut.Histories[2], "TestC")

	sut.Add(NewTarget("./bar/bar_test.go", "TestF", false), 3)

	assertHistoriesCount(t, sut, 3)
	assertHistoryTestName(t, sut.Histories[0], "TestF")
	assertHistoryTestName(t, sut.Histories[1], "TestE")
	assertHistoryTestName(t, sut.Histories[2], "TestD")

	sut.Add(NewTarget("./foo/foo_test.go", "TestE", false), 3)

	assertHistoriesCount(t, sut, 3)
	assertHistoryTestName(t, sut.Histories[0], "TestE")
	assertHistoryTestName(t, sut.Histories[1], "TestF")
	assertHistoryTestName(t, sut.Histories[2], "TestD")

	sut.Add(NewTarget("./foo/foo_test.go", "TestE", false), 5)

	assertHistoriesCount(t, sut, 3)
	assertHistoryTestName(t, sut.Histories[0], "TestE")
	assertHistoryTestName(t, sut.Histories[1], "TestF")
	assertHistoryTestName(t, sut.Histories[2], "TestD")
}

func assertHistoriesCount(t *testing.T, histories *Histories, wantCount int) {
	if len(histories.Histories) != wantCount {
		t.Errorf("want %d histories, got %d", wantCount, len(histories.Histories))
	}
}

func assertHistoryTestName(t *testing.T, history *History, wantName string) {
	if history.TestNamePattern != wantName {
		t.Errorf("want history TestNamePattern to be %s, got %s", wantName, history.TestNamePattern)
	}
}

func TestHistoryReferToSameHistory(t *testing.T) {
	history := &History{
		Path:            "./foo/foo_test.go",
		PackageName:     "foo",
		TestNamePattern: "TestA",
		IsPrefix:        false,
		RunAt:           time.Date(2025, 7, 19, 18, 0, 0, 0, time.UTC),
	}

	tests := []struct {
		name  string
		other *History
		want  bool
	}{
		{
			name: "same history",
			other: &History{
				Path:            "./foo/foo_test.go",
				PackageName:     "foo",
				TestNamePattern: "TestA",
				IsPrefix:        false,
				RunAt:           time.Date(2025, 7, 20, 12, 0, 0, 0, time.UTC), // RunAt can differ
			},
			want: true,
		},
		{
			name: "different path",
			other: &History{
				Path:            "./bar/bar_test.go",
				PackageName:     "foo",
				TestNamePattern: "TestA",
				IsPrefix:        false,
				RunAt:           time.Date(2025, 7, 20, 12, 0, 0, 0, time.UTC),
			},
			want: false,
		},
		{
			name: "different package name",
			other: &History{
				Path:            "./foo/foo_test.go",
				PackageName:     "bar",
				TestNamePattern: "TestA",
				IsPrefix:        false,
				RunAt:           time.Date(2025, 7, 20, 12, 0, 0, 0, time.UTC),
			},
			want: false,
		},
		{
			name: "different test name pattern",
			other: &History{
				Path:            "./foo/foo_test.go",
				PackageName:     "foo",
				TestNamePattern: "TestB",
				IsPrefix:        false,
				RunAt:           time.Date(2025, 7, 20, 12, 0, 0, 0, time.UTC),
			},
			want: false,
		},
		{
			name: "different isPrefix",
			other: &History{
				Path:            "./foo/foo_test.go",
				PackageName:     "foo",
				TestNamePattern: "TestA",
				IsPrefix:        true,
				RunAt:           time.Date(2025, 7, 20, 12, 0, 0, 0, time.UTC),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := history.referToSameHistory(tt.other); got != tt.want {
				t.Errorf("want %v, got %v", tt.want, got)
			}
		})
	}
}
