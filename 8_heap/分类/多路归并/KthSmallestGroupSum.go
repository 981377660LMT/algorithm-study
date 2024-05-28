// 多路归并, 第k小的和，第k个最小的数对和.
// 给定n个数组，从每个数组中选一个数，求这些数的和的第k小的值.

package main

import "sort"

func main() {
	groups := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	ks := NewKthSmallestGroupSum(groups...)
	ks.FirstKthSmallestSet(5, func(sum int) { println(sum) })
}

type KthSmallestGroupSum struct {
	groups [][]int
	start  int32
	sum    int
}

func NewKthSmallestGroupSum(groups ...[]int) *KthSmallestGroupSum {
	res := &KthSmallestGroupSum{groups: groups}
	sort.Slice(res.groups, func(i, j int) bool { return len(res.groups[i]) < len(res.groups[j]) })
	m := int32(len(res.groups))
	for res.start < m && len(res.groups[res.start]) == 1 {
		res.start++
	}
	for _, g := range res.groups {
		sort.Ints(g)
		for i := len(g) - 1; i >= 1; i-- {
			g[i] -= g[i-1]
		}
		res.sum += int(g[0])
	}
	arr := groups[res.start:]
	sort.Slice(arr, func(i, j int) bool { return arr[i][1] < arr[j][1] })
	return res
}

// k>=1.
func (ks *KthSmallestGroupSum) FirstKthSmallestSet(k int32, consumer func(sum int)) {
	if ks.start == int32(len(ks.groups)) {
		consumer(ks.sum)
		return
	}
	pq := newHeap(func(a, b *state) bool { return a.val < b.val }, nil)
	pq.Push(newState(ks.sum, ks.start, 0))
	m := int32(len(ks.groups))
	for k > 0 && pq.Len() > 0 {
		k--
		head := pq.Pop()
		consumer(head.val)
		if head.cid+1 < int32(len(ks.groups[head.gid])) {
			pq.Push(newState(head.val+ks.groups[head.gid][head.cid+1], head.gid, head.cid+1))
		}
		if head.cid > 0 && head.gid+1 < m {
			pq.Push(newState(head.val+ks.groups[head.gid+1][1], head.gid+1, 1))
			if head.cid == 1 {
				pq.Push(newState(head.val-ks.groups[head.gid][1]+ks.groups[head.gid+1][1], head.gid+1, 1))
			}
		}
	}
}

type state struct {
	val      int
	gid, cid int32
}

func newState(val int, gid, cid int32) *state {
	return &state{val: val, gid: gid, cid: cid}
}

func newHeap[H any](less func(a, b H) bool, nums []H) *heap[H] {
	nums = append(nums[:0:0], nums...)
	heap := &heap[H]{less: less, data: nums}
	heap.heapify()
	return heap
}

type heap[H any] struct {
	data []H
	less func(a, b H) bool
}

func (h *heap[H]) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *heap[H]) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}
	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *heap[H]) Top() (value H) {
	value = h.data[0]
	return
}

func (h *heap[H]) Len() int { return len(h.data) }

func (h *heap[H]) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *heap[H]) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *heap[H]) pushDown(root int) {
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
