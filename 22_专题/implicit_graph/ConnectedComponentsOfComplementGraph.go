// 完全图联通分量(完全图补图联通分量)

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	yosupo()
}

// https://judge.yosupo.jp/problem/connected_components_of_complement_graph
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int32
	fmt.Fscan(in, &n, &m)
	edges := make([][2]int32, 0, m)
	for i := int32(0); i < m; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		edges = append(edges, [2]int32{u, v})
	}

	groups := ConnectedComponentsOfComplementGraph(n, edges)
	fmt.Fprintln(out, len(groups))
	for _, g := range groups {
		fmt.Fprint(out, len(g))
		for _, v := range g {
			fmt.Fprint(out, " ", v)
		}
		fmt.Fprintln(out)
	}
}

func ConnectedComponentsOfComplementGraph(n int32, edges [][2]int32) (groups [][]int32) {
	graph := make([][]int32, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		graph[u] = append(graph[u], v)
		graph[v] = append(graph[v], u)
	}
	queue, ban, yet := make([]int32, n), make([]bool, n), make([]int32, n)
	for v := int32(0); v < n; v++ {
		yet[v] = v
	}

	yetPtr := n
	cut := []int32{0}
	head, tail := int32(0), int32(0)
	for yetPtr > 0 {
		yetPtr--
		queue[tail] = yet[yetPtr]
		tail++
		for head < tail {
			cur := queue[head]
			head++
			for _, to := range graph[cur] {
				ban[to] = true
			}
			for i := yetPtr - 1; i >= 0; i-- {
				to := yet[i]
				if ban[to] {
					continue
				}
				queue[tail] = to
				tail++
				yetPtr--
				yet[i], yet[yetPtr] = yet[yetPtr], yet[i]
			}
			for _, to := range graph[cur] {
				ban[to] = false
			}
		}
		cut = append(cut, head)
	}

	groups = make([][]int32, len(cut)-1)
	for i := 0; i < len(cut)-1; i++ {
		l, r := cut[i], cut[i+1]
		groups[i] = queue[l:r]
	}
	return
}
