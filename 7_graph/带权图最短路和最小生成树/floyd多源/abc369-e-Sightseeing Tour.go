// abc369-e-Sightseeing Tour (经过指定边的最短路)
// https://atcoder.jp/contests/abc369/tasks/abc369_e
//
// 给定一个n个点m条边的无向带权图.
// 每次查询给出ki条边.
// 求0到n-1的最短路, 且必须经过这ki条边.
// n<=400,ki<=5,q<=3000.
//
// !全排列枚举经过的顺序，二进制枚举边的方向.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([][3]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &edges[i][0], &edges[i][1], &edges[i][2])
		edges[i][0]--
		edges[i][1]--
	}

	dist := Floyd(n, edges, false)

	start, target := 0, n-1

	query := func(eids []int) int {
		res := INF
		k := len(eids)
		enumeratePermutationsAll(k, func(indicesView []int) bool {
			for s := 0; s < 1<<k; s++ {
				sum, node := 0, start
				for i := 0; i < k; i++ {
					eid := eids[indicesView[i]]
					e := edges[eid]
					sum += e[2]
					if s>>i&1 == 1 {
						sum += dist[node][e[0]]
						node = e[1]
					} else {
						sum += dist[node][e[1]]
						node = e[0]
					}
				}
				sum += dist[node][target]
				res = min(res, sum)
			}
			return false
		})
		return res
	}

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var k int
		fmt.Fscan(in, &k)
		eids := make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &eids[j])
			eids[j]--
		}
		fmt.Fprintln(out, query(eids))
	}
}

const INF int = 1e18

func Floyd(n int, edges [][3]int, directed bool) [][]int {
	dist := make([][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int, n)
		for j := 0; j < n; j++ {
			dist[i][j] = INF
		}
		dist[i][i] = 0
	}

	for _, road := range edges {
		u, v, w := road[0], road[1], road[2]
		dist[u][v] = min(w, dist[u][v])
		if !directed {
			dist[v][u] = min(w, dist[v][u])
		}
	}

	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			if dist[i][k] == INF {
				continue
			}
			for j := 0; j < n; j++ {
				if dist[k][j] == INF {
					continue
				}
				dist[i][j] = min(dist[i][j], dist[i][k]+dist[k][j])
			}
		}
	}

	return dist
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// 生成全排列(不保证字典序).
func enumeratePermutationsAll(n int, f func(indicesView []int) (shouldBreak bool)) {
	var dfs func([]int, int) bool
	dfs = func(a []int, i int) bool {
		if i == len(a) {
			return f(a)
		}
		dfs(a, i+1)
		for j := i + 1; j < len(a); j++ {
			a[i], a[j] = a[j], a[i]
			if dfs(a, i+1) {
				return true
			}
			a[i], a[j] = a[j], a[i]
		}
		return false
	}
	ids := make([]int, n)
	for i := range ids {
		ids[i] = i
	}
	dfs(ids, 0)
}
