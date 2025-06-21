package foo

import "testing"

func TestB1(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		a := 1
		b := 2
		got := a + b
		want := 3
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
	t.Run("test2", func(t *testing.T) {
		a := 2
		b := 3
		got := a + b
		want := 5
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
}
