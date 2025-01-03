// 并行排序
// 核心思想是利用多个线程对数据进行并行的分块排序（bucket sort），
// 然后递归地进行对称归并（Symmetric Merge）.

package merge

import (
	"math"
	"runtime"
	"sort"
	"sync"
)

func main() {

}

func sortBucket(comparators Comparators) {
	sort.Sort(comparators)
}

func copyChunk(chunk []Comparators) []Comparators {
	cp := make([]Comparators, len(chunk))
	copy(cp, chunk)
	return cp
}

// MultithreadedSortComparators will take a list of comparators
// and sort it using as many threads as are available.  The list
// is split into buckets for a bucket sort and then recursively
// merged using SymMerge.
func MultithreadedSortComparators(comparators Comparators) Comparators {
	toBeSorted := make(Comparators, len(comparators))
	copy(toBeSorted, comparators)

	var wg sync.WaitGroup

	numCPU := int64(runtime.NumCPU())
	if numCPU == 1 { // single core machine
		numCPU = 2
	} else {
		// otherwise this algo only works with a power of two
		numCPU = int64(prevPowerOfTwo(uint64(numCPU)))
	}

	chunks := chunk(toBeSorted, numCPU)
	wg.Add(len(chunks))
	for i := 0; i < len(chunks); i++ {
		go func(i int) {
			sortBucket(chunks[i])
			wg.Done()
		}(i)
	}

	wg.Wait()

	// 对每对相邻的块进行对称归并（SymMerge），结果存储在 todo 切片中
	todo := make([]Comparators, len(chunks)/2)
	for {
		todo = todo[:len(chunks)/2]
		wg.Add(len(chunks) / 2)
		for i := 0; i < len(chunks); i += 2 {
			go func(i int) {
				todo[i/2] = SymMerge(chunks[i], chunks[i+1])
				wg.Done()
			}(i)
		}

		wg.Wait()

		chunks = copyChunk(todo)
		if len(chunks) == 1 {
			break
		}
	}

	return chunks[0]
}

func chunk(comparators Comparators, numParts int64) []Comparators {
	parts := make([]Comparators, numParts)
	for i := int64(0); i < numParts; i++ {
		parts[i] = comparators[i*int64(len(comparators))/numParts : (i+1)*int64(len(comparators))/numParts]
	}
	return parts
}

func prevPowerOfTwo(x uint64) uint64 {
	x = x | (x >> 1)
	x = x | (x >> 2)
	x = x | (x >> 4)
	x = x | (x >> 8)
	x = x | (x >> 16)
	x = x | (x >> 32)
	return x - (x >> 1)
}

// #region interface

// Comparators defines a typed list of type Comparator.
type Comparators []Comparator

// Less returns a bool indicating if the comparator at index i
// is less than the comparator at index j.
func (c Comparators) Less(i, j int) bool {
	return c[i].Compare(c[j]) < 0
}

// Len returns an int indicating the length of this list
// of comparators.
func (c Comparators) Len() int {
	return len(c)
}

// Swap swaps the values at positions i and j.
func (c Comparators) Swap(i, j int) {
	c[j], c[i] = c[i], c[j]
}

// Comparator defines items that can be sorted.  It contains
// a single method allowing the compare logic to compare one
// comparator to another.
type Comparator interface {
	// Compare will return a value indicating how this comparator
	// compares with the provided comparator.  A negative number
	// indicates this comparator is less than the provided comparator,
	// a 0 indicates equality, and a positive number indicates this
	// comparator is greater than the provided comparator.
	Compare(Comparator) int
}

// #endregion

// #region symmerge

// symSearch is like symBinarySearch but operates
// on two sorted lists instead of a sorted list and an index.
// It's duplication of code but you buy performance.
func symSearch(u, w Comparators) int {
	start, stop, p := 0, len(u), len(w)-1
	for start < stop {
		mid := (start + stop) / 2
		if u[mid].Compare(w[p-mid]) <= 0 {
			start = mid + 1
		} else {
			stop = mid
		}
	}

	return start
}

// swap will swap positions of the two lists from index
// to the end of the list.  It expects that these lists
// are the same size or one different.
func swap(u, w Comparators, index int) {
	for i := index; i < len(u); i++ {
		u[i], w[i-index] = w[i-index], u[i]
	}
}

