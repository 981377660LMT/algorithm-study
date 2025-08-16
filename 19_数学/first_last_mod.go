package main

import (
	"fmt"
)

// 在 [start, end) 区间内，寻找第一个和最后一个满足 x % mod == remainder 的数。
func FirstLastMod(start, end, mod, remainder int) (int, int, bool) {
	if start >= end {
		return 0, 0, false
	}
	if remainder < 0 || remainder >= mod {
		return 0, 0, false
	}
	r := start % mod
	delta := (remainder - r + mod) % mod
	first := start + delta
	if first >= end {
		return 0, 0, false
	}
	last := first + ((end-1-first)/mod)*mod
	return first, last, true
}

func main() {
	tests := []struct {
		start, end, mod, remainder int
		expectFirst, expectLast    int
		expectOk                   bool
	}{
		{1, 10, 3, 2, 2, 8, true},
		{1, 10, 3, 1, 1, 7, true},
		{1, 11, 3, 1, 1, 10, true},
		{1, 10, 3, 0, 3, 9, true},
		{0, 0, 3, 0, 0, 0, false},
		{5, 5, 3, 2, 0, 0, false},
		{10, 20, 1, 0, 10, 19, true},
		{7, 8, 2, 1, 7, 7, true},
		{7, 8, 2, 0, 0, 0, false},
		{1, 10, 3, 3, 0, 0, false},
		{1, 10, 3, -1, 0, 0, false},
	}
	for _, test := range tests {
		first, last, ok := FirstLastMod(test.start, test.end, test.mod, test.remainder)
		if ok != test.expectOk || (ok && (first != test.expectFirst || last != test.expectLast)) {
			fmt.Printf("Test failed: %+v, got (%d, %d, %v)\n", test, first, last, ok)
		}
	}
	fmt.Println("All tests passed.")
}
