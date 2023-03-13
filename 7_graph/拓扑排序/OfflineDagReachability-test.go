// Offline Dag Reachability(DAGの到達可能性クエリ)
// https://ei1333.github.io/library/graph/others/offline-dag-reachability.hpp

// 如果图上有环,先SCC分解成DAG
// 然后64个查询一组批处理

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=0275
	// 给定n给顶点m边的有向图,当前地铁票的地铁进站口start和出站口end.
	// 对于每个查询(a,b)判断从现在的票是否能满足从a到b.
	// 限制条件:
	// 1.a和b在start到end的最短路上 -> 检查每条边(u,v)是否满足dist1[u]+w+dist2[v]=dist1[end] 或者 dist2[u]+w+dist1[v]=dist1[end]
	// 2.a到b的路径必须走最短路 -> 无环
	// n,m<=1e5 q<=4e4

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	graph := make([][]Edge, n)
	edges := make([][3]int, 0, m)
	for i := 0; i < m; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		u--
		v--
		graph[u] = append(graph[u], Edge{v, w})
		graph[v] = append(graph[v], Edge{u, w})
		edges = append(edges, [3]int{u, v, w})
	}

	var start, end int
	fmt.Fscan(in, &start, &end)
	start--
	end--
	dist1, dist2 := Dijkstra(n, graph, start), Dijkstra(n, graph, end)

	dag := make([][]int, n)
	d := dist1[end]
	for _, e := range edges {
		u, v := e[0], e[1]
		if dist1[u]+e[2]+dist2[v] == d {
			dag[u] = append(dag[u], v)
		}
		if dist2[u]+e[2]+dist1[v] == d {
			dag[v] = append(dag[v], u)
		}
	}

	var q int
	fmt.Fscan(in, &q)
	queries := make([][2]int, 0, q)
	for i := 0; i < q; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		queries = append(queries, [2]int{a, b})
	}

	res := offlineDagReachability(dag, queries)
	for _, yn := range res {
		if yn {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}

// 在有向无环图graph上,对每个查询(u,v)判断从u是否能到达v.
//  O((E+V)*Q/64)
func offlineDagReachability(dag [][]int, queries [][2]int) []bool {
	n, q := len(dag), len(queries)
	order := topoSort(dag)
	res := make([]bool, q)
	for i := 0; i < q; i += 64 {
		upper := min(q, i+64)
		dp := make([]uint64, n)
		for k := i; k < upper; k++ {
			dp[queries[k][0]] |= 1 << uint(k-i)
		}

		for _, cur := range order {
			for _, next := range dag[cur] {
				dp[next] |= dp[cur]
			}
		}

		for k := i; k < upper; k++ {
			res[k] = dp[queries[k][1]]&(1<<uint(k-i)) > 0
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

func topoSort(dag [][]int) []int {
	n := len(dag)
	deg := make([]int, n)
	for i := 0; i < n; i++ {
		for _, v := range dag[i] {
			deg[v]++
		}
	}

	queue := []int{}
	for i := 0; i < n; i++ {
		if deg[i] == 0 {
			queue = append(queue, i)
		}
	}

	order := []int{}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		order = append(order, cur)
		for _, next := range dag[cur] {
			deg[next]--
			if deg[next] == 0 {
				queue = append(queue, next)
			}
		}
	}

	return order
}

const INF int = 1e18

type Edge struct{ to, weight int }

func Dijkstra(n int, adjList [][]Edge, start int) (dist []int) {
	dist = make([]int, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0

	pq := nhp(func(a, b H) int {
		return a.dist - b.dist
	}, []H{{start, 0}})

	for pq.Len() > 0 {
		curNode := pq.Pop()
		cur, curDist := curNode.node, curNode.dist
		if curDist > dist[cur] {
			continue
		}

		for _, edge := range adjList[cur] {
			next, weight := edge.to, edge.weight
			if cand := curDist + weight; cand < dist[next] {
				dist[next] = cand
				pq.Push(H{next, cand})
			}
		}
	}

	return
}

type H = struct{ node, dist int }

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
