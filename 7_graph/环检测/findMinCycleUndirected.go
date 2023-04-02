// 无向图找环(返回任意一个极小的环(即环里不存在弦))
// https://github.com/maspypy/library/issues/3

package main

import (
	"bufio"
	"fmt"
	"os"
)

func yosupo() {
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
	vs, es := FindCycleUndirected(n, edges)
	if len(vs) == 0 {
		fmt.Fprintln(out, -1)
		return
	}
	fmt.Fprintln(out, len(vs))
	for _, v := range vs {
		fmt.Fprint(out, v, " ")
	}
	fmt.Fprintln(out)
	for _, e := range es {
		fmt.Fprint(out, e, " ")
	}
}

func main() {
	// [[1,3],[3,5],[5,1],[0,2],[2,4],[4,6],[6,0]]
	fmt.Println(FindCycleUndirected(7, [][2]int{{1, 3}, {3, 5}, {5, 1}, {0, 2}, {2, 4}, {4, 6}, {6, 0}}))

}

// 无向图找环(返回任意的一个极小的环).如果不存在环，返回空切片.
func FindCycleUndirected(n int, edges [][2]int) (vs []int, es []int) {
	type edge struct{ to, id int }
	graph := make([][]edge, n)
	for i, e := range edges {
		u, v := e[0], e[1]
		graph[u] = append(graph[u], edge{v, i})
		graph[v] = append(graph[v], edge{u, i})
	}

	m := len(edges)
	dep := make([]int, n)
	for i := 0; i < n; i++ {
		dep[i] = -1
	}
	usedE := make([]bool, m)
	parent := make([]int, n)

	var dfs func(int, int) int
	dfs = func(v, d int) int {
		dep[v] = d
		for _, e := range graph[v] {
			if usedE[e.id] {
				continue
			}
			if dep[e.to] != -1 {
				return v
			}
			usedE[e.id] = true
			parent[e.to] = e.id
			res := dfs(e.to, d+1)
			if res != -1 {
				return res
			}
		}
		return -1
	}

	for v := 0; v < n; v++ {
		if dep[v] != -1 {
			continue
		}
		w := dfs(v, 0)
		if w == -1 {
			continue
		}

		b, backE := -1, -1
		for {
			for _, e := range graph[w] {
				if usedE[e.id] {
					continue
				}
				if dep[e.to] > dep[w] || dep[e.to] == -1 {
					continue
				}
				b = w
				backE = e.id
			}
			if w == v {
				break
			}
			e := edges[parent[w]]
			w = e[1] + e[0] - w
		}

		a := edges[backE][1] + edges[backE][0] - b
		es = append(es, backE)
		vs = append(vs, a)
		for {
			x := vs[len(vs)-1]
			e := edges[es[len(es)-1]]
			y := e[1] + e[0] - x
			if y == a {
				break
			}
			vs = append(vs, y)
			es = append(es, parent[y])
		}
		return vs, es
	}

	return vs, es
}
