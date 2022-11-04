package maxflow

import "io"

// 最大流 Dinic's algorithm O(n^2 * m)  二分图上为 O(m√n)
// 如果容量是浮点数，下面代码中 > 0 的判断要改成 > eps
// https://en.wikipedia.org/wiki/Dinic%27s_algorithm
// https://oi-wiki.org/graph/flow/max-flow/#dinic
// https://cp-algorithms.com/graph/dinic.html
// 模板题 https://www.luogu.com.cn/problem/P3376 https://www.luogu.com.cn/problem/P2740
func (*graph) maxFlowDinic(in io.Reader, n, m, st, end int, min func(int, int) int) int {
	const inf int = 1e9 // 1e18
	st--
	end--

	type neighbor struct{ to, rid, cap, eid int } // rid 为反向边在邻接表中的下标
	g := make([][]neighbor, n)
	addEdge := func(from, to, cap, eid int) {
		g[from] = append(g[from], neighbor{to, len(g[to]), cap, eid})
		g[to] = append(g[to], neighbor{from, len(g[from]) - 1, 0, -1}) // 无向图上 0 换成 cap
	}
	for i := 0; i < m; i++ {
		var v, w, cp int
		Fscan(in, &v, &w, &cp)
		v--
		w--
		addEdge(v, w, cp, i)
	}

	var d []int // 从源点 st 出发的距离
	bfs := func() bool {
		d = make([]int, len(g))
		d[st] = 1
		q := []int{st}
		for len(q) > 0 {
			v := q[0]
			q = q[1:]
			for _, e := range g[v] {
				if w := e.to; e.cap > 0 && d[w] == 0 {
					d[w] = d[v] + 1
					q = append(q, w)
				}
			}
		}
		return d[end] > 0
	}
	// 寻找增广路
	var iter []int // 当前弧，在其之前的边已经没有用了，避免对没有用的边进行多次检查
	var dfs func(int, int) int
	dfs = func(v int, minF int) int {
		if v == end {
			return minF
		}
		for ; iter[v] < len(g[v]); iter[v]++ {
			e := &g[v][iter[v]]
			if w := e.to; e.cap > 0 && d[w] > d[v] {
				if f := dfs(w, min(minF, e.cap)); f > 0 {
					e.cap -= f
					g[w][e.rid].cap += f
					return f
				}
			}
		}
		return 0
	}
	dinic := func() (maxFlow int) { // int64
		for bfs() {
			iter = make([]int, len(g))
			for {
				if f := dfs(st, inf); f > 0 {
					maxFlow += f
				} else {
					break
				}
			}
		}
		return
	}
	maxFlow := dinic()

	// EXTRA: 容量复原（不存原始容量的写法）
	for _, es := range g {
		for i, e := range es {
			if e.eid >= 0 { // 正向边
				es[i].cap += g[e.to][e.rid].cap
				g[e.to][e.rid].cap = 0
			}
		}
	}

	// EXTRA: 求流的分配方案（即反向边上的 cap）
	// https://loj.ac/p/115 https://www.acwing.com/problem/content/2190/
	ans := make([]int, m)
	for _, es := range g { // v
		for _, e := range es {
			w, i := e.to, e.eid
			if i >= 0 { // 正向边
				ans[i] = g[w][e.rid].cap
			}
		}
	}

	// EXTRA: 求关键边（扩容后可以增加最大流的边）的数量
	// 关键边 v-w 需满足，在跑完最大流后：
	// 1. 这条边的流量等于其容量
	// 2. 在残余网络上，从源点可以到达 v，从 w 可以到达汇点（即从汇点顺着反向边可以到达 w）
	// http://poj.org/problem?id=3204 https://www.acwing.com/problem/content/2238/
	{
		// 在残余网络上跑 DFS，看看哪些点能从源点和汇点访问到（从汇点出发的要判断反向边的流量）
		vis1 := make([]bool, len(g))
		var dfs1 func(int)
		dfs1 = func(v int) {
			vis1[v] = true
			for _, e := range g[v] {
				if w := e.to; e.cap > 0 && !vis1[w] {
					dfs1(w)
				}
			}
		}
		dfs1(st)

		vis2 := make([]bool, len(g))
		var dfs2 func(int)
		dfs2 = func(v int) {
			vis2[v] = true
			for _, e := range g[v] {
				if w := e.to; !vis2[w] && g[w][e.rid].cap > 0 {
					dfs2(w)
				}
			}
		}
		dfs2(end)

		ans := 0
		for v, es := range g {
			if !vis1[v] {
				continue
			}
			for _, e := range es {
				// 原图的边，流量为 0（说明该边满流），且边的两端点能分别从源汇访问到
				if e.eid >= 0 && e.cap == 0 && vis2[e.to] {
					ans++
				}
			}
		}
	}

	return maxFlow
}
