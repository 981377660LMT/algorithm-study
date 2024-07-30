package main

import (
	"fmt"
	"math/rand"
)

func main() {
	rng := func(l, r int) int { return l + rand.Intn(r-l+1) }
	sum := func(a []int) int {
		res := 0
		for _, v := range a {
			res += v
		}
		return res
	}
	n := rng(1, 3)
	s := make([]int, n)
	for i := 0; i < n; i++ {
		s[i] = rng(0, 3)
	}
	dat := make([][]int, n)
	a := []int{}
	for i := 0; i < n; i++ {
		for j := 0; j < s[i]; j++ {
			a = append(a, i)
		}
	}
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	for _, x := range a {
		dat[x] = append(dat[x], len(dat[x]))
	}
	k := rng(0, sum(s)+1)
	fmt.Println(dat, k)
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

const INF int = 1e18

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

	pq := NewHeap(func(a, b [2]int) bool { return a[0] < b[0] }, nil)
	for i := 0; i < n; i++ {
		pq.Push([2]int{g(i, res[i]), i})
	}
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

func NewHeap[H any](less func(a, b H) bool, nums []H) *Heap[H] {
	nums = append(nums[:0:0], nums...)
	heap := &Heap[H]{less: less, data: nums}
	heap.heapify()
	return heap
}

type Heap[H any] struct {
	data []H
	less func(a, b H) bool
}

func (h *Heap[H]) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap[H]) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}
	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap[H]) Top() (value H) {
	value = h.data[0]
	return
}

func (h *Heap[H]) Len() int { return len(h.data) }

func (h *Heap[H]) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *Heap[H]) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap[H]) pushDown(root int) {
	n := h.Len()
	for left := (root<<1 + 1); left < n; left = (root<<1 + 1) {
		right := left + 1
		minIndex := root
		if h.less(h.data[left], h.data[minIndex]) {
			minIndex = left
		}
		if right < n && h.less(h.data[right], h.data[minIndex]) {
			minIndex = right
		}
		if minIndex == root {
			return
		}
		h.data[root], h.data[minIndex] = h.data[minIndex], h.data[root]
		root = minIndex
	}
}
