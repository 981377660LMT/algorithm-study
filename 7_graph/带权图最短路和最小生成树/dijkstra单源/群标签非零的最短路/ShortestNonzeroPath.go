// !群标签非零的最短路
// https://ygussany.hatenablog.com/entry/2019/12/04/000000
// https://ei1333.github.io/library/graph/shortest-path/shortest-nonzero-path.hpp
// 无向图的边上携带群元素
// 求出起点到各个终点的最短路,满足使得路径上的群元素op之和不为幺元
// O(mlogn)

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://yukicoder.me/problems/no/1602
	// 给定n个点m条边的无向图,每条边上有一个xor标签
	// 求各个点0-n-2到n-1的最短路,满足路径上的xor标签异或和不为0
	// 如果不存在这样的路径,输出-1
	// n<=1e5 m<=2e5 k<=30
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)
	g := NewShortestNonzeroPath(n)
	for i := 0; i < m; i++ {
		var u, v, cost int
		var s string
		fmt.Fscan(in, &u, &v, &cost, &s)
		u, v = u-1, v-1
		label := 0
		for j := 0; j < k; j++ {
			if s[j] == '1' {
				label |= 1 << j
			}
		}
		g.AddUndirectedEdge(u, v, cost, label)
	}

	res := g.Build(n - 1)
	for i := 0; i <= n-2; i++ {
		if res[i] == INF {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, res[i])
		}
	}

}

const INF int = 1e18

type Label = int // xor标签

func (*ShortestNonzeroPath) e() Label                      { return 0 }
func (*ShortestNonzeroPath) op(label1, label2 Label) Label { return label1 ^ label2 }

type ShortestNonzeroPath struct {
	g  [][]edge
	uf []int
}

type edge = struct {
	to    int
	cost  int
	label Label
}

type SP struct {
	dist   []int
	depth  []int
	parent []int
	label  []Label
}

func NewShortestNonzeroPath(n int) *ShortestNonzeroPath {
	return &ShortestNonzeroPath{g: make([][]edge, n)}
}

func (s *ShortestNonzeroPath) AddUndirectedEdge(u, v, cost int, label Label) {
	s.AddDirectedEdge(u, v, cost, label)
	s.AddDirectedEdge(v, u, cost, label)
}

func (s *ShortestNonzeroPath) AddDirectedEdge(u, v, cost int, label Label) {
	s.g[u] = append(s.g[u], edge{v, cost, label})
}

func (snp *ShortestNonzeroPath) Build(start int) (dist []int) {
	n := len(snp.g)
	sp := snp.dijkstra(start)
	snp.uf = make([]int, n)
	for i := range snp.uf {
		snp.uf[i] = -1
	}

	type tuple = [3]int
	pq := NewHeap(func(a, b H) int {
		return a.(tuple)[0] - b.(tuple)[0]
	}, nil)
	for u := 0; u < n; u++ {
		if sp.dist[u] != INF {
			for i := range snp.g[u] {
				e := snp.g[u][i]
				if u < e.to && snp.op(sp.label[u], e.label) != sp.label[e.to] {
					pq.Push(tuple{sp.dist[u] + sp.dist[e.to] + e.cost, u, i})
				}
			}
		}
	}

	dist = make([]int, n)
	for i := range dist {
		dist[i] = INF
	}
	bs := []int{}
	for pq.Len() > 0 {
		tmp := pq.Pop().(tuple)
		cost, u0, i := tmp[0], tmp[1], tmp[2]
		v0 := snp.g[u0][i].to
		u, v := snp.findUf(u0), snp.findUf(v0)
		for u != v {
			if sp.depth[u] > sp.depth[v] {
				bs = append(bs, u)
				u = snp.findUf(sp.parent[u])
			} else {
				bs = append(bs, v)
				v = snp.findUf(sp.parent[v])
			}
		}
		for _, x := range bs {
			snp.uniteUf(u, x)
			dist[x] = cost - sp.dist[x]
			for j := range snp.g[x] {
				e := snp.g[x][j]
				if snp.op(sp.label[x], e.label) == sp.label[e.to] {
					pq.Push(tuple{dist[x] + sp.dist[e.to] + e.cost, x, j})
				}
			}
		}
		bs = bs[:0]
	}

	for i := 0; i < n; i++ {
		if sp.label[i] != snp.e() && sp.dist[i] < dist[i] {
			dist[i] = sp.dist[i]
		}
	}

	return
}

func (snp *ShortestNonzeroPath) dijkstra(s int) *SP {
	n := len(snp.g)
	type pair = [2]int
	dist := make([]int, n)
	for i := range dist {
		dist[i] = INF
	}
	depth, parent := make([]int, n), make([]int, n)
	for i := range parent {
		parent[i] = -1
		depth[i] = -1
	}
	label := make([]Label, n)
	for i := range label {
		label[i] = snp.e()
	}

	pq := NewHeap(func(a, b H) int {
		return a.([2]int)[0] - b.([2]int)[0]
	}, nil)
	pq.Push(pair{0, s})
	dist[s] = 0
	depth[s] = 0

	for pq.Len() > 0 {
		p := pq.Pop().([2]int)
		cost, cur := p[0], p[1]
		if dist[cur] < cost {
			continue
		}
		for _, e := range snp.g[cur] {
			to, nextCost := e.to, cost+e.cost
			if dist[to] > nextCost {
				dist[to] = nextCost
				parent[to] = cur
				depth[to] = depth[cur] + 1
				label[to] = snp.op(label[cur], e.label)
				pq.Push(pair{nextCost, to})
			}
		}
	}

	return &SP{dist, depth, parent, label}
}

func (s *ShortestNonzeroPath) findUf(k int) int {
	if s.uf[k] == -1 {
		return k
	}
	s.uf[k] = s.findUf(s.uf[k])
	return s.uf[k]
}

func (s *ShortestNonzeroPath) uniteUf(r, c int) {
	s.uf[c] = r
}

type H = interface{}

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
