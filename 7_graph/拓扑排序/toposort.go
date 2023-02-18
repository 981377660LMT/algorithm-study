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
	g := make([][]int, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		g[u] = append(g[u], v)
	}

	order, ok := topoSort(g)
	if !ok {
		fmt.Println("No")
		return
	}

	fmt.Println("Yes")
	for _, v := range order {
		fmt.Println(v)
	}
}

func topoSort(dag [][]int) (order []int, ok bool) {
	n := len(dag)
	visited, temp := make([]bool, n), make([]bool, n)
	var dfs func(int) bool
	dfs = func(i int) bool {
		if temp[i] {
			return false
		}
		if !visited[i] {
			temp[i] = true
			for _, v := range dag[i] {
				if !dfs(v) {
					return false
				}
			}
			visited[i] = true
			order = append(order, i)
			temp[i] = false
		}
		return true
	}

	for i := 0; i < n; i++ {
		if !visited[i] {
			if !dfs(i) {
				return nil, false
			}
		}
	}

	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}
	return order, true
}
