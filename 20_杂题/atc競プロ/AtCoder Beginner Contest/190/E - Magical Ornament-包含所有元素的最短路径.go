//  魔法饰品
//  bfs+状压dp
//  n<=1e5 k<=17
//  !包含所有元素的最短路径和(路径可以不连续)
//  !建图，用bfs预处理出k种颜元素两两之间的最短距离，再用状压DP求解

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e18

func magicalOrnament(n int, edges [][]int, need []int) int {
	adjList := make([][]int, n)
	for _, edge := range edges {
		adjList[edge[0]] = append(adjList[edge[0]], edge[1])
		adjList[edge[1]] = append(adjList[edge[1]], edge[0])
	}

	dist := make([][]int, len(need)) // dist[i][j]表示i,j两种颜元素之间的最短距离
	for i, start := range need {
		curDist := bfs(start, adjList)
		dist[i] = make([]int, len(need))
		for j, end := range need {
			dist[i][j] = curDist[end]
		}
	}

	memo := [20][1 << 20]int{}
	for i := 0; i < 20; i++ {
		for j := 0; j < (1 << 20); j++ {
			memo[i][j] = -1
		}
	}

	var dfs func(index int, visited int) int
	dfs = func(index int, visited int) int {
		if visited == (1<<len(need))-1 {
			return 0
		}
		if memo[index][visited] != -1 {
			return memo[index][visited]
		}

		res := INF
		for next := 0; next < len(need); next++ {
			if visited&(1<<next) == 0 {
				res = min(res, dist[index][next]+dfs(next, visited|(1<<next)))
			}
		}
		memo[index][visited] = res
		return res
	}

	res := INF
	for i := 0; i < len(need); i++ {
		res = min(res, dfs(i, 1<<i))
	}
	if res == INF {
		return -1
	}
	return res + 1
}

func bfs(start int, adjList [][]int) []int {
	n := len(adjList)
	dist := make([]int, n)
	for i := 0; i < n; i++ {
		dist[i] = INF
	}
	dist[start] = 0
	queue := []int{start}
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([][]int, m)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		edges[i] = []int{u, v}
	}

	var k int
	fmt.Fscan(in, &k)
	need := make([]int, k)
	for i := 0; i < k; i++ {
		var x int
		fmt.Fscan(in, &x)
		need[i] = x - 1
	}

	fmt.Fprintln(out, magicalOrnament(n, edges, need))
}
