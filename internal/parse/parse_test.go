package parse

import (
	"testing"

	"github.com/lusingander/gotip/internal/tip"
)

func TestProcessFile(t *testing.T) {
	skipSubtests := false
	tests := []struct {
		name     string
		filePath string
		want     []*tip.TestFunction
	}{
		{"a", "testdata/foo/a_test.go", wantTestA()},
		{"b", "testdata/foo/b_test.go", wantTestB()},
		{"c", "testdata/bar/c_test.go", wantTestC()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processFile(tt.filePath, skipSubtests)
			if err != nil {
				t.Errorf("ProcessFile(%s) error = %v", tt.filePath, err)
				return
			}
			assertEqualTests(t, got, tt.want)
		})
	}
}

func wantTestA() []*tip.TestFunction {
	return []*tip.TestFunction{
		{
			Name: "TestA1",
			Subs: []*tip.SubTest{},
		},
		{
			Name: "TestA2_1",
			Subs: []*tip.SubTest{
				{Name: "test1", Subs: []*tip.SubTest{}},
				{Name: "test2", Subs: []*tip.SubTest{}},
				{Name: "test3", Subs: []*tip.SubTest{}},
			},
		},
		{
			Name: "TestA2_2",
			Subs: []*tip.SubTest{
				{Name: "test1", Subs: []*tip.SubTest{}},
				{Name: "test2", Subs: []*tip.SubTest{}},
				{Name: "test3", Subs: []*tip.SubTest{}},
			},
		},
		{
			Name: "TestA2_3",
			Subs: []*tip.SubTest{
				{Name: "test1", Subs: []*tip.SubTest{}},
				{Name: "test2", Subs: []*tip.SubTest{}},
				{Name: "test3", Subs: []*tip.SubTest{}},
			},
		},
		{
			Name: "TestA2_4",
			Subs: []*tip.SubTest{
				{Name: "test1", Subs: []*tip.SubTest{}},
				{Name: "test2", Subs: []*tip.SubTest{}},
			},
		},
		{
			Name: "TestA2_5",
			Subs: []*tip.SubTest{
				{Name: "???", Subs: []*tip.SubTest{}},
			},
		},
		{
			Name: "TestA3",
			Subs: []*tip.SubTest{
				{Name: "???", Subs: []*tip.SubTest{}},
			},
		},
		{
			Name: "TestA4",
			Subs: []*tip.SubTest{
				{Name: "???", Subs: []*tip.SubTest{}},
			},
		},
		{
			Name: "TestA5",
			Subs: []*tip.SubTest{
				{Name: "???", Subs: []*tip.SubTest{}},
			},
		},
	}
}

func wantTestB() []*tip.TestFunction {
	return []*tip.TestFunction{
		{
			Name: "TestB1",
			Subs: []*tip.SubTest{
				{Name: "test1", Subs: []*tip.SubTest{}},
				{Name: "test2", Subs: []*tip.SubTest{}},
			},
		},
	}
}

func wantTestC() []*tip.TestFunction {
	return []*tip.TestFunction{
		{
			Name: "TestC1",
			Subs: []*tip.SubTest{
				{
					Name: "test1",
					Subs: []*tip.SubTest{
						{
							Name: "subtest1",
							Subs: []*tip.SubTest{},
						},
						{
							Name: "subtest2",
							Subs: []*tip.SubTest{},
						},
						{
							Name: "subtest3",
							Subs: []*tip.SubTest{
								{Name: "subsubtest1", Subs: []*tip.SubTest{}},
								{Name: "subsubtest2", Subs: []*tip.SubTest{}},
							},
						},
					},
				},
			},
		},
	}
}

func TestProcessFile_skipSubtests(t *testing.T) {
	skipSubtests := true
	tests := []struct {
		name     string
		filePath string
		want     []*tip.TestFunction
	}{
		{"a", "testdata/foo/a_test.go", wantSkipSubtestsTestA()},
		{"b", "testdata/foo/b_test.go", wantSkipSubtestsTestB()},
		{"c", "testdata/bar/c_test.go", wantSkipSubtestsTestC()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processFile(tt.filePath, skipSubtests)
			if err != nil {
				t.Errorf("ProcessFile(%s) error = %v", tt.filePath, err)
				return
			}
			assertEqualTests(t, got, tt.want)
		})
	}
}

func wantSkipSubtestsTestA() []*tip.TestFunction {
	return []*tip.TestFunction{
		{Name: "TestA1", Subs: []*tip.SubTest{}},
		{Name: "TestA2_1", Subs: []*tip.SubTest{}},
		{Name: "TestA2_2", Subs: []*tip.SubTest{}},
		{Name: "TestA2_3", Subs: []*tip.SubTest{}},
		{Name: "TestA2_4", Subs: []*tip.SubTest{}},
		{Name: "TestA2_5", Subs: []*tip.SubTest{}},
		{Name: "TestA3", Subs: []*tip.SubTest{}},
		{Name: "TestA4", Subs: []*tip.SubTest{}},
		{Name: "TestA5", Subs: []*tip.SubTest{}},
	}
}

func wantSkipSubtestsTestB() []*tip.TestFunction {
	return []*tip.TestFunction{{Name: "TestB1", Subs: []*tip.SubTest{}}}
}

func wantSkipSubtestsTestC() []*tip.TestFunction {
	return []*tip.TestFunction{{Name: "TestC1", Subs: []*tip.SubTest{}}}
}

func assertEqualTests(t *testing.T, got, want []*tip.TestFunction) {
	if len(got) != len(want) {
		t.Errorf("got tests length = %d, want %d", len(got), len(want))
		return
	}
	for i := range got {
		assertEqualTest(t, got[i], want[i])
	}
}

func assertEqualTest(t *testing.T, got, want *tip.TestFunction) {
	if got.Name != want.Name {
		t.Errorf("got name = %s, want %s", got.Name, want.Name)
		return
	}
	assertEqualSubTests(t, got.Subs, want.Subs)
}

func assertEqualSubTests(t *testing.T, got, want []*tip.SubTest) {
	if len(got) != len(want) {
		t.Errorf("got subs length = %d, want %d", len(got), len(want))
		return
	}
	for i := range got {
		assertEqualSubTest(t, got[i], want[i])
	}
}

func assertEqualSubTest(t *testing.T, got, want *tip.SubTest) {
	if got.Name != want.Name {
		t.Errorf("got name = %s, want %s", got.Name, want.Name)
		return
	}
	if len(got.Subs) != len(want.Subs) {
		t.Errorf("got subs length = %d, want %d", len(got.Subs), len(want.Subs))
		return
	}
	for i := range got.Subs {
		assertEqualSubTest(t, got.Subs[i], want.Subs[i])
	}
}
