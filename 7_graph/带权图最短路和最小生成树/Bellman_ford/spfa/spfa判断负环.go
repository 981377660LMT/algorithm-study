package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct{ to, weight int }

// spfa 判断负环，时间复杂度 O(V*E)
func spfa(n int, graph [][]Edge) bool {
	dist := make([]int, n)
	queue := make([]int, 0, n)
	inQueue := make([]bool, n)
	relaxedCount := make([]int, n)
	for i := 0; i < n; i++ {
		queue = append(queue, i)
		inQueue[i] = true
		relaxedCount[i] = 1
	}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		inQueue[cur] = false
		for _, e := range graph[cur] {
			next, weight := e.to, e.weight
			cand := dist[cur] + weight
			if cand < dist[next] { // !如果要正环这里需要改成 >
				dist[next] = cand
				if !inQueue[next] {
					relaxedCount[next]++
					if relaxedCount[next] >= n+1 { // +1是虚拟源点
						return true
					}
					inQueue[next] = true
					queue = append(queue, next)
				}
			}
		}
	}

	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	graph := make([][]Edge, n)
	for i := 0; i < m; i++ {
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)
		a, b = a-1, b-1
		graph[a] = append(graph[a], Edge{b, c})
	}

	if spfa(n, graph) {
		fmt.Fprintln(out, "Yes")
	} else {
		fmt.Fprintln(out, "No")
	}
}
