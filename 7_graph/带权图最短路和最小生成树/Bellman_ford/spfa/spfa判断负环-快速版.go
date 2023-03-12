package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct{ to, weight int }

// spfa 判负环
// https://www.acwing.com/solution/content/87368/
func spfaBetter(n int, graph [][]Edge) bool {
	dist := make([]int, n)
	queue := make([]int, 0, n)
	inQueue := make([]bool, n)
	pre := make([]int, n)
	for i := 0; i < n; i++ {
		queue = append(queue, i)
		inQueue[i] = true
		pre[i] = -1
	}

	idx := 0
	var detectCycle func() bool
	detectCycle = func() bool {
		vec := []int{}
		inStack := make([]bool, n)
		vis := make([]bool, n)
		for i := 0; i < n; i++ {
			if !vis[i] {
				for j := i; j != -1; j = pre[j] {
					if !vis[j] {
						vis[j] = true
						vec = append(vec, j)
						inStack[j] = true
					} else {
						if inStack[j] {
							return true
						}
						break
					}
				}
				for _, j := range vec {
					inStack[j] = false
				}
				vec = vec[:0]
			}
		}
		return false
	}

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		inQueue[u] = false
		for _, e := range graph[u] {
			v, w := e.to, e.weight
			if dist[u]+w < dist[v] {
				pre[v] = u
				dist[v] = dist[u] + w
				if idx++; idx == n {
					idx = 0
					if detectCycle() {
						return true
					}
				}
				if !inQueue[v] {
					queue = append(queue, v)
					inQueue[v] = true
				}
			}
		}
	}

	if detectCycle() {
		return true
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

	if spfaBetter(n, graph) {
		fmt.Fprintln(out, "Yes")
	} else {
		fmt.Fprintln(out, "No")
	}
}
