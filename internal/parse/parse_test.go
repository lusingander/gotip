package parse

import "testing"

func TestProcessFile(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     []*TestFunction
	}{
		{"a", "testdata/foo/a_test.go", wantTestA()},
		{"b", "testdata/foo/b_test.go", wantTestB()},
		{"c", "testdata/bar/c_test.go", wantTestC()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProcessFile(tt.filePath)
			if err != nil {
				t.Errorf("ProcessFile(%s) error = %v", tt.filePath, err)
				return
			}
			assertEqualTests(t, got, tt.want)
		})
	}
}

func wantTestA() []*TestFunction {
	return []*TestFunction{
		{
			name: "TestA1",
			subs: []*SubTest{},
		},
		{
			name: "TestA2_1",
			subs: []*SubTest{
				{name: "test1", subs: []*SubTest{}},
				{name: "test2", subs: []*SubTest{}},
				{name: "test3", subs: []*SubTest{}},
			},
		},
		{
			name: "TestA2_2",
			subs: []*SubTest{
				{name: "test1", subs: []*SubTest{}},
				{name: "test2", subs: []*SubTest{}},
				{name: "test3", subs: []*SubTest{}},
			},
		},
		{
			name: "TestA2_3",
			subs: []*SubTest{
				{name: "test1", subs: []*SubTest{}},
				{name: "test2", subs: []*SubTest{}},
				{name: "test3", subs: []*SubTest{}},
			},
		},
		{
			name: "TestA3",
			subs: []*SubTest{
				{name: "<unknown>", subs: []*SubTest{}},
			},
		},
		{
			name: "TestA4",
			subs: []*SubTest{
				{name: "<unknown>", subs: []*SubTest{}},
			},
		},
		{
			name: "TestA5",
			subs: []*SubTest{
				{name: "<unknown>", subs: []*SubTest{}},
			},
		},
	}
}

func wantTestB() []*TestFunction {
	return []*TestFunction{
		{
			name: "TestB1",
			subs: []*SubTest{
				{name: "test1", subs: []*SubTest{}},
				{name: "test2", subs: []*SubTest{}},
			},
		},
	}
}

func wantTestC() []*TestFunction {
	return []*TestFunction{
		{
			name: "TestC1",
			subs: []*SubTest{
				{
					name: "test1",
					subs: []*SubTest{
						{
							name: "subtest1",
							subs: []*SubTest{},
						},
						{
							name: "subtest2",
							subs: []*SubTest{},
						},
						{
							name: "subtest3",
							subs: []*SubTest{
								{name: "subsubtest1", subs: []*SubTest{}},
								{name: "subsubtest2", subs: []*SubTest{}},
							},
						},
					},
				},
			},
		},
	}
}

func assertEqualTests(t *testing.T, got, want []*TestFunction) {
	if len(got) != len(want) {
		t.Errorf("got tests length = %d, want %d", len(got), len(want))
		return
	}
	for i := range got {
		assertEqualTest(t, got[i], want[i])
	}
}

func assertEqualTest(t *testing.T, got, want *TestFunction) {
	if got.name != want.name {
		t.Errorf("got name = %s, want %s", got.name, want.name)
		return
	}
	assertEqualSubTests(t, got.subs, want.subs)
}

func assertEqualSubTests(t *testing.T, got, want []*SubTest) {
	if len(got) != len(want) {
		t.Errorf("got subs length = %d, want %d", len(got), len(want))
		return
	}
	for i := range got {
		assertEqualSubTest(t, got[i], want[i])
	}
}

func assertEqualSubTest(t *testing.T, got, want *SubTest) {
	if got.name != want.name {
		t.Errorf("got name = %s, want %s", got.name, want.name)
		return
	}
	if len(got.subs) != len(want.subs) {
		t.Errorf("got subs length = %d, want %d", len(got.subs), len(want.subs))
		return
	}
	for i := range got.subs {
		assertEqualSubTest(t, got.subs[i], want.subs[i])
	}
}
