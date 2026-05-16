package parse

import (
	"testing"

	"github.com/lusingander/gotip/internal/tip"
)

func TestUnresolvedSubTestResolve(t *testing.T) {
	tests := []struct {
		name string
		sub  *unresolvedSubTest
		want []*tip.SubTest
	}{
		{
			name: "literal name",
			sub: &unresolvedSubTest{
				name: &literalSubTestName{name: "outer"},
			},
			want: []*tip.SubTest{{Name: "outer", Resolved: true, Subs: []*tip.SubTest{}}},
		},
		{
			name: "unresolved name",
			sub: &unresolvedSubTest{
				name: &unknownSubTestName{},
			},
			want: []*tip.SubTest{{Name: "", Resolved: false, Subs: []*tip.SubTest{}}},
		},
		{
			name: "expanded names share resolved children",
			sub: &unresolvedSubTest{
				name: &selectorSubTestName{cases: []string{"first", "second"}},
				subs: []*unresolvedSubTest{
					{name: &literalSubTestName{name: "child"}},
				},
			},
			want: []*tip.SubTest{
				{
					Name:     "first",
					Resolved: true,
					Subs: []*tip.SubTest{
						{Name: "child", Resolved: true, Subs: []*tip.SubTest{}},
					},
				},
				{
					Name:     "second",
					Resolved: true,
					Subs: []*tip.SubTest{
						{Name: "child", Resolved: true, Subs: []*tip.SubTest{}},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertEqualSubTests(t, tt.sub.resolve(), tt.want)
		})
	}
}

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
		{"d", "testdata/baz/d_test.go", wantTestD()},
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
			Name: "TestSimpleAddition",
			Subs: []*tip.SubTest{},
		},
		{
			Name: "TestStructSlicePositionalFields",
			Subs: []*tip.SubTest{
				{Name: "test1", Resolved: true, Subs: []*tip.SubTest{}},
				{Name: "test2", Resolved: true, Subs: []*tip.SubTest{}},
				{Name: "test3", Resolved: true, Subs: []*tip.SubTest{}},
			},
		},
		{
			Name: "TestStructSliceKeyedFields",
			Subs: []*tip.SubTest{
				{Name: "test1", Resolved: true, Subs: []*tip.SubTest{}},
				{Name: "test2", Resolved: true, Subs: []*tip.SubTest{}},
				{Name: "test3", Resolved: true, Subs: []*tip.SubTest{}},
			},
		},
		{
			Name: "TestNamedStructSlice",
			Subs: []*tip.SubTest{
				{Name: "test1", Resolved: true, Subs: []*tip.SubTest{}},
				{Name: "test2", Resolved: true, Subs: []*tip.SubTest{}},
				{Name: "test3", Resolved: true, Subs: []*tip.SubTest{}},
			},
		},
		{
			Name: "TestVarStructSlice",
			Subs: []*tip.SubTest{
				{Name: "test1", Resolved: true, Subs: []*tip.SubTest{}},
				{Name: "test2", Resolved: true, Subs: []*tip.SubTest{}},
			},
		},
		{
			Name: "TestNestedStructFieldName",
			Subs: []*tip.SubTest{
				{Name: "", Resolved: false, Subs: []*tip.SubTest{}},
			},
		},
		{
			Name: "TestMapKeySubtests",
			Subs: []*tip.SubTest{
				{Name: "test1", Resolved: true, Subs: []*tip.SubTest{}},
				{Name: "test2", Resolved: true, Subs: []*tip.SubTest{}},
			},
		},
		{
			Name: "TestConcatGeneratedName",
			Subs: []*tip.SubTest{
				{Name: "", Resolved: false, Subs: []*tip.SubTest{}},
			},
		},
		{
			Name: "TestSprintfGeneratedName",
			Subs: []*tip.SubTest{
				{Name: "", Resolved: false, Subs: []*tip.SubTest{}},
			},
		},
		{
			Name: "TestStringIdentNames",
			Subs: []*tip.SubTest{
				{Name: "test1", Resolved: true, Subs: []*tip.SubTest{}},
				{Name: "test2", Resolved: true, Subs: []*tip.SubTest{}},
				{Name: "test3", Resolved: true, Subs: []*tip.SubTest{}},
			},
		},
		{
			Name: "TestNonTestingRunIgnored",
			Subs: []*tip.SubTest{},
		},
		{
			Name: "TestNestedSubtestsWithRenamedTestingReceiver",
			Subs: []*tip.SubTest{
				{
					Name:     "outer",
					Resolved: true,
					Subs: []*tip.SubTest{
						{Name: "inner", Resolved: true, Subs: []*tip.SubTest{}},
					},
				},
			},
		},
	}
}

