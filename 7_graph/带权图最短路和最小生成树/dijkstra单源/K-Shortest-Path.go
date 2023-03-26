// !求从s到t的k条最短路的长度,不经过重复点
// https://ei1333.github.io/library/test/verify/yukicoder-1069.test.cpp
// https://maspypy.github.io/library/graph/shortest_path/K_shortest_path.hpp
// O(K*V*E*LogV)
// n<=2e3 m<=2e3 k<=10

package main

import "fmt"

func main() {
	res := KthShortestPath(4,
		[]Edge{{0, 1, 1, 0}, {0, 2, 2, 1}, {1, 2, 1, 2}, {1, 3, 3, 3}, {2, 3, 1, 4}},
		0, 3, 2,
	)
	fmt.Println(res) // [{3 [0 2 3] [1 4]} {3 [0 1 2 3] [0 2 4]}]
}

const INF int = 1e18

type PathInfo struct {
	cost   int
	vs, es []int
}

func KthShortestPath(n int, edges []Edge, start, target, k int) (res []PathInfo) {
	graph := make([][]Edge, n)
	for _, e := range edges {
		graph[e.from] = append(graph[e.from], e)
	}

	type nodeItem struct{ es, ngEdges []int }
	type pathItem struct {
		es, ngEs []int
		cost, n  int
	}
	var nodes []nodeItem
	var paths []pathItem

	nodes = append(nodes, nodeItem{[]int{}, []int{}})
	var dist, parent []int
	var ngV, ngE []bool

	for len(res) < k {
		for _, item := range nodes {
			es, ngEs := item.es, item.ngEdges
			dist = make([]int, n)
			parent = make([]int, n)
			ngV = make([]bool, n)
			ngE = make([]bool, len(edges))
			for i := range dist {
				dist[i] = INF
				parent[i] = -1
			}
			prefCost := 0
			for _, x := range es {
				prefCost += edges[x].weight
			}

			for _, x := range es {
				ngV[edges[x].from] = true
				ngE[x] = true
			}
			for _, x := range ngEs {
				ngE[x] = true
			}

			type pqItem struct{ node, dist int }
			que := nhp(func(a, b H) int { return a.(pqItem).dist - b.(pqItem).dist }, nil)
			add := func(v, d, p int) {
				if dist[v] > d {
					dist[v] = d
					que.Push(pqItem{v, d})
					parent[v] = p
				}
			}
			s0 := start
			if len(es) > 0 {
				s0 = edges[es[len(es)-1]].to
			}

			add(s0, prefCost, -1)
			for que.Len() > 0 {
				popped := que.Pop().(pqItem)
				dv, v := popped.dist, popped.node
				if dv != dist[v] {
					continue
				}
				if v == target {
					break
				}
				for _, e := range graph[v] {
					if ngE[e.id] || ngV[e.to] {
						continue
					}
					add(e.to, dv+e.weight, e.id)
				}
			}

			if parent[target] == -1 {
				continue
			}
			addE := []int{}
			v := target
			for v != s0 {
				addE = append(addE, parent[v])
				v = edges[parent[v]].from
			}
			for i := 0; i < len(addE)/2; i++ {
				addE[i], addE[len(addE)-1-i] = addE[len(addE)-1-i], addE[i]
			}
			n := len(es)
			es = append(es, addE...)
			paths = append(paths, pathItem{es, ngEs, dist[target], n})
		}

		if len(paths) == 0 {
			break
		}

		best := [2]int{-1, INF}
		for i := 0; i < len(paths); i++ {
			cost := paths[i].cost
			if cost < best[1] {
				best[0], best[1] = i, cost
			}
		}
		idx := best[0]
		paths[idx], paths[len(paths)-1] = paths[len(paths)-1], paths[idx]
		path := paths[len(paths)-1]
		paths = paths[:len(paths)-1]
		cost, es, ngEs, n := path.cost, path.es, path.ngEs, path.n
		vs := []int{start}
		for _, x := range es {
			vs = append(vs, edges[x].to)
		}
		res = append(res, PathInfo{cost, vs, es})
		nodes = nodes[:0]
		for k := n; k < len(es); k++ {
			newEs := append([]int{}, es[:k]...)
			newNg := append([]int{}, ngEs...)
			newNg = append(newNg, es[k])
			nodes = append(nodes, nodeItem{newEs, newNg})
		}
	}
	return
}

type Edge struct{ from, to, weight, id int }

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
