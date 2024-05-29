// DynamicShortestPath
// 动态最短路，带修改的有向图最短路/带修改的最短路问题

package main

const INF int = 2e18
const INF32 int32 = 1e9 + 10

type DirectedModifiableShortestPath struct {
	curTag   int32
	seg      *segment
	nodes    []*node
	edges    []*edge
	src      *node
	dst      *node
	prepared bool
}

func NewDirectedModifiableShortestPath(n, m, s, t int32) *DirectedModifiableShortestPath {
	nodes := make([]*node, n)
	for i := int32(0); i < n; i++ {
		nodes[i] = newNode()
		nodes[i].id = i
		nodes[i].distToSrc = INF
		nodes[i].distToDst = INF
	}
	edges := make([]*edge, 0, m)
	src := nodes[s]
	dst := nodes[t]
	return &DirectedModifiableShortestPath{nodes: nodes, edges: edges, src: src, dst: dst}
}

func (sp *DirectedModifiableShortestPath) AddEdge(u, v int32, cost int) {
	if sp.prepared {
		panic("can't add edge after prepare")
	}
	e := newEdge()
	e.a = sp.nodes[u]
	e.b = sp.nodes[v]
	e.w = cost
	e.a.adj = append(e.a.adj, e)
	e.b.adj = append(e.b.adj, e)
	sp.edges = append(sp.edges, e)
}

func (sp *DirectedModifiableShortestPath) MinDist() int {
	sp.prepare()
	return sp.dst.distToSrc
}

// O(1)
func (sp *DirectedModifiableShortestPath) QueryOnAddEdge(a, b int32, dist int) int {
	sp.prepare()
	res := sp.dst.distToSrc
	res = min(res, sp.nodes[a].distToSrc+dist+sp.nodes[b].distToDst)
	return res
}

// O(log n)
func (sp *DirectedModifiableShortestPath) QueryOnModifyEdge(i int32, cost int) int {
	sp.prepare()
	e := sp.edges[i]
	w := cost
	res := e.a.distToSrc + e.b.distToDst + w
	if e.tag == -1 {
		res = min(res, sp.dst.distToSrc)
	} else {
		val := sp.seg.Query(e.tag, e.tag, 1, sp.curTag)
		res = min(res, val)
	}
	return res
}

// O(log n)
func (sp *DirectedModifiableShortestPath) QueryOnDeleteEdge(i int32) int {
	return sp.QueryOnModifyEdge(i, INF)
}

func (sp *DirectedModifiableShortestPath) update(a, b *node, w int) {
	dist := a.distToSrc + b.distToDst + w
	l := sp.prev(a)
	r := sp.post(b)
	sp.seg.Update(l+1, r-1, 1, sp.curTag, dist)
}

func (sp *DirectedModifiableShortestPath) post(root *node) int32 {
	if root.r == -INF32 {
		root.r = sp.curTag + 1
		for _, e := range root.adj {
			if e.a != root {
				continue
			}
			node := e.Other(root)
			if node.distToDst+e.w == root.distToDst {
				if e.tag != -1 {
					root.r = e.tag
				} else {
					root.r = sp.post(node)
				}
				break
			}
		}
	}
	return root.r
}

func (sp *DirectedModifiableShortestPath) prev(root *node) int32 {
	if root.l == -INF32 {
		root.l = -1
		for _, e := range root.adj {
			if e.b != root {
				continue
			}
			node := e.Other(root)
			if node.distToSrc+e.w == root.distToSrc {
				if e.tag != -1 {
					root.l = e.tag
				} else {
					root.l = sp.prev(node)
				}
				break
			}
		}
	}
	return root.l
}

func (sp *DirectedModifiableShortestPath) prepare() {
	if sp.prepared {
		return
	}
	sp.prepared = true

	pq := newErasableHeapGeneric(func(a, b *node) bool {
		if a.distToSrc == b.distToSrc {
			return a.id < b.id
		}
		return a.distToSrc < b.distToSrc
	})

	sp.src.distToSrc = 0
	pq.Push(sp.src)
	for pq.Len() > 0 {
		head := pq.Pop()
		for _, e := range head.adj {
			if e.a != head {
				continue
			}
			node := e.Other(head)
			if tmp := head.distToSrc + e.w; node.distToSrc > tmp {
				pq.Erase(node)
				node.distToSrc = tmp
				pq.Push(node)
			}
		}
	}

	sp.dst.distToDst = 0
	pq = newErasableHeapGeneric(func(a, b *node) bool {
		if a.distToDst == b.distToDst {
			return a.id < b.id
		}
		return a.distToDst < b.distToDst
	})
	pq.Push(sp.dst)
	for pq.Len() > 0 {
		head := pq.Pop()
		for _, e := range head.adj {
			if e.b != head {
				continue
			}
			node := e.Other(head)
			if tmp := head.distToDst + e.w; node.distToDst > tmp {
				pq.Erase(node)
				node.distToDst = tmp
				pq.Push(node)
			}
		}
	}

	for trace := sp.src; trace != sp.dst; {
		var next *edge
		for _, e := range trace.adj {
			if e.a != trace {
				continue
			}
			node := e.Other(trace)
			if node.distToDst+e.w == trace.distToDst {
				next = e
				break
			}
		}
		sp.curTag++
		next.tag = sp.curTag
		trace = next.Other(trace)
	}

	sp.seg = newSegment(1, sp.curTag)
	for _, e := range sp.edges {
		if e.tag != -1 {
			continue
		}
		sp.update(e.a, e.b, e.w)
	}
}

