// 有向图找环(任意一个环)

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://judge.yosupo.jp/problem/cycle_detection
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([][2]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &edges[i][0], &edges[i][1])
	}
	vs, es := FindCycleDirected(n, edges)
	if len(vs) == 0 {
		fmt.Fprintln(out, -1)
		return
	}
	fmt.Fprintln(out, len(vs))
	for _, e := range es {
		fmt.Fprintln(out, e)
	}
}

// 有向图找一个环.如果不存在环，返回空切片.
func FindCycleDirected(n int, edges [][2]int) (vs []int, es []int) {
	used := make([]int, n)
	parent := make([][2]int, n)
	type edge struct{ to, id int }
	graph := make([][]edge, n)
	for i, e := range edges {
		u, v := e[0], e[1]
		graph[u] = append(graph[u], edge{v, i})
	}

	var dfs func(int)
	dfs = func(v int) {
		used[v] = 1
		for _, e := range graph[v] {
			if len(es) > 0 {
				return
			}
			if used[e.to] == 0 {
				parent[e.to] = [2]int{v, e.id}
				dfs(e.to)
			} else if used[e.to] == 1 {
				es = []int{e.id}
				cur := v
				for cur != e.to {
					es = append(es, parent[cur][1])
					cur = parent[cur][0]
				}
				for i, j := 0, len(es)-1; i < j; i, j = i+1, j-1 {
					es[i], es[j] = es[j], es[i]
				}
				return
			}
		}
		used[v] = 2
	}

	for v := 0; v < n; v++ {
		if used[v] == 0 {
			dfs(v)
		}
	}

	if len(es) == 0 {
		return
	}

	vs = make([]int, len(es))
	for i := range es {
		vs[i] = edges[es[i]][0]
	}
	return
}
