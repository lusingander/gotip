package tip

import "testing"

func TestDropLastSegment(t *testing.T) {
	tests := []struct {
		pattern  string
		expected string
	}{
		{
			pattern:  "TestFoo/Bar/Baz/",
			expected: "TestFoo/Bar/",
		},
		{
			pattern:  "TestFoo/Bar/Baz",
			expected: "TestFoo/Bar/",
		},
		{
			pattern:  "TestFoo/Bar/",
			expected: "TestFoo/",
		},
		{
			pattern:  "TestFoo/",
			expected: "",
		},
		{
			pattern:  "TestFoo",
			expected: "",
		},
		{
			pattern:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			sut := &Target{
				TestNamePattern: tt.pattern,
			}
			sut.DropLastSegment()
			if sut.TestNamePattern != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, sut.TestNamePattern)
			}
		})
	}
}
