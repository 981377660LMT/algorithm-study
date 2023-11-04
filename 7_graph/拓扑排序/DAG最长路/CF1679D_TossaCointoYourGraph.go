// https://www.luogu.com.cn/problem/CF1679D
// Toss a Coin to Your Graph...
// 给定一个有向图，
// 每个点有一个点权，任选起点，走k步，问经过的点的最大权值最小能是多少
// 无解输出-1，没有重边和自环，但是会有环
// n,m<=2e5,k<=1e18.
//
// 二分，等价于问能否在图中走k-1步。
// 如果存在环，一定可以。
// 如果不存在环，最长链(dag最长路)需要>=k-1.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func Solve(n int, edges [][2]int, values []int, k int) int {
	if k == 1 {
		return mins(values...)
	}

	check := func(mid int) bool {
		alive := make([]bool, n)
		aliveCount := 0
		for i := 0; i < n; i++ {
			alive[i] = values[i] <= mid
			if alive[i] {
				aliveCount++
			}
		}

		adjList := make([][]int, n)
		deg := make([]int, n)
		for _, e := range edges {
			a, b := e[0], e[1]
			if alive[a] && alive[b] {
				adjList[a] = append(adjList[a], b)
				deg[b]++
			}
		}

		queue := []int{}
		for i := 0; i < n; i++ {
			if alive[i] && deg[i] == 0 {
				queue = append(queue, i)
			}
		}
		count := 0
		dp := make([]int, n)
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			count++
			for _, next := range adjList[cur] {
				dp[next] = max(dp[next], dp[cur]+1)
				deg[next]--
				if deg[next] == 0 {
					queue = append(queue, next)
				}
			}
		}

		if count < aliveCount {
			return true // 存在环
		}
		return maxs(dp...) >= k-1 //dag最长路>=k-1
	}

	left, right := 0, maxs(values...)
	ok := false
	for left <= right {
		mid := (left + right) / 2
		if check(mid) {
			right = mid - 1
			ok = true
		} else {
			left = mid + 1
		}
	}

	if ok {
		return left
	}
	return -1
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

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)
	values := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &values[i])
	}

	edges := make([][2]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &edges[i][0], &edges[i][1])
		edges[i][0]--
		edges[i][1]--
	}

	fmt.Fprintln(out, Solve(n, edges, values, k))
}
