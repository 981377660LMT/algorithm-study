// 给定一张n点m边的有向带权图，
// !其中有k个特殊点，问这k个点之间两两最短路的最小值是多少
// https://www.luogu.com.cn/problem/solution/P5304
// 1. 考虑将特殊点分为两个集合 A,B，s 连向 A 中的所有点， t 连向 B 中的所有点，那么 s 到 t 的最短路就是 A,B 两个集合的最短路的最小值
// 2. 对于 k 个特殊点，枚举二进制里的第 i 位，如果第 i 位为 1，那么就把这个点放到 A 集合，否则放到 B 集合，然后求出 A,B 两个集合的最短路，取最小值即可
// 3. 原理是，假设 k 个特殊点里最近的两个点是 a,b，那么 a,b 一定有一个二进制位是不同的，那么那次分组时一定被分到不同的集合里，从而肯定被算进了最后的答案之中最短路
// 4. 注意是有向图，因此分组时需要正反跑两遍

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)

	for i := 0; i < T; i++ {
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)
		edges := make([][3]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &edges[j][0], &edges[j][1], &edges[j][2])
			edges[j][0]--
			edges[j][1]--
		}
		specials := make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &specials[j])
			specials[j]--
		}
		fmt.Fprintln(out, Solve(n, edges, specials))
	}
}

func Solve(n int, edges [][3]int, specials []int) int {
	START, END := n, n+1
	res := INF
	log := bits.Len(uint(maxs(specials...)))
	adjList := NewInternalCsr(n+2, len(edges)*2)
	cal := func(rev bool) {
		for bit := 0; bit < log; bit++ {
			var groupA, groupB []int
			for _, v := range specials {
				if v&(1<<bit) > 0 {
					groupA = append(groupA, v)
				} else {
					groupB = append(groupB, v)
				}
			}

			if rev {
				groupA, groupB = groupB, groupA
			}

			adjList.Clear()

			for _, edge := range edges {
				from, to, weight := edge[0], edge[1], edge[2]
				adjList.AddDirectedEdge(from, to, weight)
			}
			for _, v := range groupA {
				adjList.AddDirectedEdge(START, v, 0)
			}
			for _, v := range groupB {
				adjList.AddDirectedEdge(v, END, 0)
			}

			dist := DijkstraInternalCsr(n+2, adjList, START)
			res = min(res, dist[END])
		}
	}

	// 注意是有向图，因此分组时需要正反跑两遍
	cal(true)
	cal(false)
	return res
}

const INF int = 1e18

type InternalCsr struct {
	to     []int32
	next   []int32
	weight []int
	head   []int32
	eid    int32
}

func NewInternalCsr(n int, m int) *InternalCsr {
	res := &InternalCsr{
		to:     make([]int32, m),
		next:   make([]int32, m),
		weight: make([]int, m),
		head:   make([]int32, n),
		eid:    0,
	}
	for i := range res.head {
		res.head[i] = -1
	}
	return res
}

func (csr *InternalCsr) AddDirectedEdge(from, to int, weight int) {
	csr.to[csr.eid] = int32(to)
	csr.next[csr.eid] = csr.head[from]
	csr.weight[csr.eid] = weight
	csr.head[from] = csr.eid
	csr.eid++
}

func (csr *InternalCsr) EnumerateNeighbors(cur int, f func(next int, weight int)) {
	for i := csr.head[cur]; i != -1; i = csr.next[i] {
		f(int(csr.to[i]), csr.weight[i])
	}
}

func (csr *InternalCsr) Clear() {
	csr.eid = 0
	for i := range csr.head {
		csr.head[i] = -1
	}
}

func (csr *InternalCsr) Copy() *InternalCsr {
	clone := &InternalCsr{}
	clone.eid = csr.eid
	clone.to = make([]int32, len(csr.to))
	copy(clone.to, csr.to)
	clone.next = make([]int32, len(csr.next))
	copy(clone.next, csr.next)
	clone.weight = make([]int, len(csr.weight))
	copy(clone.weight, csr.weight)
	clone.head = make([]int32, len(csr.head))
	copy(clone.head, csr.head)
	return clone
}

type Neighbor struct{ to, weight int }

func DijkstraInternalCsr(n int, adjList *InternalCsr, start int) (dist []int) {
	dist = make([]int, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0

	pq := NewHeap(func(a, b H) bool { return a.dist < b.dist }, []H{{dist: 0, node: start}})

	for pq.Len() > 0 {
		curNode := pq.Pop()
		cur, curDist := curNode.node, curNode.dist
		if curDist > dist[cur] {
			continue
		}

		adjList.EnumerateNeighbors(cur, func(next, weight int) {
			if cand := curDist + weight; cand < dist[next] {
				dist[next] = cand
				pq.Push(H{dist: cand, node: next})
			}
		})
	}

	return
}

type H = struct{ dist, node int }

func NewHeap(less func(a, b H) bool, data []H) *Heap {
	data = append(data[:0:0], data...)
	heap := &Heap{less: less, data: data}
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

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}
