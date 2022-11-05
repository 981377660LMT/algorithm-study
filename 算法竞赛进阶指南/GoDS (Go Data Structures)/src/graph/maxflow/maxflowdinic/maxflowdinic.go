package maxflowdinic

import (
	"fmt"
	"io"
)

// 最大流 Dinic's algorithm O(V^2 * E)  二分图上为 O(E√V)
// 如果容量是浮点数，下面代码中 > 0 的判断要改成 > eps
// !https://github.dev/EndlessCheng/codeforces-go/blob/016834c19c4289ae5999988585474174224f47e2/copypasta/graph.go#L3207
func MaxFlowDinic(in io.Reader, n, m, start, end int) int {
	const inf int = 1e18
	start--
	end--

	type neighbor struct {
		to int
		// rid 为反向边在邻接表中的下标
		rid int
		// 边的残量(初始值为容量)
		cap int
		// -1表示是反向边
		eid int
	}
	graph := make([][]neighbor, n)
	addEdge := func(from, to, cap, eid int) {
		graph[from] = append(graph[from], neighbor{to, len(graph[to]), cap, eid})
		graph[to] = append(graph[to], neighbor{from, len(graph[from]) - 1, 0, -1}) // !无向图上 0 换成 cap
	}

	for i := 0; i < m; i++ {
		var v, w, cp int
		fmt.Fscan(in, &v, &w, &cp)
		v--
		w--
		addEdge(v, w, cp, i)
	}

	var dist []int // 从源点出发的距离
	bfs := func() bool {
		dist = make([]int, len(graph))
		dist[start] = 1
		queue := []int{start}
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			for _, edge := range graph[cur] {
				if next := edge.to; edge.cap > 0 && dist[next] == 0 {
					dist[next] = dist[cur] + 1
					queue = append(queue, next)
				}
			}
		}
		return dist[end] > 0
	}

	// 寻找增广路
	var iter []int // 当前弧，在其之前的边已经没有用了，避免对没有用的边进行多次检查
	var dfs func(int, int) int
	dfs = func(cur int, flow int) int {
		if cur == end {
			return flow
		}
		for ; iter[cur] < len(graph[cur]); iter[cur]++ {
			edge := &graph[cur][iter[cur]]
			if next := edge.to; edge.cap > 0 && dist[next] > dist[cur] {
				if delta := dfs(next, min(flow, edge.cap)); delta > 0 {
					edge.cap -= delta
					graph[next][edge.rid].cap += delta
					return delta
				}
			}
		}
		return 0
	}

	dinic := func() (maxFlow int) { // int64
		for bfs() {
			iter = make([]int, len(graph))
			for {
				if delta := dfs(start, inf); delta > 0 {
					maxFlow += delta
				} else {
					break
				}
			}
		}
		return
	}
	maxFlow := dinic()

	// EXTRA: 容量复原（不存原始容量的写法）
	for _, edges := range graph {
		for i, edge := range edges {
			if edge.eid >= 0 { // 正向边
				edges[i].cap += graph[edge.to][edge.rid].cap
				graph[edge.to][edge.rid].cap = 0
			}
		}
	}

	// EXTRA: 求流的分配方案（即反向边上的 cap）
	// https://loj.ac/p/115 https://www.acwing.com/problem/content/2190/
	res := make([]int, m)
	for _, edges := range graph { // v
		for _, edge := range edges {
			next, ei := edge.to, edge.eid
			if ei >= 0 { // 正向边
				res[ei] = graph[next][edge.rid].cap
			}
		}
	}

	// EXTRA: 求关键边（扩容后可以增加最大流的边）的数量
	// 关键边 v-w 需满足，在跑完最大流后：
	// 1. 这条边的流量等于其容量
	// !2. 在残余网络上，从源点可以到达 v，从 w 可以到达汇点（即从汇点顺着反向边可以到达 w）
	// http://poj.org/problem?id=3204 https://www.acwing.com/problem/content/2238/
	{
		// 在残余网络上跑 DFS，看看哪些点能从源点和汇点访问到（从汇点出发的要判断反向边的流量）
		visited1 := make([]bool, len(graph))
		var dfs1 func(int)
		dfs1 = func(cur int) {
			visited1[cur] = true
			for _, edge := range graph[cur] {
				if next := edge.to; edge.cap > 0 && !visited1[next] {
					dfs1(next)
				}
			}
		}
		dfs1(start)

		visited2 := make([]bool, len(graph))
		var dfs2 func(int)
		dfs2 = func(cur int) {
			visited2[cur] = true
			for _, edge := range graph[cur] {
				if next := edge.to; !visited2[next] && graph[next][edge.rid].cap > 0 {
					dfs2(next)
				}
			}
		}
		dfs2(end)

		res := 0
		for i, edges := range graph {
			if !visited1[i] {
				continue
			}
			for _, edge := range edges {
				// 原图的边，流量为 0（说明该边满流），且边的两端点能分别从源汇访问到
				if edge.eid >= 0 && edge.cap == 0 && visited2[edge.to] {
					res++
				}
			}
		}
	}

	return maxFlow
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
