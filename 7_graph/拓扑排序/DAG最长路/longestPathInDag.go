package main

// 2050. 并行课程 III
// https://leetcode.cn/problems/parallel-courses-iii/description/
func minimumTime(n int, relations [][]int, time []int) int {
	dummy := n
	dag := make([][]int, n+1)
	for _, pair := range relations {
		u, v := pair[0]-1, pair[1]-1
		dag[u] = append(dag[u], v)
	}
	for i := 0; i < n; i++ {
		dag[dummy] = append(dag[dummy], i)
	}
	dp, _ := LongestPathInDag(n+1, dag, func(_, to int) int { return time[to] })
	res := 0
	for _, d := range dp {
		res = max(res, d)
	}
	return res
}

// 2770. 达到末尾下标所需的最大跳跃次数
// https://leetcode.cn/problems/maximum-number-of-jumps-to-reach-the-last-index/
func maximumJumps(nums []int, target int) int {
	n := len(nums)
	dag := make([][]int, n)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			diff := nums[j] - nums[i]
			if diff >= -target && diff <= target {
				dag[i] = append(dag[i], j)
			}
		}
	}
	return LongestPathInDagWithStart(n, dag, func(_, to int) int { return 1 }, 0, n-1)
}

// dag最长路, 并检验是否为dag(是否有环).
func LongestPathInDag(n int, dag [][]int, getWeight func(from, to int) int) (dp []int, ok bool) {
	indeg := make([]int, n)
	for _, nexts := range dag {
		for _, j := range nexts {
			indeg[j]++
		}
	}

	count := 0
	queue := []int{}
	for i, d := range indeg {
		if d == 0 {
			queue = append(queue, i)
		}
	}

	dp = make([]int, n)
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		count++
		for _, next := range dag[cur] {
			dp[next] = max(dp[next], dp[cur]+getWeight(cur, next))
			indeg[next]--
			if indeg[next] == 0 {
				queue = append(queue, next)
			}
		}
	}

	ok = count == n
	return
}

const INF int = 1e18

// 从起点出发到终点的最长路.如果无法到达,则返回-1.
func LongestPathInDagWithStart(n int, dag [][]int, getWeight func(from, to int) int, start, end int) int {
	dp := make([]int, n)
	visited := make([]bool, n)
	var dfs func(cur int) int
	dfs = func(cur int) int {
		if visited[cur] {
			return dp[cur]
		}
		visited[cur] = true
		if cur == end {
			return 0
		}

		res := -INF
		for _, next := range dag[cur] {
			res = max(res, dfs(next)+getWeight(cur, next))
		}

		dp[cur] = res
		return res
	}

	res := dfs(start)
	if visited[end] {
		return res
	}
	return -1
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
