package mincostmaxflow

import "io"

// 最小费用流 MCFP
// 最小费用最大流 MCMF（即满流时的费用）
// 将 Edmonds-Karp 中的 BFS 改成 SPFA O(fnm) 或 Dijkstra O(fmlogn)
// 要求初始网络中无负权圈
// 性能对比（洛谷 P3381，由于数据不强所以 SPFA 很快）：SPFA 1.05s(max 365ms)   Dijkstra 1.91s(max 688ms)
// https://en.wikipedia.org/wiki/Edmonds%E2%80%93Karp_algorithm
// https://oi-wiki.org/graph/flow/min-cost/
// https://cp-algorithms.com/graph/min_cost_flow.html
// 最小费用流的不完全算法博物馆 https://www.luogu.com.cn/blog/Atalod/zui-xiao-fei-yong-liu-di-fou-wan-quan-suan-fa-bo-wu-guan
// 模板题 https://www.luogu.com.cn/problem/P3381
func (*graph) minCostFlowSPFA(in io.Reader, n, m, st, end int) (int, int64) {
	const inf int = 1e9 // 1e18
	st--
	end--

	type neighbor struct{ to, rid, cap, cost, eid int } // rid 为反向边在邻接表中的下标
	g := make([][]neighbor, n)
	addEdge := func(from, to, cap, cost, eid int) {
		g[from] = append(g[from], neighbor{to, len(g[to]), cap, cost, eid})
		g[to] = append(g[to], neighbor{from, len(g[from]) - 1, 0, -cost, -1}) // 无向图上 0 换成 cap
	}
	for i := 0; i < m; i++ {
		var v, w, cp, cost int
		Fscan(in, &v, &w, &cp, &cost)
		v--
		w--
		addEdge(v, w, cp, cost, i)
	}

	dist := make([]int64, len(g))
	type vi struct{ v, i int }
	fa := make([]vi, len(g))
	spfa := func() bool {
		const _inf int64 = 1e18
		for i := range dist {
			dist[i] = _inf
		}
		dist[st] = 0
		inQ := make([]bool, len(g))
		inQ[st] = true
		q := []int{st}
		for len(q) > 0 {
			v := q[0]
			q = q[1:]
			inQ[v] = false
			for i, e := range g[v] {
				if e.cap == 0 {
					continue
				}
				w := e.to
				if newD := dist[v] + int64(e.cost); newD < dist[w] {
					dist[w] = newD
					fa[w] = vi{v, i}
					if !inQ[w] {
						q = append(q, w)
						inQ[w] = true
					}
				}
			}
		}
		return dist[end] < _inf
	}
	ek := func() (maxFlow int, minCost int64) {
		for spfa() {
			// 沿 st-end 的最短路尽量增广
			minF := inf
			for v := end; v != st; {
				p := fa[v]
				if c := g[p.v][p.i].cap; c < minF {
					minF = c
				}
				v = p.v
			}
			for v := end; v != st; {
				p := fa[v]
				e := &g[p.v][p.i]
				e.cap -= minF
				g[v][e.rid].cap += minF
				v = p.v
			}
			maxFlow += minF
			minCost += dist[end] * int64(minF)
		}
		return
	}
	return ek()
}
