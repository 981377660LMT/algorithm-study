package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

// https://atcoder.jp/contests/abc245/tasks/abc245_g
// 给 N 个点，M 条边的无向带权图，每个点的颜色（值域为 [1,K]）。并给定 L 个点作为特殊点。
// 询问 每个点到最近的与其颜色不同的特殊点的距离（无解输出 -1） 。

// 枚举颜色的每个二进制位，把所有特殊点这一位上颜色是 1 的加入起点，跑最短路，更新所有终点中这一位上颜色是 0 的终点，
// 然后倒过来，把所有特殊点这一位上颜色是 0 的加入起点，跑最短路，更新所有终点中这一位上颜色是 1 的终点。
// 由于两个颜色不同一定至少有一个二进制位不同，因此上述算法可以保证所有终点都被颜色不同的起点更新到。 复杂度 O(mlognlogk)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, M, K, L int
	fmt.Fscan(in, &N, &M, &K, &L)
	colors := make([]int, N)
	for i := 0; i < N; i++ {
		fmt.Fscan(in, &colors[i])
		colors[i]--
	}
	criticals := make([]int, L)
	for i := 0; i < L; i++ {
		fmt.Fscan(in, &criticals[i])
		criticals[i]--
	}
	edges := make([][3]int, M)
	for i := 0; i < M; i++ {
		fmt.Fscan(in, &edges[i][0], &edges[i][1], &edges[i][2])
		edges[i][0]--
		edges[i][1]--
	}

	res := Solve(N, edges, colors, criticals)
	for _, v := range res {
		if v == INF {
			v = -1
		}
		fmt.Fprint(out, v, " ")
	}
}

const INF int = 1e18

func Solve(
	n int, edges [][3]int, colors []int, criticals []int,
) []int {
	D := NewDictionary()
	colors = append(colors[:0:0], colors...)
	for i := range colors {
		colors[i] = D.Id(colors[i])
	}

	res := make([]int, n)
	for i := range res {
		res[i] = INF
	}
	log := bits.Len(uint(maxs(colors...)))
	adjList := NewInternalCsr(n, len(edges)*2)
	cal := func(rev bool) {
		for bit := 0; bit < log; bit++ {
			var starts []int
			if rev {
				for _, v := range criticals {
					if colors[v]&(1<<bit) == 0 {
						starts = append(starts, v)
					}
				}
			} else {
				for _, v := range criticals {
					if colors[v]&(1<<bit) > 0 {
						starts = append(starts, v)
					}
				}
			}

			adjList.Clear()

			for _, edge := range edges {
				from, to, weight := edge[0], edge[1], edge[2]
				adjList.AddDirectedEdge(from, to, weight)
				adjList.AddDirectedEdge(to, from, weight)
			}

			dist := DijkstraInternalCsr(n, adjList, starts)

			if rev {
				for i := range res {
					if colors[i]&(1<<bit) > 0 {
						res[i] = min(res[i], dist[i])
					}
				}
			} else {
				for i := range res {
					if colors[i]&(1<<bit) == 0 {
						res[i] = min(res[i], dist[i])
					}
				}
			}
		}
	}

	// 1 -> 0, 0 -> 1
	cal(true)
	cal(false)
	return res
}

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

func DijkstraInternalCsr(n int, adjList *InternalCsr, starts []int) (dist []int) {
	dist = make([]int, n)
	for i := range dist {
		dist[i] = INF
	}
	pq := NewHeap(func(a, b H) bool { return a.dist < b.dist }, nil)
	for _, v := range starts {
		dist[v] = 0
		pq.Push(H{dist: 0, node: v})
	}

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

type V = int
type Dictionary struct {
	_idToValue []V
	_valueToId map[V]int
}

// A dictionary that maps values to unique ids.
func NewDictionary() *Dictionary {
	return &Dictionary{
		_valueToId: map[V]int{},
	}
}
func (d *Dictionary) Id(value V) int {
	res, ok := d._valueToId[value]
	if ok {
		return res
	}
	id := len(d._idToValue)
	d._idToValue = append(d._idToValue, value)
	d._valueToId[value] = id
	return id
}
func (d *Dictionary) Value(id int) V {
	return d._idToValue[id]
}
func (d *Dictionary) Has(value V) bool {
	_, ok := d._valueToId[value]
	return ok
}
func (d *Dictionary) Size() int {
	return len(d._idToValue)
}
