// 最小斯坦纳树求方案
// !一个带权的无向图上有k个关键点，求联通k个关键点最小的代价(边权之和)。以及哪些边/点被选中。
// n≤100,m≤500,k≤10。

// O(n*3^k+mlogm*2^k)

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
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

	cost, _, _ := MinimumSteinerTree(n, edges, criticals, nil)
	fmt.Fprintln(out, cost)
}

const INF int = 1e18

// 一个联通的无向带权图上有k个关键点 criticals，求联通所有点最小的代价(边权之和)。
//  vWeights: 每个顶点的附加权重(一般为make([]int, n))。
func MinimumSteinerTree(N int, edges [][]int, criticals, vWeights []int) (cost int, es, vs []int) {
	graph := make([][][3]int, N) // (to,w,ei)
	for i, e := range edges {
		u, v, w := e[0], e[1], e[2]
		graph[u] = append(graph[u], [3]int{v, w, i})
		graph[v] = append(graph[v], [3]int{u, w, i})
	}
	if vWeights == nil {
		vWeights = make([]int, N)
	}

	M := len(edges)
	K := len(criticals)

	dp := make([][]int, 1<<K)
	for i := range dp {
		dp[i] = make([]int, N)
		for j := range dp[i] {
			dp[i][j] = INF
		}
	}
	for i := 0; i < N; i++ {
		dp[0][i] = vWeights[i]
	}

	par := make([][]int, 1<<K)
	for i := range par {
		par[i] = make([]int, N)
		for j := range par[i] {
			par[i][j] = -1
		}
	}

	for s := 1; s < 1<<K; s++ {
		curDp := dp[s]
		curS := uint(s)
		for curS > 0 {
			bit := bits.TrailingZeros(curS)
			v := criticals[bit]
			curDp[v] = min(curDp[v], dp[s^(1<<bit)][v])
			curS ^= 1 << bit
		}
		for t := (s - 1) & s; t > 0; t = (t - 1) & s {
			for v := 0; v < N; v++ {
				cand := dp[t][v] + dp[s^t][v] - vWeights[v]
				if cand < curDp[v] {
					curDp[v] = cand
					par[s][v] = 2 * t
				}
			}
		}

		inits := make([]H, N)
		for i := 0; i < N; i++ {
			inits[i] = H{curDp[i], i}
		}
		pq := NewHeap(func(a, b H) bool {
			return a[0] < b[0]
		}, inits)
		for pq.Len() > 0 {
			item := pq.Pop()
			dv, v := item[0], item[1]
			if dv != curDp[v] {
				continue
			}
			for _, e := range graph[v] {
				to, cost, id := e[0], e[1], e[2]
				cand := dv + cost + vWeights[to]
				if cand < curDp[to] {
					curDp[to] = cand
					par[s][to] = 2*id + 1
					pq.Push(H{cand, to})
				}
			}
		}
	}

	// 复元
	usedV, usedE := make([]bool, N), make([]bool, M)
	vToK := make([]int, N)
	for i := range vToK {
		vToK[i] = -1
	}
	for i := 0; i < K; i++ {
		vToK[criticals[i]] = i
	}

	root, min_ := 0, INF
	for i, v := range dp[len(dp)-1] {
		if v < min_ {
			root, min_ = i, v
		}
	}
	queue := [][2]int{{(1 << K) - 1, root}}
	usedV[root] = true

	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]
		s, v := item[0], item[1]
		if s == 0 {
			continue
		}
		if par[s][v] == -1 {
			k := vToK[v]
			queue = append(queue, [2]int{s ^ (1 << k), v})
			continue
		} else if par[s][v]&1 == 1 {
			eid := par[s][v] / 2
			e := edges[eid]
			w := v ^ e[0] ^ e[1]
			usedV[w] = true
			usedE[eid] = true
			queue = append(queue, [2]int{s, w})
			continue
		} else {
			t := par[s][v] / 2
			queue = append(queue, [2]int{t, v}, [2]int{s ^ t, v})
		}

	}

	for i := 0; i < N; i++ {
		if usedV[i] {
			vs = append(vs, i)
		}
	}
	for i := 0; i < M; i++ {
		if usedE[i] {
			es = append(es, i)
		}
	}
	for _, v := range vs {
		cost += vWeights[v]
	}
	for _, e := range es {
		cost += edges[e][2]
	}

	return
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

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}

type H = [2]int

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
	if h.Len() == 0 {
		panic("heap is empty")
	}
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
