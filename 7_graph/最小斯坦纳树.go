// 最小斯坦纳树
// !一个带权的无向图上有k个关键点，求联通k个关键点最小的代价(边权之和)。
// n≤100,m≤500,k≤10。

// O(n*3^k+mlogm*2^k)

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// abc364_g()
	p6192()
}

const INF int = 1e18

// G - Last Major City
// https://atcoder.jp/contests/abc364/tasks/abc364_g
// n个点 m条边，边有边权。
// 前 k−1个点是重要点。
//
// 依次解决以下问题。
// 分别选定 k→n的节点为重要点，问每个情况下，由重要点构成的最小生成树的权值。
//
// !dp[bit][i] 表示 bit状态下，以i为根的最小生成树的权值。
func abc364_g() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)
	k--
	graph := make([][][2]int, n)
	for i := 0; i < m; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		u, v = u-1, v-1
		graph[u] = append(graph[u], [2]int{v, w})
		graph[v] = append(graph[v], [2]int{u, w})
	}

	dp := make([][]int, 1<<k) // !dp[bit][i] 表示 bit状态下，以i为根的最小生成树的权值。
	for i := range dp {
		dp[i] = make([]int, n)
		for j := range dp[i] {
			dp[i][j] = INF
		}
	}
	for i := 0; i < k; i++ {
		dp[1<<i][i] = 0
	}

	pq := NewHeap(func(a, b [2]int) bool { return a[0] < b[0] }, nil)

	for bit := 1; bit < (1 << k); bit++ {
		for sub := bit; sub > 0; sub = (sub - 1) & bit {
			for i := 0; i < n; i++ {
				dp[bit][i] = min(dp[bit][i], dp[sub][i]+dp[bit-sub][i])
			}
		}

		pq.Clear()
		for i := 0; i < n; i++ {
			pq.Push([2]int{dp[bit][i], i})
		}
		for pq.Len() > 0 {
			tmp := pq.Pop()
			d, u := tmp[0], tmp[1]
			if d > dp[bit][u] {
				continue
			}
			for _, e := range graph[u] {
				v, w := e[0], e[1]
				if dp[bit][u]+w < dp[bit][v] {
					dp[bit][v] = dp[bit][u] + w
					pq.Push([2]int{dp[bit][v], v})
				}
			}
		}
	}

	for i := k; i < n; i++ {
		fmt.Fprintln(out, dp[len(dp)-1][i])
	}
}

func p6192() {
	// https://www.luogu.com.cn/problem/P6192
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)
	edges := make([][]int, 0, m)
	for i := 0; i < m; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		u, v = u-1, v-1
		edges = append(edges, []int{u, v, w})
	}

	criticals := make([]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &criticals[i])
		criticals[i]--
	}

	fmt.Fprintln(out, MinimumSteinerTree(n, edges, criticals))
}

// 一个联通的无向带权图上有k个关键点 criticals，求联通k个关键点最小的代价(边权之和)。
func MinimumSteinerTree(n int, edges [][]int, criticals []int) int {
	k := len(criticals)
	graph := make([][][2]int, n)
	for _, e := range edges {
		u, v, w := e[0], e[1], e[2]
		graph[u] = append(graph[u], [2]int{v, w})
		graph[v] = append(graph[v], [2]int{u, w})
	}

	dp := make([][]int, 1<<k) // dp[bit][i] 表示 bit状态下，以i为根的最小生成树的权值。
	for i := range dp {
		dp[i] = make([]int, n)
		for j := range dp[i] {
			dp[i][j] = INF
		}
	}

	for i := 0; i < k; i++ {
		dp[1<<i][criticals[i]] = 0
	}

	pq := NewHeap(func(a, b [2]int) bool { return a[0] < b[0] }, nil)

	for bit := 1; bit < (1 << k); bit++ {
		for sub := bit; sub > 0; sub = (sub - 1) & bit {
			for i := 0; i < n; i++ {
				dp[bit][i] = min(dp[bit][i], dp[sub][i]+dp[bit-sub][i])
			}
		}

		pq.Clear()
		for i := 0; i < n; i++ {
			pq.Push([2]int{dp[bit][i], i})
		}
		for pq.Len() > 0 {
			tmp := pq.Pop()
			d, u := tmp[0], tmp[1]
			if d > dp[bit][u] {
				continue
			}
			for _, e := range graph[u] {
				v, w := e[0], e[1]
				if dp[bit][u]+w < dp[bit][v] {
					dp[bit][v] = dp[bit][u] + w
					pq.Push([2]int{dp[bit][v], v})
				}
			}
		}
	}

	res := INF
	for i := 0; i < n; i++ {
		res = min(res, dp[len(dp)-1][i])
	}
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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

func (h *Heap[H]) Clear() {
	h.data = h.data[:0]
}
