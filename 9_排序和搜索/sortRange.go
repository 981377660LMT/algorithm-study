package sort

import (
	"internal/reflectlite"
	"math/bits"
)

// lessSwap is a pair of Less and Swap function for use with the
// auto-generated func-optimized variant of sort.go in
// zfuncversion.go.
type lessSwap struct {
	Less func(i, j int) bool
	Swap func(i, j int)
}

// Slice sorts the slice x given the provided less function.
// It panics if x is not a slice.
//
// The sort is not guaranteed to be stable: equal elements
// may be reversed from their original order.
// For a stable sort, use SliceStable.
//
// The less function must satisfy the same requirements as
// the Interface type's Less method.
func Slice(x any, less func(i, j int) bool) {
	rv := reflectlite.ValueOf(x)
	swap := reflectlite.Swapper(x)
	length := rv.Len()
	limit := bits.Len(uint(length))
	pdqsort_func(lessSwap{less, swap}, 0, length, limit)
}

// SliceStable sorts the slice x using the provided less
// function, keeping equal elements in their original order.
// It panics if x is not a slice.
//
// The less function must satisfy the same requirements as
// the Interface type's Less method.
func SliceStable(x any, less func(i, j int) bool) {
	rv := reflectlite.ValueOf(x)
	swap := reflectlite.Swapper(x)
	stable_func(lessSwap{less, swap}, rv.Len())
}

// SliceIsSorted reports whether the slice x is sorted according to the provided less function.
// It panics if x is not a slice.
func SliceIsSorted(x any, less func(i, j int) bool) bool {
	rv := reflectlite.ValueOf(x)
	n := rv.Len()
	for i := n - 1; i > 0; i-- {
		if less(i, i-1) {
			return false
		}
	}
	return true
}