// decomposeForSymMerge pulls an active site out of the list
// of length in size.  W becomes the active site for future sym
// merges and v1, v2 are decomposed and split among the other
// list to be merged and w.
func decomposeForSymMerge(length int, comparators Comparators) (v1 Comparators, w Comparators, v2 Comparators) {

	if length >= len(comparators) {
		panic(`INCORRECT PARAMS FOR SYM MERGE.`)
	}

	overhang := (len(comparators) - length) / 2
	v1 = comparators[:overhang]
	w = comparators[overhang : overhang+length]
	v2 = comparators[overhang+length:]
	return
}

// symBinarySearch will perform a binary search between the provided
// indices and find the index at which a rotation should occur.
func symBinarySearch(u Comparators, start, stop, total int) int {
	for start < stop {
		mid := (start + stop) / 2
		if u[mid].Compare(u[total-mid]) <= 0 {
			start = mid + 1
		} else {
			stop = mid
		}
	}

	return start
}

// symSwap will perform a rotation or swap between the provided
// indices.  Again, there is duplication here with swap, but
// we are buying performance.
func symSwap(u Comparators, start1, start2, end int) {
	for i := 0; i < end; i++ {
		u[start1+i], u[start2+i] = u[start2+i], u[start1+i]
	}
}

// symRotate determines the indices to use in a symSwap and
// performs the swap.
func symRotate(u Comparators, start1, start2, end int) {
	i := start2 - start1
	if i == 0 {
		return
	}

	j := end - start2
	if j == 0 {
		return
	}

	if i == j {
		symSwap(u, start1, start2, i)
		return
	}

	p := start1 + i
	for i != j {
		if i > j {
			symSwap(u, p-i, p, j)
			i -= j
		} else {
			symSwap(u, p-i, p+j-i, i)
			j -= i
		}
	}
	symSwap(u, p-i, p, i)
}

// symMerge is the recursive and internal form of SymMerge.
func symMerge(u Comparators, start1, start2, last int) {
	if start1 < start2 && start2 < last {
		mid := (start1 + last) / 2
		n := mid + start2
		var start int
		if start2 > mid {
			start = symBinarySearch(u, n-last, mid, n-1)
		} else {
			start = symBinarySearch(u, start1, start2, n-1)
		}
		end := n - start

		symRotate(u, start, start2, end)
		symMerge(u, start1, start, mid)
		symMerge(u, mid, end, last)
	}
}

// SymMerge will perform a symmetrical merge of the two provided
// lists.  It is expected that these lists are pre-sorted.  Failure
// to do so will result in undefined behavior.  This function does
// make use of goroutines, so multithreading can aid merge time.
// This makes M*log(N/M+1) comparisons where M is the length
// of the shorter list and N is the length of the longer list.
func SymMerge(u, w Comparators) Comparators {
	lenU, lenW := len(u), len(w)
	if lenU == 0 {
		return w
	}

	if lenW == 0 {
		return u
	}

	diff := lenU - lenW
	if math.Abs(float64(diff)) > 1 {
		u1, w1, u2, w2 := prepareForSymMerge(u, w)

		lenU1 := len(u1)
		lenU2 := len(u2)
		u = append(u1, w1...)
		w = append(u2, w2...)
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			symMerge(u, 0, lenU1, len(u))
			wg.Done()
		}()
		go func() {
			symMerge(w, 0, lenU2, len(w))
			wg.Done()
		}()

		wg.Wait()
		u = append(u, w...)
		return u
	}

	u = append(u, w...)
	symMerge(u, 0, lenU, len(u))
	return u
}

// prepareForSymMerge performs a symmetrical decomposition on two
// lists of different sizes.  It breaks apart the longer list into
// an active site (equal to the size of the shorter list) and performs
// a symmetrical rotation with the active site and the shorter list.
// The two stubs are then split between the active site and shorter list
// ensuring two equally sized lists where every value in u' is less
// than w'.
func prepareForSymMerge(u, w Comparators) (u1, w1, u2, w2 Comparators) {
	if u.Len() > w.Len() {
		u, w = w, u
	}
	v1, w, v2 := decomposeForSymMerge(len(u), w)

	i := symSearch(u, w)

	u1 = make(Comparators, i)
	copy(u1, u[:i])
	w1 = append(v1, w[:len(w)-i]...)

	u2 = make(Comparators, len(u)-i)
	copy(u2, u[i:])

	w2 = append(w[len(w)-i:], v2...)
	return
}

// #endregion
