// 最短哈密尔顿路径
// 指定起点和终点，恰好经过每个点一次。
// [旅行商问题](https://github.com/EndlessCheng/codeforces-go/blob/dabe70c566ef898eed871a3359e592fbd771a4b5/copypasta/dp.go#L1948)

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

func main() {
	// acwing91()
	// words = ["alex","loves","leetcode"]
	// fmt.Println(shortestSuperstring([]string{"alex", "loves", "leetcode"}))
	// words = ["catg","ctaagt","gcta","ttca","atgcatc"]

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

	dist, _ := MininumHamiltonianPath(n, func(i, j int32) int { return grid[i][j] }, []int32{0})
	fmt.Fprintln(out, dist[n-1])
}

// https://leetcode.cn/problems/find-the-shortest-superstring/description/
// 943. 最短超级串
// 给定一个字符串数组 words，找到以 words 中每个字符串作为子字符串的最短字符串。
// 如果有多个有效最短字符串满足题目条件，返回其中 任意一个 即可。
// 我们可以假设 words 中没有字符串是 words 中另一个字符串的子字符串。
func shortestSuperstring(words []string) string {
	compressStringNaive := func(pre, post string) int {
		m, n := len(pre), len(post)
		for res := min(m, n); res >= 0; res-- {
			if pre[m-res:] == post[:res] {
				return res
			}
		}
		return 0
	}

	n := int32(len(words))
	dist := make([][]int, n)
	for i := range dist {
		dist[i] = make([]int, n)
		for j := range dist[i] {
			dist[i][j] = -compressStringNaive(words[i], words[j])
		}
	}

	starts := make([]int32, n)
	for i := range starts {
		starts[i] = int32(i)
	}
	bestDist, bestTarget := INF, int32(-1)
	minDist, restorePath := MininumHamiltonianPath(int32(n), func(i, j int32) int { return dist[i][j] }, starts)
	for target := int32(0); target < n; target++ {
		if minDist[target] < bestDist {
			bestDist, bestTarget = minDist[target], target
		}
	}

	res := strings.Builder{}
	path := restorePath(bestTarget)
	res.WriteString(words[path[0]])
	for i := 1; i < len(path); i++ {
		common := -dist[path[i-1]][path[i]]
		curS := words[path[i]]
		res.WriteString(curS[common:])
	}
	return res.String()
}

const INF int = 1e18

// 返回最短的哈密尔顿路径的长度和路径.
func MininumHamiltonianPathFromGraph(graph [][][2]int, starts []int32) (minDist []int, restorePath func(target int32) []int32) {
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
	return MininumHamiltonianPath(n, func(i, j int32) int { return dist[i][j] }, starts)
}

// 返回最短的哈密尔顿路径的长度和路径.
func MininumHamiltonianPath(n int32, cost func(i, j int32) int, starts []int32) (minDist []int, restorePath func(target int32) []int32) {
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
	for _, start := range starts {
		dp[1<<start][start] = 0
	}

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
