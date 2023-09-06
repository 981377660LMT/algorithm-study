// https://yukicoder.me/problems/no/1320
// n,m<=2000,wi<=1e9

// !O(n^2) 寻找包含某个点的最小环

package main

import (
	"bufio"
	"fmt"
	"os"
)

// https://atcoder.jp/contests/abc308/tasks/abc308_h
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	graph := make([][]int, n)
	for i := 0; i < n; i++ {
		graph[i] = make([]int, n)
		for j := 0; j < n; j++ {
			graph[i][j] = INF
		}
	}

	for i := 0; i < m; i++ {
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)
		a--
		b--
		graph[a][b] = c
		graph[b][a] = c
	}

	res := INF
	for a := 0; a < n; a++ {
		minCost, cycle := MinCostCycleWithPoint(graph, a)
		if minCost == INF {
			continue
		}
		adj := []int{cycle[1], cycle[len(cycle)-1]}
		for i := 0; i < n; i++ {
			if i == a || i == adj[0] || i == adj[1] {
				continue
			}
			res = min(res, minCost+graph[a][i])
		}
		for i := 0; i < 2; i++ {
			ori := graph[a][adj[i]]
			graph[a][adj[i]] = INF
			graph[adj[i]][a] = INF
			minCost, cycle = MinCostCycleWithPoint(graph, a)
			res = min(res, minCost+ori)
			graph[a][adj[i]] = ori
			graph[adj[i]][a] = ori
		}
	}

	if res == INF {
		fmt.Fprintln(out, -1)
	} else {
		fmt.Fprintln(out, res)
	}
}

const INF int = 1e18

// 寻找包含某个点 withPoint 的最小环. O(n^2).
//
//	返回最小环的长度和路径.
//	如果不存在最小环, 返回 INF 和 nil.
func MinCostCycleWithPoint(graph [][]int, withPoint int) (minCost int, cycle []int) {
	n := len(graph)
	dist := make([]int, n)
	parent := make([]int, n)
	g := make([]int, n)
	visited := make([]bool, n)
	for i := 0; i < n; i++ {
		dist[i] = INF
		parent[i] = -1
	}
	dist[withPoint] = 0
	g[withPoint] = withPoint

	// 构造最短路树
	for t := 0; t < n; t++ {
		mn, pos := INF+1, -1
		for i := 0; i < n; i++ {
			if !visited[i] && dist[i] < mn {
				mn = dist[i]
				pos = i
			}
		}

		if pos == -1 {
			panic("no path")
		}
		visited[pos] = true
		for i := 0; i < n; i++ {
			if cand := dist[pos] + graph[pos][i]; dist[i] > cand {
				dist[i] = cand
				parent[i] = pos
				if pos == withPoint {
					g[i] = i
				} else {
					g[i] = g[pos]
				}
			}
		}
	}

	// 找最小环
	mn := INF
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if parent[i] == j || parent[j] == i {
				continue
			}
			if g[i] == g[j] {
				continue
			}
			mn = min(mn, dist[i]+dist[j]+graph[i][j])
		}
	}

	// 求具体路径
	// 可以删掉这一段，直接返回mn来只求最小环长度
	if mn == INF {
		return INF, nil
	}

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if parent[i] == j || parent[j] == i {
				continue
			}
			if g[i] == g[j] {
				continue
			}
			if mn != dist[i]+dist[j]+graph[i][j] {
				continue
			}
			res := []int{}
			a, b := i, j
			for a != withPoint {
				res = append(res, a)
				a = parent[a]
			}
			res = append(res, a)
			for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
				res[i], res[j] = res[j], res[i]
			}
			for b != withPoint {
				res = append(res, b)
				b = parent[b]
			}
			return mn, res
		}
	}

	return INF, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
