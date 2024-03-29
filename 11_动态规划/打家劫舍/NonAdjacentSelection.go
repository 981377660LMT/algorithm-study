package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	// nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// nas := NewNonAdjacentSelection(nums, true)
	// fmt.Println(nas.Restore(5))

	CF958E2()
}

// 给定n个时间点。每个区间都以某两个时间点为左右端点，且每个区间的「代价」定义端点的时间之差。
// !你要选择k个连续的区间，保证这个k个连续的区间没有交集，且代价总和最小。
// n<=5e5,k<=5000
// https://www.luogu.com.cn/problem/CF958E2
func CF958E2() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &k, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	sort.Ints(nums)
	diff := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		diff[i] = nums[i+1] - nums[i]
	}
	S := NewNonAdjacentSelection(diff, true)
	res := S.Solve()
	fmt.Fprintln(out, res[k])
}

// 从数组不相邻选择 k(0<=k<=(n+1/2)) 个数,最大化和/最小化和.
type NonAdjacentSelection struct {
	n        int
	nums     []int
	minimize bool
	history  [][2]int
	solved   bool
}

func NewNonAdjacentSelection(nums []int, minimize bool) *NonAdjacentSelection {
	return &NonAdjacentSelection{
		n:        len(nums),
		nums:     nums,
		minimize: minimize,
	}
}

func (nas *NonAdjacentSelection) Solve() []int {
	if nas.minimize {
		tmp := make([]int, len(nas.nums))
		for i := 0; i < len(nas.nums); i++ {
			tmp[i] = -nas.nums[i]
		}
		nas.nums = tmp
	}

	nums := nas.nums
	history := [][2]int{}
	n := nas.n
	rest := make([]bool, n+2)
	for i := 1; i < n+1; i++ {
		rest[i] = true
	}
	left, right := make([]int, n+2), make([]int, n+2)
	for i := 0; i < n+2; i++ {
		left[i] = i - 1
		right[i] = i + 1
	}
	range_ := make([][2]int, n+2)
	for i := 1; i < n+1; i++ {
		range_[i] = [2]int{i - 1, i}
	}
	val := make([]int, n+2)
	for i := 1; i < n+1; i++ {
		val[i] = nums[i-1]
	}

	pqNums := make([]H, n)
	for i := 0; i < n; i++ {
		pqNums[i] = H{value: val[i+1], index: i + 1}
	}
	pq := NewHeap(func(a, b H) bool { return a.value > b.value }, pqNums)

	res := make([]int, 0, ((n+1)/2)+1)
	res = append(res, 0)
	for pq.Len() > 0 {
		item := pq.Pop()
		add, i := item.value, item.index
		if !rest[i] {
			continue
		}
		res = append(res, res[len(res)-1]+add)
		L, R := left[i], right[i]
		history = append(history, range_[i])
		if 1 <= L {
			right[left[L]] = i
			left[i] = left[L]
		}
		if R <= n {
			left[right[R]] = i
			right[i] = right[R]
		}
		if rest[L] && rest[R] {
			val[i] = val[L] + val[R] - val[i]
			pq.Push(H{value: val[i], index: i})
			range_[i] = [2]int{range_[L][0], range_[R][1]}
		} else {
			rest[i] = false
		}
		rest[L] = false
		rest[R] = false
	}

	if nas.minimize {
		for i := range res {
			res[i] = -res[i]
		}
	}

	nas.history = history
	nas.solved = true
	return res
}

// 选择k个数,使得和最大/最小,返回选择的数的下标.
// 0 <= k <= (n+1) / 2.
func (nas *NonAdjacentSelection) Restore(k int) []int {
	if k < 0 || k > (nas.n+1)/2 {
		panic("k must be in [0,(n+1)/2]")
	}
	if !nas.solved {
		nas.Solve()
	}
	diff := make([]int, nas.n+1)
	for i := 0; i < k; i++ {
		item := nas.history[i]
		a, b := item[0], item[1]
		diff[a]++
		diff[b]--
	}
	for i := 1; i < nas.n+1; i++ {
		diff[i] += diff[i-1]
	}
	res := make([]int, 0, k)
	for i := 0; i < nas.n; i++ {
		if diff[i]&1 == 1 {
			res = append(res, i)
		}
	}
	return res
}

type H = struct {
	value int
	index int
}

func NewHeap(less func(a, b H) bool, nums []H) *Heap {
	nums = append(nums[:0:0], nums...)
	heap := &Heap{less: less, data: nums}
	heap.heapify()
	return heap
}

type Heap struct {
	data []H
	less func(a, b H) bool
}

func (h *Heap) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}

	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap) Top() (value H) {
	value = h.data[0]
	return
}

func (h *Heap) Len() int { return len(h.data) }

func (h *Heap) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *Heap) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap) pushDown(root int) {
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
