// 最短哈密尔顿回路
package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=DPL_2_A
	// TSP邮递员问题
	// n<=15
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	graph := make([][][2]int, n)
	for i := 0; i < m; i++ {
		var u, v, c int
		fmt.Fscan(in, &u, &v, &c)
		graph[u] = append(graph[u], [2]int{v, c})
	}
	res, _ := MininumHamiltonianCycle(graph)
	fmt.Fprintln(out, res)
}

const INF int = 1e18

// 返回最短的哈密尔顿回路的长度和路径.
//  不存在时返回{-1, nil}
func MininumHamiltonianCycle(graph [][][2]int) (res int, cycle []int) {
	n := len(graph)
	dist := make([][]int, n)
	for i := range dist {
		dist[i] = make([]int, n)
		for j := range dist[i] {
			dist[i][j] = INF
		}
	}
	for v, edges := range graph {
		for _, e := range edges {
			to, cost := e[0], e[1]
			dist[v][to] = min(dist[v][to], cost)
		}
	}
	n -= 1
	FULL := (1 << n) - 1
	dp := make([][]int, 1<<n)
	for i := range dp {
		dp[i] = make([]int, n)
		for j := range dp[i] {
			dp[i][j] = INF
		}
	}
	for v := 0; v < n; v++ {
		dp[1<<v][v] = dist[n][v]
	}
	for s := 0; s < 1<<n; s++ {
		for from := 0; from < n; from++ {
			if dp[s][from] < INF {
				enumerateBits(uint(FULL-s), func(to int) {
					t := s | 1<<to
					cost := dist[from][to]
					if cost < INF {
						dp[t][to] = min(dp[t][to], dp[s][from]+cost)
					}
				})
			}
		}
	}
	s := (1 << n) - 1
	res = INF
	bestV := -1
	for v := 0; v < n; v++ {
		if dist[v][n] < INF && dp[s][v] < INF {
			cand := dp[s][v] + dist[v][n]
			if cand < res {
				res = cand
				bestV = v
			}
		}
	}
	if res == INF {
		return -1, nil
	}
	cycle = []int{n, bestV}
	t := s
	for len(cycle) <= n {
		to := cycle[len(cycle)-1]
		from := func() int {
			for from := 0; from < n; from++ {
				s := t ^ (1 << to)
				if dp[s][from] < INF && dist[from][to] < INF && dp[s][from]+dist[from][to] == dp[t][to] {
					return from
				}
			}
			return -1
		}()
		cycle = append(cycle, from)
		t ^= 1 << to
	}

	for i, j := 0, len(cycle)-1; i < j; i, j = i+1, j-1 {
		cycle[i], cycle[j] = cycle[j], cycle[i]
	}
	return res, cycle
}

// 遍历每个为1的比特位
func enumerateBits(s uint, f func(bit int)) {
	for s != 0 {
		i := bits.TrailingZeros(s)
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
