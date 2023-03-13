// https://ei1333.github.io/library/graph/shortest-path/k-shortest-walk.hpp
// 求从s到t的k条最短路的长度
// O((V+E)*LogV+K*LogK)
// n,m,k<=3e5

// !还有一个求从s到t的k条最短路的长度,不经过重复点 的版本
// https://ei1333.github.io/library/test/verify/yukicoder-1069.test.cpp
// O(K*V*E*LogV)
// n<=2e3 m<=2e3 k<=10

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const INF int = 1e18

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, start, end, k int
	fmt.Fscan(in, &n, &m, &start, &end, &k)
	edges := make([]Edge, m)
	for i := 0; i < m; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		edges[i] = Edge{u, v, w, i}
	}
	res := kShortestWalk(n, edges, start, end, k)
	for i := 0; i < len(res); i++ {
		fmt.Fprintln(out, res[i])
	}
	for i := 0; i < k-len(res); i++ {
		fmt.Fprintln(out, -1)
	}

}

type Edge struct{ from, to, weight, id int }

// 按照升序返回k条最短路的路径长度
//  不满k条时全部返回
func kShortestWalk(n int, edges []Edge, start, target, k int) (res []int) {
	g := make([][]Edge, n)
	rg := make([][]Edge, n)
	for _, e := range edges {
		g[e.from] = append(g[e.from], e)
		rg[e.to] = append(rg[e.to], Edge{e.to, e.from, e.weight, e.id})
	}

	dist, from := Dijkstra(n, rg, target)
	if from[start] == -1 {
		return
	}

	ch := make([][]int, n)
	for i := 0; i < n; i++ {
		if from[i] >= 0 {
			ch[from[i]] = append(ch[from[i]], i)
		}
	}

	pHeap := NewPersistentLeftistHeap() // 持久化小根堆
	h := make([]*heap, n)

	{
		que := []int{target}
		for len(que) > 0 {
			idx := que[0]
			que = que[1:]
			if from[idx] >= 0 {
				h[idx] = pHeap.Meld(h[idx], h[from[idx]])
			}
			used := true
			for _, e := range g[idx] {
				if e.to != target && from[e.to] == -1 {
					continue
				}
				if used && from[idx] == e.to && dist[e.to]+e.weight == dist[idx] {
					used = false
					continue
				}
				h[idx] = pHeap.Push(h[idx], e.weight-dist[idx]+dist[e.to], e.to)
			}
			for _, to := range ch[idx] {
				que = append(que, to)
			}
		}
	}

	type pair = struct {
		first  int
		second *heap
	}
	pq := nhp(func(v1, v2 H) int { return v1.(pair).first - v2.(pair).first }, nil)
	var root *heap
	root = pHeap.Push(root, dist[start], start)
	pq.Push(pair{dist[start], root})

	for pq.Len() > 0 {
		item := pq.Pop().(pair)
		cost, cur := item.first, item.second
		res = append(res, cost)
		if len(res) == k {
			break
		}
		if cur.left != nil {
			pq.Push(pair{cost + cur.left.Value - cur.Value, cur.left})
		}
		if cur.right != nil {
			pq.Push(pair{cost + cur.right.Value - cur.Value, cur.right})
		}
		if h[cur.Id] != nil {
			pq.Push(pair{cost + h[cur.Id].Value, h[cur.Id]})
		}
	}
	return

}

type V = int

type leftistHeap struct {
}

type heap struct {
	Value       V
	Id          int
	height      int // 维持平衡
	left, right *heap
}

func NewPersistentLeftistHeap() *leftistHeap {
	return &leftistHeap{}
}

func (lh *leftistHeap) Alloc(value V, id int) *heap {
	res := &heap{Value: value, Id: id, height: 1}
	return res
}

func (lh *leftistHeap) Build(nums []V) []*heap {
	res := make([]*heap, len(nums))
	for i, num := range nums {
		res[i] = lh.Alloc(num, i)
	}
	return res
}

func (lh *leftistHeap) Push(heap *heap, value V, id int) *heap {
	return lh.Meld(heap, lh.Alloc(value, id))
}

func (lh *leftistHeap) Pop(heap *heap) *heap {
	return lh.Meld(heap.left, heap.right)
}

func (lh *leftistHeap) Top(heap *heap) V {
	return heap.Value
}

// 合并两个堆,返回合并后的堆.
func (lh *leftistHeap) Meld(heap1, heap2 *heap) *heap {
	if heap1 == nil {
		return heap2
	}
	if heap2 == nil {
		return heap1
	}
	if heap2.Value < heap1.Value {
		heap1, heap2 = heap2, heap1
	}
	heap1 = lh.clone(heap1)
	heap1.right = lh.Meld(heap1.right, heap2)
	if heap1.left == nil || heap1.left.height < heap1.right.height {
		heap1.left, heap1.right = heap1.right, heap1.left
	}
	heap1.height = 1
	if heap1.right != nil {
		heap1.height += heap1.right.height
	}
	return heap1
}

func (h *heap) String() string {
	var sb []string
	var dfs func(h *heap)
	dfs = func(h *heap) {
		if h == nil {
			return
		}
		sb = append(sb, fmt.Sprintf("%d", h.Value))
		dfs(h.left)
		dfs(h.right)
	}
	dfs(h)
	return strings.Join(sb, " ")
}

// 持久化,拷贝一份结点.
func (lh *leftistHeap) clone(h *heap) *heap {
	if h == nil {
		return h
	}
	res := &heap{height: h.height, Value: h.Value, Id: h.Id, left: h.left, right: h.right}
	return res
}

func Dijkstra(n int, adjList [][]Edge, start int) (dist, preV []int) {
	type pqItem struct{ node, dist int }
	dist = make([]int, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0
	preV = make([]int, n)
	for i := range preV {
		preV[i] = -1
	}

	pq := nhp(func(a, b H) int {
		return a.(pqItem).dist - b.(pqItem).dist
	}, nil)
	pq.Push(pqItem{start, 0})

	for pq.Len() > 0 {
		curNode := pq.Pop().(pqItem)
		cur, curDist := curNode.node, curNode.dist
		if curDist > dist[cur] {
			continue
		}

		for _, edge := range adjList[cur] {
			next, weight := edge.to, edge.weight
			if cand := curDist + weight; cand < dist[next] {
				dist[next] = cand
				preV[next] = cur
				pq.Push(pqItem{next, cand})
			}
		}
	}

	return
}

type H = interface{}

// Should return a number:
//    negative , if a < b
//    zero     , if a == b
//    positive , if a > b
type Comparator func(a, b H) int

func nhp(comparator Comparator, nums []H) *Heap {
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
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
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
