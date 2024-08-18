// 带流量限制的最小费用流.
// https://ei1333.github.io/luzhiled/snippets/graph/primal-dual.html
// !Primal-Dual O(f*ElogV)
// 限制流量为f时的最小费用,如果不能满足流量等于f,则返回-1

package main

// 100401. 放三个车的价值之和最大 II
// https://leetcode.cn/problems/maximum-value-sum-by-placing-three-rooks-ii/
func maximumValueSum(board [][]int) int64 {
	ROW, COL := int32(len(board)), int32(len(board[0]))
	M := NewPrimalDual(ROW + COL + 2)
	S, T := ROW+COL, ROW+COL+1
	for r := int32(0); r < ROW; r++ {
		M.AddEdge(S, r, 1, 0)
	}
	for r := int32(0); r < ROW; r++ {
		for c := int32(0); c < COL; c++ {
			M.AddEdge(r, ROW+c, 1, -board[r][c]) // 求最大值，取负数
		}
	}
	for c := int32(0); c < COL; c++ {
		M.AddEdge(ROW+c, T, 1, 0)
	}
	res := -M.MinCostFlow(S, T, 3) // !放三个车，流量为3
	return int64(res)
}

const INF int = 1e18

type PrimalDual struct {
	prevv, preve       []int32
	potential, minCost []int
	graph              [][]edge
}

type edge struct {
	isRev bool
	to    int32
	rev   int32
	cap   int
	cost  int
}

func NewPrimalDual(n int32) *PrimalDual {
	return &PrimalDual{graph: make([][]edge, n)}
}

// 顶点 from 到顶点 to 的容量为 cap, 费用为 cost 的边.
func (p *PrimalDual) AddEdge(from, to int32, cap, cost int) {
	p.graph[from] = append(p.graph[from], edge{isRev: false, to: to, cap: cap, cost: cost, rev: int32(len(p.graph[to]))})
	p.graph[to] = append(p.graph[to], edge{isRev: true, to: from, cap: 0, cost: -cost, rev: int32(len(p.graph[from]) - 1)})
}

// 从顶点 s 到顶点 t 流量为 f 的最小费用流, 返回其费用.
// 如果不存在，返回-1.
func (pd *PrimalDual) MinCostFlow(start, target int32, f int) int {
	v := len(pd.graph)
	res := 0
	type pair = struct {
		first  int
		second int32
	}
	que := NewHeap[pair](func(a, b pair) bool { return a.first < b.first }, nil)
	pd.potential = make([]int, v)
	pd.prevv = make([]int32, v)
	pd.preve = make([]int32, v)
	for i := 0; i < v; i++ {
		pd.prevv[i] = -1
		pd.preve[i] = -1
	}

	for f > 0 {
		pd.minCost = make([]int, v)
		for i := 0; i < v; i++ {
			pd.minCost[i] = INF
		}

		que.Push(pair{0, start})
		pd.minCost[start] = 0
		for que.Len() > 0 {
			p := que.Pop()
			if pd.minCost[p.second] < p.first {
				continue
			}

			for i := 0; i < len(pd.graph[p.second]); i++ {
				e := pd.graph[p.second][i]
				nextCost := pd.minCost[p.second] + e.cost + pd.potential[p.second] - pd.potential[e.to]
				if e.cap > 0 && pd.minCost[e.to] > nextCost {
					pd.minCost[e.to] = nextCost
					pd.prevv[e.to] = p.second
					pd.preve[e.to] = int32(i)
					que.Push(pair{pd.minCost[e.to], e.to})
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

type edgeInfo struct {
	from, to  int32
	flow, cap int
}

// 输出原图的边信息.
func (p *PrimalDual) GetEdges() (res []edgeInfo) {
	for i := int32(0); i < int32(len(p.graph)); i++ {
		for _, e := range p.graph[i] {
			if e.isRev {
				continue
			}
			revEdge := p.graph[e.to][e.rev]
			res = append(res, edgeInfo{from: i, to: e.to, flow: revEdge.cap, cap: e.cap + revEdge.cap})
		}
	}
	return
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
