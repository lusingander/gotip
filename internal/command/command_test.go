package command

import "testing"

func TestTestNameToTestRunRegex(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		isPrefix bool
		want     string
	}{
		{
			name:    "plain name",
			pattern: "TestFoo/case",
			want:    "^TestFoo$/^case$",
		},
		{
			name:    "plus",
			pattern: "TestFoo/a+b",
			want:    `^TestFoo$/^a\+b$`,
		},
		{
			name:    "parentheses",
			pattern: "TestFoo/case(1)",
			want:    `^TestFoo$/^case\(1\)$`,
		},
		{
			name:    "dot",
			pattern: "TestFoo/foo.bar",
			want:    `^TestFoo$/^foo\.bar$`,
		},
		{
			name:    "brackets",
			pattern: "TestFoo/[]int",
			want:    `^TestFoo$/^\[\]int$`,
		},
		{
			name:    "metacharacters in each segment",
			pattern: "Test.Foo/foo|bar/baz?",
			want:    `^Test\.Foo$/^foo\|bar$/^baz\?$`,
		},
		{
			name:     "prefix",
			pattern:  "Test.Foo/a+b",
			isPrefix: true,
			want:     `^Test\.Foo$/^a\+b`,
		},
		{
			name:     "prefix with trailing slash",
			pattern:  "Test.Foo/a+b/",
			isPrefix: true,
			want:     `^Test\.Foo$/^a\+b$/`,
		},
		{
			name: "empty pattern",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := testNameToTestRunRegex(tt.pattern, tt.isPrefix)
			if got != tt.want {
				t.Errorf("testNameToTestRunRegex(%q, %t) = %q, want %q", tt.pattern, tt.isPrefix, got, tt.want)
			}
		})
	}
}
