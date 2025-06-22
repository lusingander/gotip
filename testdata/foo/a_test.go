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

func TestA3(t *testing.T) {
	for i := 0; i < 3; i++ {
		t.Run("test"+string(i), func(t *testing.T) {
			a := i + 1
			b := i + 2
			got := a + b
			want := 3 + i
			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}
		})
	}
}

func TestA4(t *testing.T) {
	tests := map[string]struct {
		a    int
		b    int
		want int
	}{
		"test1": {a: 1, b: 2, want: 3},
		"test2": {a: 2, b: 3, want: 5},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.a + tt.b
			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}
