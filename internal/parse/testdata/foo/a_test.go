package foo

import (
	"fmt"
	"strconv"
	"testing"
)

type fakeRunner struct{}

func (fakeRunner) Run(name string, f func()) {
	f()
}

func TestSimpleAddition(t *testing.T) {
	a := 1
	b := 2
	got := a + b
	want := 3
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestStructSlicePositionalFields(t *testing.T) {
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

func TestStructSliceKeyedFields(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{name: "test1", a: 1, b: 2, want: 3},
		{name: "test2", a: 2, b: 3, want: 5},
		{name: "test3", a: 3, b: 4, want: 7},
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

func TestNamedStructSlice(t *testing.T) {
	type fixture struct {
		name string
		a    int
		b    int
		want int
	}
	tests := []fixture{
		{name: "test1", a: 1, b: 2, want: 3},
		{name: "test2", a: 2, b: 3, want: 5},
		{name: "test3", a: 3, b: 4, want: 7},
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

func TestVarStructSlice(t *testing.T) {
	var tests = []struct {
		name string
		a    int
		b    int
		want int
	}{
		{"test1", 1, 2, 3},
		{"test2", 2, 3, 5},
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

func TestNestedStructFieldName(t *testing.T) {
	type reqParam struct {
		name string
		a    int
		b    int
	}
	tests := []struct {
		req  reqParam
		want int
	}{
		{reqParam{"test1", 1, 2}, 3},
		{reqParam{"test2", 2, 3}, 5},
	}

	for _, tt := range tests {
		t.Run(tt.req.name, func(t *testing.T) {
			got := tt.req.a + tt.req.b
			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestMapKeySubtests(t *testing.T) {
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

func TestConcatGeneratedName(t *testing.T) {
	for i := 0; i < 3; i++ {
		t.Run("test"+strconv.Itoa(i), func(t *testing.T) {
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

func TestSprintfGeneratedName(t *testing.T) {
	for i := 0; i < 3; i++ {
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
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

func TestStringIdentNames(t *testing.T) {
	name := "test1"
	t.Run(name, func(t *testing.T) {
		a := 1
		b := 2
		got := a + b
		want := 3
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	const constName = "test2"
	t.Run(constName, func(t *testing.T) {
		a := 2
		b := 3
		got := a + b
		want := 5
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	var varName = "test3"
	t.Run(varName, func(t *testing.T) {
		a := 3
		b := 4
		got := a + b
		want := 7
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
}

func TestNonTestingRunIgnored(t *testing.T) {
	runner := fakeRunner{}
	runner.Run("not-a-subtest", func() {})
}

func TestNestedSubtestsWithRenamedTestingReceiver(t *testing.T) {
	runner := fakeRunner{}

	t.Run("outer", func(st *testing.T) {
		st.Run("inner", func(t *testing.T) {})
		runner.Run("not-a-subtest", func() {})
	})
}
