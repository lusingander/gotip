package foo

import "testing"

func TestA1(t *testing.T) {
	a := 1
	b := 2
	got := a + b
	want := 3
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestA2(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{"test1", 1, 2, 3},
		{"test2", 2, 3, 5},
		{"test3", 3, 4, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a + tt.b
			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}