func wantTestB() []*tip.TestFunction {
	return []*tip.TestFunction{
		{
			Name: "TestLiteralSubtestsWithHelper",
			Subs: []*tip.SubTest{
				{Name: "test1", Resolved: true, Subs: []*tip.SubTest{}},
				{Name: "test2", Resolved: true, Subs: []*tip.SubTest{}},
			},
		},
	}
}

func wantTestC() []*tip.TestFunction {
	return []*tip.TestFunction{
		{
			Name: "TestNestedLiteralSubtests",
			Subs: []*tip.SubTest{
				{
					Name:     "test1",
					Resolved: true,
					Subs: []*tip.SubTest{
						{
							Name:     "subtest1",
							Resolved: true,
							Subs:     []*tip.SubTest{},
						},
						{
							Name:     "subtest2",
							Resolved: true,
							Subs:     []*tip.SubTest{},
						},
						{
							Name:     "subtest3",
							Resolved: true,
							Subs: []*tip.SubTest{
								{Name: "subsubtest1", Resolved: true, Subs: []*tip.SubTest{}},
								{Name: "subsubtest2", Resolved: true, Subs: []*tip.SubTest{}},
							},
						},
					},
				},
			},
		},
	}
}

func wantTestD() []*tip.TestFunction {
	return []*tip.TestFunction{
		{
			Name: "TestValid",
			Subs: []*tip.SubTest{},
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
		{"d", "testdata/baz/d_test.go", wantSkipSubtestsTestD()},
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
		{Name: "TestSimpleAddition", Subs: []*tip.SubTest{}},
		{Name: "TestStructSlicePositionalFields", Subs: []*tip.SubTest{}},
		{Name: "TestStructSliceKeyedFields", Subs: []*tip.SubTest{}},
		{Name: "TestNamedStructSlice", Subs: []*tip.SubTest{}},
		{Name: "TestVarStructSlice", Subs: []*tip.SubTest{}},
		{Name: "TestNestedStructFieldName", Subs: []*tip.SubTest{}},
		{Name: "TestMapKeySubtests", Subs: []*tip.SubTest{}},
		{Name: "TestConcatGeneratedName", Subs: []*tip.SubTest{}},
		{Name: "TestSprintfGeneratedName", Subs: []*tip.SubTest{}},
		{Name: "TestStringIdentNames", Subs: []*tip.SubTest{}},
		{Name: "TestNonTestingRunIgnored", Subs: []*tip.SubTest{}},
		{Name: "TestNestedSubtestsWithRenamedTestingReceiver", Subs: []*tip.SubTest{}},
	}
}

func wantSkipSubtestsTestB() []*tip.TestFunction {
	return []*tip.TestFunction{{Name: "TestLiteralSubtestsWithHelper", Subs: []*tip.SubTest{}}}
}

func wantSkipSubtestsTestC() []*tip.TestFunction {
	return []*tip.TestFunction{{Name: "TestNestedLiteralSubtests", Subs: []*tip.SubTest{}}}
}

func wantSkipSubtestsTestD() []*tip.TestFunction {
	return []*tip.TestFunction{{Name: "TestValid", Subs: []*tip.SubTest{}}}
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
	if got.Resolved != want.Resolved {
		t.Errorf("got resolved = %t, want %t", got.Resolved, want.Resolved)
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
