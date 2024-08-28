// 最短哈密尔顿回路
package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// abc180e()
	judge()
}

func judge() {
	// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=DPL_2_A
	// TSP邮递员问题
	// n<=15
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int32
	fmt.Fscan(in, &n, &m)
	graph := make([][][2]int, n)
	for i := int32(0); i < m; i++ {
		var u, v, c int
		fmt.Fscan(in, &u, &v, &c)
		graph[u] = append(graph[u], [2]int{v, c})
	}
	cost, _ := MininumHamiltonianCycleFromGraph(graph)
	fmt.Fprintln(out, cost)
}

// E - Traveling Salesman among Aerial Cities
// https://atcoder.jp/contests/abc180/tasks/abc180_e
func abc180e() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	edges := make([][3]int, n)
	for i := int32(0); i < n; i++ {
		var x, y, z int
		fmt.Fscan(in, &x, &y, &z)
		edges[i] = [3]int{x, y, z}
	}

	f := func(i, j int32) int {
		x1, y1, z1 := edges[i][0], edges[i][1], edges[i][2]
		x2, y2, z2 := edges[j][0], edges[j][1], edges[j][2]
		return abs(x1-x2) + abs(y1-y2) + max(0, z2-z1)
	}

	cost, _ := MininumHamiltonianCycle(n, f)
	fmt.Fprintln(out, cost)
}

const INF int = 1e18

// 返回最短的哈密尔顿回路的长度和路径.
//
//	不存在时返回{-1, nil}
func MininumHamiltonianCycleFromGraph(graph [][][2]int) (res int, cycle []int32) {
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
	return MininumHamiltonianCycle(n, func(i, j int32) int { return dist[i][j] })
}

// 返回最短的哈密尔顿回路的长度和路径.
//
//	不存在时返回{-1, nil}
func MininumHamiltonianCycle(n int32, cost func(i, j int32) int) (res int, cycle []int32) {
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
	n -= 1
	FULL := int32((1 << n) - 1)
	dp := make([][]int, 1<<n)
	for i := range dp {
		dp[i] = make([]int, n)
		for j := range dp[i] {
			dp[i][j] = INF
		}
	}
	for v := int32(0); v < n; v++ {
		dp[1<<v][v] = dist[n][v]
	}
	for s := int32(0); s < 1<<n; s++ {
		for from := int32(0); from < n; from++ {
			if dp[s][from] < INF {
				enumerateBits32(FULL-s, func(to int32) {
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
	bestV := int32(-1)
	for v := int32(0); v < n; v++ {
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
	cycle = []int32{n, bestV}
	t := s
	for int32(len(cycle)) <= n {
		to := cycle[len(cycle)-1]
		from := func() int32 {
			for from := int32(0); from < n; from++ {
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
