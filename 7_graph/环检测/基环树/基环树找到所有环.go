// !基环树森林找环

package main

func maximumInvitations(favorite []int) int {
	n := len(favorite)
	adjList := make([][]int, n)
	deg := make([]int, n)
	for i, v := range favorite {
		adjList[i] = append(adjList[i], v)
		adjList[v] = append(adjList[v], i)
		deg[i] += 1
		deg[v] += 1
	}

	cycleGroup, depth := FindCycleAndCalDepth(n, adjList, deg, false)
	// 两种情况:1.所有的二元基环树里的最长链之和;2.唯一的最长环的长度
	cand1, cand2 := 0, 0
	for i := 0; i < n; i++ {
		if favorite[favorite[i]] == i {
			cand1 += 1 + depth[i]
		}
	}
	for _, cycle := range cycleGroup {
		cand2 = max(cand2, len(cycle))
	}
	return max(cand1, cand2)
}

// 无/有向基环树森林找环上的点,并记录每个点在拓扑排序中的最大深度,最外层的点深度为0.
func FindCycleAndCalDepth(n int, adjList [][]int, deg []int, isDirected bool) (cycles [][]int, depth []int) {
	depth = make([]int, n)
	startDeg := 0
	if !isDirected {
		startDeg = 1
	}
	queue := make([]int, 0)
	for i := 0; i < n; i++ {
		if deg[i] == startDeg {
			queue = append(queue, i)
		}
	}
	visited := make([]bool, n)
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
		cycles = append(cycles, path)
	}
	return
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
