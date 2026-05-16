package listfmt

import (
	"testing"

	"github.com/lusingander/gotip/internal/tip"
)

func TestFormatText(t *testing.T) {
	tests := fixtureTests()
	got, err := FormatText(tests)
	if err != nil {
		t.Fatalf("FormatText() error = %v", err)
	}
	want := `# ./a/a_test.go
- TestA
  - alpha
  - ??? [unresolved]

# ./b/b_test.go
- TestB
  - outer
    - inner
`
	if got != want {
		t.Errorf("FormatText() = %q, want %q", got, want)
	}
}

func TestFormatJSON(t *testing.T) {
	tests := fixtureTests()
	got, err := FormatJSON(tests)
	if err != nil {
		t.Fatalf("FormatJSON() error = %v", err)
	}
	want := `{
  "files": [
    {
      "path": "./a/a_test.go",
      "tests": [
        {
          "name": "TestA",
          "subtests": [
            {
              "name": "alpha",
              "resolved": true,
              "subtests": []
            },
            {
              "name": null,
              "resolved": false,
              "subtests": []
            }
          ]
        }
      ]
    },
    {
      "path": "./b/b_test.go",
      "tests": [
        {
          "name": "TestB",
          "subtests": [
            {
              "name": "outer",
              "resolved": true,
              "subtests": [
                {
                  "name": "inner",
                  "resolved": true,
                  "subtests": []
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
`
	if got != want {
		t.Errorf("FormatJSON() = %q, want %q", got, want)
	}
}

func fixtureTests() map[string][]*tip.TestFunction {
	return map[string][]*tip.TestFunction{
		"./b/b_test.go": {
			{
				Name: "TestB",
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
		},
		"./a/a_test.go": {
			{
				Name: "TestA",
				Subs: []*tip.SubTest{
					{Name: "alpha", Resolved: true, Subs: []*tip.SubTest{}},
					{Name: "", Resolved: false, Subs: []*tip.SubTest{}},
				},
			},
		},
	}
}
