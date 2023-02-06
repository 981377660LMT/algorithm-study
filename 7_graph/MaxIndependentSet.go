// 无向图最大独立集(独立集是指图 G 中两两互不相邻的顶点构成的集合)
// n<=40

package main

import (
	"bufio"
	"fmt"
	"os"
)

func maxIndepentSet(n int, edges [][]int) []int {
	adjList := make([][]int, n)
	deg := make([]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adjList[u] = append(adjList[u], v)
		adjList[v] = append(adjList[v], u)
		deg[u]++
		deg[v]++
	}

	used, dead := make([]int, n), make([]int, n)
	pre, ans := []int{}, []int{}
	res, count, alive := 0, 0, n

	var dfs func()
	dfs = func() {
		if count+alive <= res {
			return
		}

		v := -1
		for i := 0; i < n; i++ {
			if used[i] != 0 || dead[i] != 0 {
				continue
			}
			if deg[i] <= 1 {
				v = i
				break
			}
			if v < 0 || deg[v] < deg[i] {
				v = i
			}
		}

		if v < 0 {
			return
		}

		if deg[v] != 1 {
			dead[v] = 1
			alive--
			for _, u := range adjList[v] {
				deg[u]--
			}
			dfs()
			dead[v] = 0
			alive++
			for _, u := range adjList[v] {
				deg[u]++
			}
		}

		used[v] = 1
		alive--
		for _, u := range adjList[v] {
			if dead[u] == 0 {
				dead[u]++
				if used[u] == 0 {
					alive--
				}
			} else {
				dead[u]++
			}
		}
		count++

		if count > res {
			pre = append(pre[:0:0], used...)
		}
		res = max(res, count)

		dfs()
		used[v] = 0
		alive++
		for _, u := range adjList[v] {
			dead[u]--
			if dead[u] == 0 {
				if used[u] == 0 {
					alive++
				}
			}
		}
		count--
	}

	dfs()
	for i := 0; i < n; i++ {
		if pre[i] != 0 {
			ans = append(ans, i)
		}
	}
	return ans
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
	edges := make([][]int, 0, m)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		edges = append(edges, []int{u, v})
	}

	res := maxIndepentSet(n, edges)
	fmt.Fprintln(out, len(res))
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}
