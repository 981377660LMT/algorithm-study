package main

import (
	"fmt"
)

func main() {
	nums := []int{2, 3, 1, 4, 5, 6, 7, 8, 9, 0}
	SliceStable(nums, func(i, j int) bool { return nums[i] < nums[j] })
	fmt.Println(nums)
}

// lessSwap is a pair of Less and Swap function for use with the
// auto-generated func-optimized variant of sort.go in
// zfuncversion.go.
type lessSwap struct {
	Less func(i, j int) bool
	Swap func(i, j int)
}

// SliceStable sorts the slice x using the provided less
// function, keeping equal elements in their original order.
// It panics if x is not a slice.
//
// The less function must satisfy the same requirements as
// the Interface type's Less method.
func SliceStable(x []int, less func(i, j int) bool) {
	swap := func(i, j int) { x[i], x[j] = x[j], x[i] }
	stable_func(lessSwap{less, swap}, len(x))
}

// insertionSort_func sorts data[a:b] using insertion sort.
func insertionSort_func(data lessSwap, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && data.Less(j, j-1); j-- {
			data.Swap(j, j-1)
		}
	}
}

func stable_func(data lessSwap, n int) {
	blockSize := 20 // must be > 0
	a, b := 0, blockSize
	for b <= n {
		insertionSort_func(data, a, b)
		a = b
		b += blockSize
	}
	insertionSort_func(data, a, n)

	for blockSize < n {
		a, b = 0, 2*blockSize
		for b <= n {
			symMerge_func(data, a, a+blockSize, b)
			a = b
			b += 2 * blockSize
		}
		if m := a + blockSize; m < n {
			symMerge_func(data, a, m, n)
		}
		blockSize *= 2
	}
}

// symMerge_func merges the two sorted subsequences data[a:m] and data[m:b] using
// the SymMerge algorithm from Pok-Son Kim and Arne Kutzner, "Stable Minimum
// Storage Merging by Symmetric Comparisons", in Susanne Albers and Tomasz
// Radzik, editors, Algorithms - ESA 2004, volume 3221 of Lecture Notes in
// Computer Science, pages 714-723. Springer, 2004.
//
// Let M = m-a and N = b-n. Wolog M < N.
// The recursion depth is bound by ceil(log(N+M)).
// The algorithm needs O(M*log(N/M + 1)) calls to data.Less.
// The algorithm needs O((M+N)*log(M)) calls to data.Swap.
//
// The paper gives O((M+N)*log(M)) as the number of assignments assuming a
// rotation algorithm which uses O(M+N+gcd(M+N)) assignments. The argumentation
// in the paper carries through for Swap operations, especially as the block
// swapping rotate uses only O(M+N) Swaps.
//
// symMerge assumes non-degenerate arguments: a < m && m < b.
// Having the caller check this condition eliminates many leaf recursion calls,
// which improves performance.
func symMerge_func(data lessSwap, a, m, b int) {
	// Avoid unnecessary recursions of symMerge
	// by direct insertion of data[a] into data[m:b]
	// if data[a:m] only contains one element.
	if m-a == 1 {
		// Use binary search to find the lowest index i
		// such that data[i] >= data[a] for m <= i < b.
		// Exit the search loop with i == b in case no such index exists.
		i := m
		j := b
		for i < j {
			h := int(uint(i+j) >> 1)
			if data.Less(h, a) {
				i = h + 1
			} else {
				j = h
			}
		}
		// Swap values until data[a] reaches the position before i.
		for k := a; k < i-1; k++ {
			data.Swap(k, k+1)
		}
		return
	}

	// Avoid unnecessary recursions of symMerge
	// by direct insertion of data[m] into data[a:m]
	// if data[m:b] only contains one element.
	if b-m == 1 {
		// Use binary search to find the lowest index i
		// such that data[i] > data[m] for a <= i < m.
		// Exit the search loop with i == m in case no such index exists.
		i := a
		j := m
		for i < j {
			h := int(uint(i+j) >> 1)
			if !data.Less(m, h) {
				i = h + 1
			} else {
				j = h
			}
		}
		// Swap values until data[m] reaches the position i.
		for k := m; k > i; k-- {
			data.Swap(k, k-1)
		}
		return
	}

	mid := int(uint(a+b) >> 1)
	n := mid + m
	var start, r int
	if m > mid {
		start = n - b
		r = mid
	} else {
		start = a
		r = m
	}
	p := n - 1

	for start < r {
		c := int(uint(start+r) >> 1)
		if !data.Less(p-c, c) {
			start = c + 1
		} else {
			r = c
		}
	}

	end := n - start
	if start < m && m < end {
		rotate_func(data, start, m, end)
	}
	if a < start && start < mid {
		symMerge_func(data, a, start, mid)
	}
	if mid < end && end < b {
		symMerge_func(data, mid, end, b)
	}
}

// rotate_func rotates two consecutive blocks u = data[a:m] and v = data[m:b] in data:
// Data of the form 'x u v y' is changed to 'x v u y'.
// rotate performs at most b-a many calls to data.Swap,
// and it assumes non-degenerate arguments: a < m && m < b.
func rotate_func(data lessSwap, a, m, b int) {
	i := m - a
	j := b - m

	for i != j {
		if i > j {
			swapRange_func(data, m-i, m, j)
			i -= j
		} else {
			swapRange_func(data, m-i, m+j-i, i)
			j -= i
		}
	}
	// i == j
	swapRange_func(data, m-i, m, i)
}

func swapRange_func(data lessSwap, a, b, n int) {
	for i := 0; i < n; i++ {
		data.Swap(a+i, b+i)
	}
}
