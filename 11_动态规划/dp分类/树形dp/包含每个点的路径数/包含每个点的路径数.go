// 给定一棵树，求包含每个点的路径数
// 每个点出现在多少条路径上
// !路径分两种情况，一种是没有父节点参与的，树形 DP 一下就行了；另一种是父节点参与的，个数就是 子树*(n-子树)

package main

import (
	"bufio"
	"fmt"
	"os"
)

// https://ac.nowcoder.com/acm/contest/272/B
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	adjList := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		adjList[u] = append(adjList[u], v)
		adjList[v] = append(adjList[v], u)
	}
	values := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &values[i])
	}

	pathCounter := CountPath(n, adjList)

	res := 0
	for i := 0; i < n; i++ {
		if pathCounter[i]&1 > 0 {
			res ^= values[i]
		}
	}
	fmt.Fprintln(out, res)
}

// 求包含每个点的路径数.路径至少有两个点.
func CountPath(n int, adjList [][]int) []int {
	res := make([]int, n)
	var dfs func(cur, pre int) int
	dfs = func(cur, pre int) int {
		count := 0
		size := 1
		for _, next := range adjList[cur] {
			if next != pre {
				subSize := dfs(next, cur)
				count += size * subSize
				size += subSize
			}
		}
		count += size * (n - size)
		res[cur] = count
		return size
	}

	dfs(0, -1)
	return res
}
