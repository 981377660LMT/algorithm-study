// Johnson全源最短路
// 给定一个包含 n 个结点和 m 条带权边的有向图
// 求所有点对间的最短路径长度，一条路径的长度定义为这条路径上所有边的权值和。
// 边权可能为负，且图中可能存在重边和自环。
// n<=3e3 m<=6e3

package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

func main() {
	// https://www.luogu.com.cn/problem/P5905
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([][]int, 0, m)
	for i := 0; i < m; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		u, v = u-1, v-1
		edges = append(edges, []int{u, v, w})
	}

	dist := shortestPathJohnson(n, edges, true)
	if dist == nil { // 有负环
		fmt.Fprintln(out, -1)
		return
	}

	for i := 0; i < n; i++ {
		sum := 0
		for j := 0; j < n; j++ {
			if dist[i][j] == INF {
				dist[i][j] = 1e9
			}
			sum += dist[i][j] * (j + 1)
		}
		fmt.Fprintln(out, sum)
	}

}

const INF int = 1e18

// 任意两点最短路 Johnson O(nmlogm)
//  若有负环返回 nil
//  https://en.wikipedia.org/wiki/Johnson%27s_algorithm
//  https://oi-wiki.org/graph/shortest-path/#johnson
//  模板题 https://www.luogu.com.cn/problem/P5905
//  https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/graph.go
func shortestPathJohnson(n int, edges [][]int, directed bool) [][]int {
	type neighbor struct{ to, weight int }
	g := make([][]neighbor, n+1)
	for _, e := range edges {
		u, v, w := e[0], e[1], e[2]
		u, v = u+1, v+1
		g[u] = append(g[u], neighbor{v, w})
		if !directed {
			g[v] = append(g[v], neighbor{u, w})
		}
	}

	// 建虚拟节点 0 并且往其他的点都连一条边权为 0 的边
	for v := 1; v <= n; v++ {
		g[0] = append(g[0], neighbor{v, 0})
		if !directed {
			g[v] = append(g[v], neighbor{})
		}
	}

	spfa := func(s int) []int {
		h := make([]int, n+1)
		for i := range h {
			h[i] = INF
		}
		h[s] = 0
		inQ := make([]bool, n+1)
		inQ[s] = true
		relaxedCnt := make([]int, n+1)
		q := make([]int, 1, n+1)
		for len(q) > 0 {
			v := q[0]
			q = q[1:]
			inQ[v] = false
			for _, e := range g[v] {
				w := e.to
				if newH := h[v] + e.weight; newH < h[w] {
					h[w] = newH
					relaxedCnt[w] = relaxedCnt[v] + 1
					if relaxedCnt[w] > n {
						return nil
					}
					if !inQ[w] {
						q = append(q, w)
						inQ[w] = true
					}
				}
			}
		}

		return h
	}

	h := spfa(0)
	if h == nil {
		return nil
	}

	// 求新的边权
	for v := 1; v <= n; v++ {
		for i, e := range g[v] {
			g[v][i].weight += h[v] - h[e.to]
		}
	}

	dijkstra := func(st int) []int {
		dist := make([]int, n+1)
		for i := range dist {
			dist[i] = INF
		}
		dist[st] = 0
		q := hp{{st, 0}}
		for len(q) > 0 {
			p := heap.Pop(&q).(pair)
			v := p.v
			if dist[v] < p.d {
				continue
			}
			for _, e := range g[v] {
				w := e.to
				if newD := dist[v] + e.weight; newD < dist[w] {
					dist[w] = newD
					heap.Push(&q, pair{w, newD})
				}
			}
		}
		return dist
	}

	// 以每个点为源点跑一遍 Dijkstra
	dist := make([][]int, n+1)
	for st := 1; st <= n; st++ {
		dist[st] = dijkstra(st)
		for end, d := range dist[st] {
			if d < INF {
				dist[st][end] -= h[st] - h[end]
			}
		}
	}

	for i := 1; i <= n; i++ {
		dist[i] = dist[i][1:]
	}
	return dist[1:]
}

type pair struct{ v, d int }
type hp []pair

func (h hp) Len() int              { return len(h) }
func (h hp) Less(i, j int) bool    { return h[i].d < h[j].d }
func (h hp) Swap(i, j int)         { h[i], h[j] = h[j], h[i] }
func (h *hp) Push(v interface{})   { *h = append(*h, v.(pair)) }
func (h *hp) Pop() (v interface{}) { a := *h; *h, v = a[:len(a)-1], a[len(a)-1]; return }
