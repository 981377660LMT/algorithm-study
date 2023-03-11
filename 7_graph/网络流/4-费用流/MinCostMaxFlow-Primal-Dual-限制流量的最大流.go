// https://ei1333.github.io/luzhiled/snippets/graph/primal-dual.html

// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=GRL_6_B&lang=jp
// !Primal-Dual O(f*ElogV)
// 限制流量为f时的最小费用,如果不能满足流量等于f,则返回-1

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, f int
	fmt.Fscan(in, &n, &m, &f)
	pd := NewPrimalDual(n)
	for i := 0; i < m; i++ {
		var from, to, cap, cost int
		fmt.Fscan(in, &from, &to, &cap, &cost)
		pd.AddEdge(from, to, cap, cost)
	}
	minCost := pd.MinCostFlow(0, n-1, f)
	fmt.Fprintln(out, minCost)
	fmt.Fprintln(out, pd.GetEdges())
}

const INF int = 1e18

type PrimalDual struct {
	graph              [][]edge
	potential, minCost []int
	prevv, preve       []int
}

type edge struct {
	to    int
	cap   int
	cost  int
	rev   int
	isRev bool
}

// 頂点数 vで初期化する.
func NewPrimalDual(n int) *PrimalDual {
	return &PrimalDual{
		graph: make([][]edge, n),
	}
}

// 頂点 from から to に容量 cap、コスト cost の有向辺を張る.
func (p *PrimalDual) AddEdge(from, to, cap, cost int) {
	p.graph[from] = append(p.graph[from], edge{to, cap, cost, len(p.graph[to]), false})
	p.graph[to] = append(p.graph[to], edge{from, 0, -cost, len(p.graph[from]) - 1, true})
}

// 頂点 s から t に流量 f の最小費用流を流し, そのコストを返す.
//  流せないとき −1を返す.
func (pd *PrimalDual) MinCostFlow(start, target, f int) int {
	v := len(pd.graph)
	res := 0
	que := NewHeap(func(a, b H) int {
		return a[0] - b[0]
	}, nil)
	pd.potential = make([]int, v)
	pd.prevv = make([]int, v)
	pd.preve = make([]int, v)
	for i := 0; i < v; i++ {
		pd.prevv[i] = -1
		pd.preve[i] = -1
	}

	for f > 0 {
		pd.minCost = make([]int, v)
		for i := 0; i < v; i++ {
			pd.minCost[i] = INF
		}

		que.Push(H{0, start})
		pd.minCost[start] = 0
		for que.Len() > 0 {
			p := que.Pop()
			if pd.minCost[p[1]] < p[0] {
				continue
			}

			for i := 0; i < len(pd.graph[p[1]]); i++ {
				e := pd.graph[p[1]][i]
				nextCost := pd.minCost[p[1]] + e.cost + pd.potential[p[1]] - pd.potential[e.to]
				if e.cap > 0 && pd.minCost[e.to] > nextCost {
					pd.minCost[e.to] = nextCost
					pd.prevv[e.to] = p[1]
					pd.preve[e.to] = i
					que.Push(H{pd.minCost[e.to], e.to})
				}
			}
		}

		if pd.minCost[target] == INF {
			return -1
		}

		for i := 0; i < v; i++ {
			pd.potential[i] += pd.minCost[i]
		}

		addFlow := f
		for v := target; v != start; v = pd.prevv[v] {
			addFlow = min(addFlow, pd.graph[pd.prevv[v]][pd.preve[v]].cap)
		}
		f -= addFlow
		res += addFlow * pd.potential[target]
		for v := target; v != start; v = pd.prevv[v] {
			e := &pd.graph[pd.prevv[v]][pd.preve[v]] // !ptr
			e.cap -= addFlow
			pd.graph[v][e.rev].cap += addFlow
		}
	}

	return res
}

// 最小費用流を復元する (from, to, flow, cap).
func (p *PrimalDual) GetEdges() [][4]int {
	res := make([][4]int, 0)
	for i := 0; i < len(p.graph); i++ {
		for _, e := range p.graph[i] {
			if e.isRev {
				continue
			}
			revEdge := p.graph[e.to][e.rev]
			res = append(res, [4]int{i, e.to, revEdge.cap, revEdge.cap + e.cap})
		}
	}
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type H = [2]int

// Should return a number:
//    negative , if a < b
//    zero     , if a == b
//    positive , if a > b
type Comparator func(a, b H) int

func NewHeap(comparator Comparator, nums []H) *Heap {
	nums = append(nums[:0:0], nums...)
	heap := &Heap{comparator: comparator, data: nums}
	heap.heapify()
	return heap
}

type Heap struct {
	data       []H
	comparator Comparator
}

func (h *Heap) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap) Pop() (value H) {
	if h.Len() == 0 {
		return
	}

	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap) Peek() (value H) {
	if h.Len() == 0 {
		return
	}
	value = h.data[0]
	return
}

func (h *Heap) Len() int { return len(h.data) }

func (h *Heap) heapify() {
	for i := (h.Len() >> 1) - 1; i >= 0; i-- {
		h.pushDown(i)
	}
}

func (h *Heap) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.comparator(h.data[root], h.data[parent]) < 0; parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap) pushDown(root int) {
	n := h.Len()
	for left := (root<<1 + 1); left < n; left = (root<<1 + 1) {
		right := left + 1
		minIndex := root

		if h.comparator(h.data[left], h.data[minIndex]) < 0 {
			minIndex = left
		}

		if right < n && h.comparator(h.data[right], h.data[minIndex]) < 0 {
			minIndex = right
		}

		if minIndex == root {
			return
		}

		h.data[root], h.data[minIndex] = h.data[minIndex], h.data[root]
		root = minIndex
	}
}
