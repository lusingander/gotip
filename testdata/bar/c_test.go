package bar

import "testing"

func TestC1(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		t.Run("subtest1", func(t *testing.T) {
			a := 1
			b := 2
			got := a + b
			want := 3
			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}
		})
		t.Run("subtest2", func(t *testing.T) {
			a := 2
			b := 3
			got := a + b
			want := 5
			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}
		})
		t.Run("subtest3", func(t *testing.T) {
			t.Run("subsubtest1", func(t *testing.T) {
				a := 3
				b := 4
				got := a + b
				want := 7
				if got != want {
					t.Errorf("got %d, want %d", got, want)
				}
			})
			t.Run("subsubtest2", func(t *testing.T) {
				a := 4
				b := 5
				got := a + b
				want := 9
				if got != want {
					t.Errorf("got %d, want %d", got, want)
				}
			})
		})
	})
}