type edge struct {
	a, b *node
	w    int
	tag  int32
}

func newEdge() *edge {
	return &edge{tag: -1}
}

func (e *edge) Other(x *node) *node {
	if x == e.a {
		return e.b
	}
	return e.a
}

type node struct {
	adj                  []*edge
	distToSrc, distToDst int
	l, r                 int32
	id                   int32
}

func newNode() *node {
	return &node{l: -INF32, r: -INF32}
}

type segment struct {
	left, right *segment
	min         int
}

func newSegment(l, r int32) *segment {
	if l < r {
		m := (l + r) >> 1
		res := &segment{left: newSegment(l, m), right: newSegment(m+1, r), min: INF}
		res.pushUp()
		return res
	}
	return &segment{min: INF}
}

func (s *segment) modify(x int) {
	s.min = min(s.min, x)
}

func (s *segment) pushUp() {}

func (s *segment) pushDown() {
	s.left.modify(s.min)
	s.right.modify(s.min)
	s.min = INF
}

func (s *segment) covered(ll, rr, l, r int32) bool {
	return ll <= l && rr >= r
}

func (s *segment) noIntersection(ll, rr, l, r int32) bool {
	return ll > r || rr < l
}

func (s *segment) Update(ll, rr, l, r int32, x int) {
	if s.noIntersection(ll, rr, l, r) {
		return
	}
	if s.covered(ll, rr, l, r) {
		s.modify(x)
		return
	}
	s.pushDown()
	m := (l + r) >> 1
	s.left.Update(ll, rr, l, m, x)
	s.right.Update(ll, rr, m+1, r, x)
	s.pushUp()
}

func (s *segment) Query(ll, rr, l, r int32) int {
	if s.noIntersection(ll, rr, l, r) {
		return INF
	}
	if s.covered(ll, rr, l, r) {
		return s.min
	}
	s.pushDown()
	m := (l + r) >> 1
	return min(s.left.Query(ll, rr, l, m), s.right.Query(ll, rr, m+1, r))
}

type erasableHeapGeneric[H comparable] struct {
	data   *heapGeneric[H]
	erased *heapGeneric[H]
	size   int32
}

func newErasableHeapGeneric[H comparable](less func(a, b H) bool, nums ...H) *erasableHeapGeneric[H] {
	return &erasableHeapGeneric[H]{newHeapGeneric(less, nums...), newHeapGeneric(less), int32(len(nums))}
}

// 从堆中删除一个元素,要保证堆中存在该元素.
func (h *erasableHeapGeneric[H]) Erase(value H) {
	h.erased.Push(value)
	h.normalize()
	h.size--
}

func (h *erasableHeapGeneric[H]) Push(value H) {
	h.data.Push(value)
	h.normalize()
	h.size++
}

func (h *erasableHeapGeneric[H]) Pop() (value H) {
	value = h.data.Pop()
	h.normalize()
	h.size--
	return
}

func (h *erasableHeapGeneric[H]) Peek() (value H) {
	value = h.data.Top()
	return
}

func (h *erasableHeapGeneric[H]) Len() int32 {
	return h.size
}

func (h *erasableHeapGeneric[H]) Clear() {
	h.data.Clear()
	h.erased.Clear()
	h.size = 0
}

func (h *erasableHeapGeneric[H]) normalize() {
	for h.data.Len() > 0 && h.erased.Len() > 0 && h.data.Top() == h.erased.Top() {
		h.data.Pop()
		h.erased.Pop()
	}
}

type heapGeneric[H comparable] struct {
	data []H
	less func(a, b H) bool
}

func newHeapGeneric[H comparable](less func(a, b H) bool, nums ...H) *heapGeneric[H] {
	nums = append(nums[:0:0], nums...)
	heap := &heapGeneric[H]{less: less, data: nums}
	if len(nums) > 1 {
		heap.heapify()
	}
	return heap
}

func (h *heapGeneric[H]) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *heapGeneric[H]) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}

	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *heapGeneric[H]) Top() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}
	value = h.data[0]
	return
}

func (h *heapGeneric[H]) Len() int32 { return int32(len(h.data)) }

func (h *heapGeneric[H]) Clear() {
	h.data = h.data[:0]
}

func (h *heapGeneric[H]) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *heapGeneric[H]) pushUp(root int32) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *heapGeneric[H]) pushDown(root int32) {
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
