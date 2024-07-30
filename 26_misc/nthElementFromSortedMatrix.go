// https://uoj.ac/problem/891
// !矩阵每行和每列都是有序的，求第 k 小的元素(在线求解).

package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

func main() {
	time1 := time.Now()
	test()
	fmt.Println(time.Since(time1))
}

const INF int = 1e9

// https://leetcode.cn/problems/kth-smallest-element-in-a-sorted-matrix/description/
// 其中每行和每列元素均按升序排序，找到矩阵中第 k 小的元素
func kthSmallest(matrix [][]int, k int) int {
	row, col := len(matrix), len(matrix[0])
	lowerCount := NthElementFronSortedMatrix(row, col, k-1, func(i, j int) int { return matrix[i][j] })
	res := INF
	for i := 0; i < row; i++ {
		if c := lowerCount[i]; c < col {
			res = min(res, matrix[i][c])
		}
	}
	return res
}

func test() {
	rng := func(l, r int) int { return l + rand.Intn(r-l+1) }
	n := rng(1, 1000)
	m := rng(1, 1000)
	a := make([][]int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int, m)
	}
	ij := make([][2]int, 0, n*m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			ij = append(ij, [2]int{i, j})
		}
	}
	rand.Shuffle(len(ij), func(i, j int) { ij[i], ij[j] = ij[j], ij[i] })
	for k := range ij {
		i, j := ij[k][0], ij[k][1]
		a[i][j] = k
	}
	for i := 0; i < n; i++ {
		sort.Ints(a[i])
	}

	for j := 0; j < m; j++ {
		col := make([]int, n)
		for i := 0; i < n; i++ {
			col[i] = a[i][j]
		}
		sort.Ints(col)
		for i := 0; i < n; i++ {
			a[i][j] = col[i]
		}
	}
	if rng(0, 2) == 1 {
		b := make([][]int, m)
		for i := 0; i < m; i++ {
			b[i] = make([]int, n)
		}
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				b[j][i] = a[i][j]
			}
		}
		a, b = b, a
		n, m = m, n
	}

	k := rng(0, n*m+1)

	f := func(i, j int) int { return a[i][j] }

	res := NthElementFronSortedMatrix(n, m, k, f)
	for i := 0; i < n; i++ {
		for j := 0; j < res[i]; j++ {
			if a[i][j] >= k {
				panic("assertion failed")
			}
		}
		for j := res[i]; j < m; j++ {
			if a[i][j] < k {
				panic("assertion failed")
			}
		}
	}

}

func NthElementFronSortedMatrix(n, m, k int, f func(int, int) int) []int {
	return nthElementFronSortedMatrixdfs(n, m, k, f, 0, 0, false)
}

func nthElementFronSortedMatrixdfs(n, m, k int, f func(int, int) int, k1, k2 int, tr bool) []int {
	if k == 0 {
		return make([]int, n)
	}
	if n > m {
		tmp := nthElementFronSortedMatrixdfs(m, n, k, f, k2, k1, !tr)
		b := make([]int, n+1)
		for i := 0; i < m; i++ {
			b[0] += 1
			b[tmp[i]] -= 1
		}
		for i := 0; i < n; i++ {
			b[i+1] += b[i]
		}
		b = b[:len(b)-1]
		return b
	}
	var res []int
	if k > n {
		res = nthElementFronSortedMatrixdfs(n, m/2, (k-n)/2, f, k1, k2+1, tr)
		for i := range res {
			res[i] *= 2
		}
		k = k - (k-n)/2*2
	} else {
		res = make([]int, n)
	}

	g := func(i int, j int) int {
		i = ((i + 1) << k1) - 1
		j = ((j + 1) << k2) - 1
		if tr {
			return f(j, i)
		} else {
			return f(i, j)
		}
	}
	var pqNums [][2]int
	if res[0] < m {
		pqNums = append(pqNums, [2]int{g(0, res[0]), 0})
	}
	for i := 1; i < n; i++ {
		if res[i] < res[i-1] {
			pqNums = append(pqNums, [2]int{g(i, res[i]), i})
		}
	}
	pq := NewHeap(func(a, b [2]int) bool { return a[0] < b[0] }, pqNums...)
	for k > 0 {
		k--
		top := pq.Pop()
		i := top[1]
		res[i]++
		if k == 0 {
			break
		}
		if res[i] < m && (i == 0 || res[i-1] > res[i]) {
			pq.Push([2]int{g(i, res[i]), i})
		}
		if i+1 < n && res[i+1] == res[i]-1 {
			pq.Push([2]int{g(i+1, res[i+1]), i + 1})
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
