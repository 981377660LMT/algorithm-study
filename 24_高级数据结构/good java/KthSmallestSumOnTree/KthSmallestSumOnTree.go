package main

import (
	"fmt"
	"sort"
)

func main() {
	//    0
	//   / \
	//  1   2
	//     / \
	//    3   4
	//       / \
	//      5   6

	n := int32(7)
	edges := [][3]int32{{0, 1, 5}, {0, 2, 4}, {2, 3, 3}, {2, 4, 2}, {4, 5, 1}, {4, 6, 0}}

	rawTree := make([][][2]int32, n)
	for _, edge := range edges {
		u, v, w := edge[0], edge[1], edge[2]
		rawTree[u] = append(rawTree[u], [2]int32{v, w})
		rawTree[v] = append(rawTree[v], [2]int32{u, w})
	}

	root := int32(0)
	rootedTree := make([][][2]int32, n)

	queue := make([]int32, 0, n)
	queue = append(queue, root)
	visited := make([]bool, n)
	visited[root] = true

	vertexes := make([]*vertex, n)

	// 转有根树和构建iterator在一次bfs中完成
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, e := range rawTree[cur] {
			next := e[0]
			if !visited[next] {
				visited[next] = true
				queue = append(queue, next)
				rootedTree[cur] = append(rootedTree[cur], e)
			}
		}

		children := rootedTree[cur]
		sort.Slice(children, func(i, j int) bool { return children[i][1] < children[j][1] })
		vertexes[cur] = &vertex{Node: cur}
		vertexes[cur].Children = func() *edgeIterator {
			ptr := 0
			return &edgeIterator{
				Next: func() *edge {
					if ptr < len(children) {
						e := children[ptr]
						ptr++
						return &edge{To: vertexes[e[0]], Weight: int(e[1])}
					}
					return nil
				},
				HasNext: func() bool {
					return ptr < len(children)
				},
			}
		}
	}

	res := KthSmallestSumOnTree(vertexes[root], 6)
	for _, state := range res {
		fmt.Println(state)
	}
}

// 给定一颗树，每条边上有一个非负权重，要求从根出发，找到路径权重最小的k条路径
func KthSmallestSumOnTree(root *vertex, k int32) []*state {
	res := make([]*state, 0, k)
	pq := NewHeap(func(a, b *state) bool { return a.Sum < b.Sum }, nil)
	pq.Push(newState(root, nil, 0))
	for pq.Len() > 0 && int32(len(res)) < k {
		state := pq.Pop()
		res = append(res, state)
		// child
		if state.iterator.HasNext() {
			e := state.iterator.Next()
			pq.Push(newState(e.To, state, state.Sum+int(e.Weight)))
		}
		// sibling
		if state.ParentState != nil && state.ParentState.iterator.HasNext() {
			e := state.ParentState.iterator.Next()
			pq.Push(newState(e.To, state.ParentState, state.ParentState.Sum+int(e.Weight)))
		}
	}
	return res
}

// 遍历的一个状态.
type state struct {
	Cur         *vertex
	ParentState *state
	Sum         int
	iterator    *edgeIterator
}

func newState(cur *vertex, parent *state, sum int) *state {
	return &state{Cur: cur, ParentState: parent, Sum: sum, iterator: cur.Children()}
}

func newStateWithIterator(cur *vertex, parent *state, sum int, iterator *edgeIterator) *state {
	return &state{Cur: cur, ParentState: parent, Sum: sum, iterator: iterator}
}

func (s *state) String() string {
	return fmt.Sprintf("vertex: %d, sum: %d", s.Cur.Node, s.Sum)
}

type edge struct {
	To     *vertex
	Weight int
}

type vertex struct {
	Node     int32
	Children func() *edgeIterator
	// CustomInfo
}

type edgeIterator struct {
	Next    func() *edge
	HasNext func() bool
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
