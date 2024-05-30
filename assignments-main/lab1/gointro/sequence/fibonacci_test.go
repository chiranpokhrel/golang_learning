package sequence

import "testing"

var fibonacciTests = []struct {
	in, want uint
}{
	{0, 0},
	{1, 1},
	{2, 1},
	{3, 2},
	{4, 3},
	{5, 5},
	{6, 8},
	{7, 13},
	{8, 21},
	{9, 34},
	{10, 55},
	{20, 6765},
}

func TestFibonacci(t *testing.T) {
	for _, ft := range fibonacciTests {
		got := fibonacci(ft.in)
		if got != ft.want {
			t.Errorf("fibonacci(%q) = %q, want %q", ft.in, got, ft.want)
		}
	}
}
