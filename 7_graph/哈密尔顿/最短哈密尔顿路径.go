// 最短哈密尔顿路径
// 指定起点和终点，恰好经过每个点一次。
// [旅行商问题](https://github.com/EndlessCheng/codeforces-go/blob/dabe70c566ef898eed871a3359e592fbd771a4b5/copypasta/dp.go#L1948)

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	acwing91()
}

// 91. 最短Hamilton路径
// https://www.acwing.com/problem/content/description/93/
func acwing91() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	grid := make([][]int, n)
	for i := int32(0); i < n; i++ {
		grid[i] = make([]int, n)
		for j := int32(0); j < n; j++ {
			fmt.Fscan(in, &grid[i][j])
		}
	}

	dist, _ := MininumHamiltonianPath(n, func(i, j int32) int { return grid[i][j] }, 0)
	fmt.Fprintln(out, dist[n-1])
}

const INF int = 1e18

// 返回最短的哈密尔顿路径的长度和路径.
func MininumHamiltonianPathFromGraph(graph [][][2]int, start int32) (minDist []int, restorePath func(target int32) []int32) {
	n := int32(len(graph))
	dist := make([][]int, n)
	for i := range dist {
		dist[i] = make([]int, n)
		for j := range dist[i] {
			dist[i][j] = INF
		}
	}
	for v, edges := range graph {
		for _, e := range edges {
			dist[v][e[0]] = min(dist[v][e[0]], e[1])
		}
	}
	return MininumHamiltonianPath(n, func(i, j int32) int { return dist[i][j] }, start)
}

// 返回最短的哈密尔顿路径的长度和路径.
func MininumHamiltonianPath(n int32, cost func(i, j int32) int, start int32) (minDist []int, restorePath func(target int32) []int32) {
	dist := make([][]int, n)
	for i := range dist {
		dist[i] = make([]int, n)
		for j := range dist[i] {
			dist[i][j] = INF
		}
	}
	for i := int32(0); i < n; i++ {
		for j := int32(0); j < n; j++ {
			if i != j {
				c := cost(i, j)
				dist[i][j] = min(dist[i][j], c)
			}
		}
	}

	full := int32((1 << n) - 1)
	dp := make([][]int, 1<<n)
	for i := range dp {
		dp[i] = make([]int, n)
		for j := range dp[i] {
			dp[i][j] = INF
		}
	}
	dp[1<<start][start] = 0

	for s := int32(0); s < 1<<n; s++ {
		for from := int32(0); from < n; from++ {
			if dp[s][from] < INF {
				enumerateBits32(full-s, func(to int32) {
					t := s | 1<<to
					cost := dist[from][to]
					if cost < INF {
						dp[t][to] = min(dp[t][to], dp[s][from]+cost)
					}
				})
			}
		}
	}

	minDist = dp[full]
	restorePath = func(target int32) []int32 {
		if minDist[target] == INF {
			return nil
		}
		path := []int32{target}
		mask := full
		for int32(len(path)) < n {
			to := path[len(path)-1]
			from := func() int32 {
				for from := int32(0); from < n; from++ {
					s := mask ^ 1<<to
					if dist[from][to] < INF && dp[s][from] < INF && dp[mask][to] == dp[s][from]+dist[from][to] {
						return from
					}
				}
				return -1
			}()
			path = append(path, from)
			mask ^= 1 << to
		}
		for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
			path[i], path[j] = path[j], path[i]
		}
		return path
	}
	return

}

// 遍历每个为1的比特位
func enumerateBits32(s int32, f func(bit int32)) {
	for s != 0 {
		i := int32(bits.TrailingZeros32(uint32(s)))
		f(i)
		s ^= 1 << i
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
