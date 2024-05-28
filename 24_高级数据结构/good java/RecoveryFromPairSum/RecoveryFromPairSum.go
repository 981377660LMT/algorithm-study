// 从数对的和还原数组(RecoverArrayFromPairSum)

package main

import (
	"fmt"
	"sort"
)

func main() {
	n := 100
	sums := make([]int, 0)
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			sums = append(sums, -i-j)
		}
	}
	fmt.Println(RecoveryFromPairSum(n, sums))
}

// 给定n*(n-1)/2个两两不同的数的和，求出这n个数。

// O(n^3\log_2n)
func RecoveryFromPairSum(n int, sums []int) []int {
	if len(sums) != n*(n-1)/2 {
		panic("illegal input")
	}
	sort.Ints(sums)
	if n == 0 {
		return []int{}
	}
	if n == 1 {
		return []int{1}
	}
	if n == 2 {
		res := []int{0, sums[0]}
		if res[0] > res[1] {
			res[0], res[1] = res[1], res[0]
		}
		return res
	}
	res := make([]int, n)
	pq := NewHeap(func(a, b int) bool { return a < b }, nil)
	s12 := sums[0]
	s13 := sums[1]
	for i := 0; i < n-2; i++ {
		s23 := sums[i+2]
		first := s12 + s13 - s23
		if first%2 != 0 {
			continue
		}
		first /= 2
		if tryCase(sums, first, res, pq) {
			return res
		}
	}
	return nil
}

func tryCase(sums []int, first int, res []int, pq *Heap[int]) bool {
	n := len(res)
	wpos := 0
	res[wpos] = first
	wpos++
	pq.Clear()
	for _, x := range sums {
		if pq.Len() > 0 {
			if pq.Top() < x {
				return false
			}
			if pq.Top() == x {
				pq.Pop()
				continue
			}
		}
		if wpos >= n {
			return false
		}
		newElement := x - first
		for j := 1; j < wpos; j++ {
			pq.Push(res[j] + newElement)
		}
		res[wpos] = newElement
		wpos++
	}
	return true
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

func (h *Heap[H]) Clear() {
	h.data = h.data[:0]
}

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
