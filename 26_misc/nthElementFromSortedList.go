// !矩阵每行都是有序的，列不要求有序，求第 k 小的元素(在线求解).

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	time1 := time.Now()
	test()
	fmt.Println(time.Since(time1))

	matrix := [][]int{{1, 2, 9}, {3, 8, 9}, {4, 5, 6}}
	k := 4
	f := func(i, j int) int { return matrix[i][j] }
	lens := []int{3, 3, 3}
	res := NthElementFromSortedList(lens, k, f)
	fmt.Println(res)
}

func test() {
	rng := func(l, r int) int { return l + rand.Intn(r-l+1) }
	sum := func(a []int) int {
		res := 0
		for _, v := range a {
			res += v
		}
		return res
	}
	n := rng(1, 1e5)
	s := make([]int, n)
	for i := 0; i < n; i++ {
		s[i] = rng(0, 10)
	}
	dat := make([][]int, n)
	a := []int{}
	for i := 0; i < n; i++ {
		for j := 0; j < s[i]; j++ {
			a = append(a, i)
		}
	}
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	for i, v := range a {
		dat[v] = append(dat[v], i)
	}
	k := rng(0, sum(s)+1)

	f := func(i, j int) int { return dat[i][j] }
	res := NthElementFromSortedList(s, k, f)
	for i := 0; i < n; i++ {
		for j := 0; j < res[i]; j++ {
			if dat[i][j] >= k {
				panic("assertion failed")
			}
		}
		for j := res[i]; j < len(dat[i]); j++ {
			if dat[i][j] < k {
				panic("assertion failed")
			}
		}
	}
}

func kthSmallest(matrix [][]int, k int) int {
	row, col := len(matrix), len(matrix[0])
	lens := make([]int, row)
	for i := 0; i < row; i++ {
		lens[i] = col
	}
	lowerCount := NthElementFromSortedList(lens, k-1, func(i, j int) int { return matrix[i][j] })
	res := INF
	for i := 0; i < row; i++ {
		if c := lowerCount[i]; c < col {
			res = min(res, matrix[i][c])
		}
	}
	return res
}

const INF int = 1e18

// NthElementFromSortedList returns the k-th element from a sorted list.
// rowLens is the length of each row.
// f(i, j) returns the j-th element of the i-th row.
// f(i, j) < f(i, j+1) for all i, j.
// The result is 0-indexed.
// Time complexity: O(n(logn+logk)).
func NthElementFromSortedList(rowLens []int, k int, f func(i, j int) int) []int {
	return nthElementFromSortedListDfs(rowLens, k, f, 0)
}

func nthElementFromSortedListDfs(rowLens []int, k int, f func(i, j int) int, curK int) []int {
	n := len(rowLens)
	sm := 0
	for _, x := range rowLens {
		sm += x >> curK
	}
	if k == 0 {
		return make([]int, n)
	}
	if k == sm {
		return rowLens
	}
	row := 0
	for _, x := range rowLens {
		if x >= (1 << curK) {
			row++
		}
	}

	g := func(i int, j int) int {
		j = ((j + 1) << curK) - 1
		if j >= rowLens[i] {
			return INF
		} else {
			return f(i, j)
		}
	}

	var res []int
	if k > row {
		res = nthElementFromSortedListDfs(rowLens, (k-row)/2, f, curK+1)
		for i := range res {
			res[i] *= 2
		}
		k = k - (k-row)/2*2
	} else {
		res = make([]int, n)
	}

	pqNums := make([][2]int, 0, n)
	for i := 0; i < n; i++ {
		pqNums = append(pqNums, [2]int{g(i, res[i]), i})
	}
	pq := NewHeap(func(a, b [2]int) bool { return a[0] < b[0] }, pqNums...)
	for k > 0 {
		k--
		top := pq.Pop()
		i := top[1]
		res[i]++
		if k > 0 {
			pq.Push([2]int{g(i, res[i]), i})
		}
	}
	return res
}

type IHeap[T any] interface {
	Push(value T)
	Pop() T
	Replace(v T) T
	PushPop(v T) T

	Top() T
	Len() int
}

var _ IHeap[any] = (*Heap[any])(nil)

func NewHeap[H any](less func(a, b H) bool, nums ...H) *Heap[H] {
	res := &Heap[H]{less: less, data: append(nums[:0:0], nums...)}
	if len(nums) > 1 {
		res.heapify()
	}
	return res
}

func NewHeapWithCapacity[H any](less func(a, b H) bool, capacity int32, nums ...H) *Heap[H] {
	if n := int32(len(nums)); capacity < n {
		capacity = n
	}
	res := &Heap[H]{less: less, data: make([]H, 0, capacity)}
	res.data = append(res.data, nums...)
	if len(nums) > 1 {
		res.heapify()
	}
	return res
}

type Heap[H any] struct {
	less func(a, b H) bool
	data []H
}

func (h *Heap[H]) Push(value H) {
	h.data = append(h.data, value)
	h.up(h.Len() - 1)
}

func (h *Heap[H]) Pop() H {
	n := h.Len() - 1
	h.data[0], h.data[n] = h.data[n], h.data[0]
	h.down(0, n)
	res := h.data[n]
	h.data = h.data[:n]
	return res
}

func (h *Heap[H]) Top() H {
	return h.data[0]
}

// replace 弹出并返回堆顶，同时将 v 入堆.
// 需保证 h 非空.
func (h *Heap[H]) Replace(v H) H {
	top := h.Top()
	h.data[0] = v
	h.fix(0)
	return top
}

// pushPop 先将 v 入堆，然后弹出并返回堆顶.
func (h *Heap[H]) PushPop(v H) H {
	data, less := h.data, h.less
	if len(data) > 0 && less(data[0], v) {
		v, data[0] = data[0], v
		h.fix(0)
	}
	return v
}

func (h *Heap[H]) Len() int { return len(h.data) }

func (h *Heap[H]) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.down(i, n)
	}
}

func (h *Heap[H]) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.less(h.data[j], h.data[i]) {
			break
		}
		h.data[i], h.data[j] = h.data[j], h.data[i]
		j = i
	}
}

func (h *Heap[H]) down(i0, n int) bool {
	i := i0
	for {
		j1 := (i << 1) | 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.less(h.data[j2], h.data[j1]) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.less(h.data[j], h.data[i]) {
			break
		}
		h.data[i], h.data[j] = h.data[j], h.data[i]
		i = j
	}
	return i > i0
}

func (h *Heap[H]) fix(i int) {
	if !h.down(i, h.Len()) {
		h.up(i)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
