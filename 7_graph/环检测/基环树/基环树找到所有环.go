// !基环树森林找环

package main

// 100075. 有向图访问计数
// https://leetcode.cn/problems/count-visited-nodes-in-a-directed-graph/description/
func countVisitedNodes(edges []int) []int {
	n := len(edges)
	adjList := make([][]int, n)
	for i := 0; i < n; i++ {
		adjList[i] = append(adjList[i], edges[i])
	}
	groups, inCycle, belong, _ := CyclePartition(n, adjList, true)

	cache := make([]int, n)
	for i := 0; i < n; i++ {
		cache[i] = -1
	}

	var dfs func(cur int) int
	dfs = func(cur int) int {
		if cache[cur] != -1 {
			return cache[cur]
		}
		if inCycle[cur] {
			return len(groups[belong[cur]])
		}
		next := edges[cur]
		res := dfs(next) + 1
		cache[cur] = res
		return res
	}

	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[i] = dfs(i)
	}
	return res
}

// 2204. 无向图中到环的距离
// https://leetcode.cn/problems/distance-to-a-cycle-in-undirected-graph
func distanceToCycle(n int, edges [][]int) []int {
	adjList := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adjList[u] = append(adjList[u], v)
		adjList[v] = append(adjList[v], u)
	}

	groups, _, _, _ := CyclePartition(n, adjList, false)

	// 多源bfs.
	bfsMultiStart := func(starts []int, adjList [][]int) []int {
		n := len(adjList)
		dist := make([]int, n)
		for i := range dist {
			dist[i] = INF
		}
		queue := append(starts[:0:0], starts...)
		for _, start := range starts {
			dist[start] = 0
		}
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			for _, next := range adjList[cur] {
				cand := dist[cur] + 1
				if cand < dist[next] {
					dist[next] = cand
					queue = append(queue, next)
				}
			}
		}
		return dist
	}

	return bfsMultiStart(groups[0], adjList)
}

const INF int = 1e18

// 返回基环树森林的环分组信息(环的大小>=2)以及每个点在拓扑排序中的最大深度.
//  n: 点数; adjList: 邻接表; directed: 是否有向.
//  groups: 所有的环分组; inCycle: 点是否在环上; belong: 点所属的环编号(不在环中,则为-1); depth: 每个点在拓扑排序中的最大深度.最外层的点深度为0.
func CyclePartition(n int, adjList [][]int, directed bool) (groups [][]int, inCycle []bool, belong []int, depth []int) {
	deg := make([]int, n)
	if directed {
		for u := 0; u < n; u++ {
			for _, v := range adjList[u] {
				deg[v] += 1
			}
		}
	} else {
		for u := 0; u < n; u++ {
			deg[u] = len(adjList[u])
		}
	}

	startDeg := 0
	if !directed {
		startDeg = 1
	}

	depth = make([]int, n)
	visited := make([]bool, n)
	queue := make([]int, 0)
	for i := 0; i < n; i++ {
		if deg[i] == startDeg {
			queue = append(queue, i)
		}
	}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		visited[cur] = true
		for _, next := range adjList[cur] {
			depth[next] = max(depth[next], depth[cur]+1)
			deg[next] -= 1
			if deg[next] == startDeg {
				queue = append(queue, next)
			}
		}
	}

	var dfs func(cur int, path *[]int)
	dfs = func(cur int, path *[]int) {
		if visited[cur] {
			return
		}
		visited[cur] = true
		*path = append(*path, cur)
		for _, next := range adjList[cur] {
			dfs(next, path)
		}
	}

	for i := 0; i < n; i++ {
		if visited[i] {
			continue
		}
		path := make([]int, 0)
		dfs(i, &path)
		groups = append(groups, path)
	}

	inCycle = make([]bool, n)
	belong = make([]int, n)
	for gid, group := range groups {
		for _, u := range group {
			inCycle[u] = true
			belong[u] = gid
		}
	}

	return
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
