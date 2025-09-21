// 3691. 最大子数组总值 II(多路归并)
// https://leetcode.cn/problems/maximum-total-subarray-value-ii/description/
//
// 给你一个长度为 n 的整数数组 nums 和一个整数 k。
//
// 你必须从 nums 中选择 恰好 k 个 不同 的非空子数组 nums[l..r]。子数组可以重叠，但同一个子数组（相同的 l 和 r）不能 被选择超过一次。
//
// 子数组 nums[l..r] 的 值 定义为：max(nums[l..r]) - min(nums[l..r])。
//
// 总值 是所有被选子数组的 值 之和。
//
// 返回你能实现的 最大 可能总值。
//
// 子数组 是数组中连续的 非空 元素序列。
//
// 1 <= n == nums.length <= 5 * 1e4
// 0 <= nums[i] <= 1e9
// 1 <= k <= min(1e5, n * (n + 1) / 2)

package main

import (
	"math/bits"
	"slices"
)

func maxTotalValue(nums []int, k int) int64 {
	n := len(nums)

	type minmax struct{ min, max int }
	seg := NewSparseTableG(
		n,
		func(i int) minmax {
			return minmax{nums[i], nums[i]}
		},
		func(a, b minmax) minmax {
			return minmax{min: min(a.min, b.min), max: max(a.max, b.max)}
		},
	)
	diff := func(l, r int) int {
		m := seg(l, r)
		return m.max - m.min
	}

	type heapItem struct{ diff, l, r int }
	initialHeapItems := make([]heapItem, n)
	for i := range n {
		initialHeapItems[i] = heapItem{diff: diff(i, n-1), l: i, r: n - 1}
	}
	pq := NewHeap(
		func(a, b heapItem) bool { return a.diff > b.diff },
		initialHeapItems,
	)

	res := 0
	for i := 0; i < k && pq.Len() > 0; i++ {
		cur := pq.Pop()
		res += cur.diff
		if cur.l <= cur.r-1 {
			pq.Push(heapItem{diff: diff(cur.l, cur.r-1), l: cur.l, r: cur.r - 1})
		}
	}

	return int64(res)
}

func NewHeap[H any](less func(a, b H) bool, nums []H) *Heap[H] {
	nums = slices.Clone(nums)
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

func NewSparseTableG[E any](
	n int, f func(i int) E,
	op func(E, E) E,
) (query func(int, int) E) {
	size := bits.Len(uint(n))
	dp := make([][]E, size)
	for i := range dp {
		dp[i] = make([]E, n)
	}

	for i := range n {
		dp[0][i] = f(i)
	}

	for i := 1; i < size; i++ {
		for j := 0; j+(1<<i) <= n; j++ {
			dp[i][j] = op(dp[i-1][j], dp[i-1][j+(1<<(i-1))])
		}
	}

	query = func(left, right int) E {
		k := bits.Len(uint(right-left+1)) - 1
		return op(dp[k][left], dp[k][right-(1<<k)+1])
	}

	return
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
