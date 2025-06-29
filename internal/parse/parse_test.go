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
			got, err := processFile(tt.filePath)
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
			Name: "TestA1",
			Subs: []*SubTest{},
		},
		{
			Name: "TestA2_1",
			Subs: []*SubTest{
				{Name: "test1", Subs: []*SubTest{}},
				{Name: "test2", Subs: []*SubTest{}},
				{Name: "test3", Subs: []*SubTest{}},
			},
		},
		{
			Name: "TestA2_2",
			Subs: []*SubTest{
				{Name: "test1", Subs: []*SubTest{}},
				{Name: "test2", Subs: []*SubTest{}},
				{Name: "test3", Subs: []*SubTest{}},
			},
		},
		{
			Name: "TestA2_3",
			Subs: []*SubTest{
				{Name: "test1", Subs: []*SubTest{}},
				{Name: "test2", Subs: []*SubTest{}},
				{Name: "test3", Subs: []*SubTest{}},
			},
		},
		{
			Name: "TestA3",
			Subs: []*SubTest{
				{Name: "<unknown>", Subs: []*SubTest{}},
			},
		},
		{
			Name: "TestA4",
			Subs: []*SubTest{
				{Name: "<unknown>", Subs: []*SubTest{}},
			},
		},
		{
			Name: "TestA5",
			Subs: []*SubTest{
				{Name: "<unknown>", Subs: []*SubTest{}},
			},
		},
	}
}

func wantTestB() []*TestFunction {
	return []*TestFunction{
		{
			Name: "TestB1",
			Subs: []*SubTest{
				{Name: "test1", Subs: []*SubTest{}},
				{Name: "test2", Subs: []*SubTest{}},
			},
		},
	}
}

func wantTestC() []*TestFunction {
	return []*TestFunction{
		{
			Name: "TestC1",
			Subs: []*SubTest{
				{
					Name: "test1",
					Subs: []*SubTest{
						{
							Name: "subtest1",
							Subs: []*SubTest{},
						},
						{
							Name: "subtest2",
							Subs: []*SubTest{},
						},
						{
							Name: "subtest3",
							Subs: []*SubTest{
								{Name: "subsubtest1", Subs: []*SubTest{}},
								{Name: "subsubtest2", Subs: []*SubTest{}},
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
	if got.Name != want.Name {
		t.Errorf("got name = %s, want %s", got.Name, want.Name)
		return
	}
	assertEqualSubTests(t, got.Subs, want.Subs)
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
