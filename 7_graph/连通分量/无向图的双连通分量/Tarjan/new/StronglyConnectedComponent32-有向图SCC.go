// StronglyConnectedComponent-有向图SCC

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	yuki1813()
	// yosupo()
	// yuki1293()
}

func yuki1813() {
	// https://yukicoder.me/problems/no/1813
	// 不等关系:有向边; 全部相等:强连通(环)
	// 给定一个DAG 求将DAG变为一个环(强连通分量)的最少需要添加的边数
	// !答案为 `max(入度为0的点的个数, 出度为0的点的个数)`

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int32
	fmt.Fscan(in, &n, &m)

	graph := make([][]int32, n)
	for i := int32(0); i < m; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		graph[u] = append(graph[u], v)
	}
	count, belong := StronglyConnectedComponent(graph)
	if count == 1 { // 缩成一个点了,说明是强连通的
		fmt.Fprintln(out, 0)
		return
	}

	dag := SCCDag(graph, count, belong)
	indeg, outDeg := make([]int32, count), make([]int32, count)
	for i := int32(0); i < count; i++ {
		for _, next := range dag[i] {
			indeg[next]++
			outDeg[i]++
		}
	}

	in0, out0 := int32(0), int32(0)
	for i := int32(0); i < count; i++ {
		if indeg[i] == 0 {
			in0++
		}
		if outDeg[i] == 0 {
			out0++
		}
	}

	fmt.Fprintln(out, max32(in0, out0))
}

// 有向图强连通分量分解.
func StronglyConnectedComponent(graph [][]int32) (count int32, belong []int32) {
	n := int32(len(graph))
	belong = make([]int32, n)
	low := make([]int32, n)
	order := make([]int32, n)
	for i := range order {
		order[i] = -1
	}
	now := int32(0)
	path := []int32{}

	var dfs func(int32)
	dfs = func(v int32) {
		low[v] = now
		order[v] = now
		now++
		path = append(path, v)
		for _, to := range graph[v] {
			if order[to] == -1 {
				dfs(to)
				low[v] = min32(low[v], low[to])
			} else {
				low[v] = min32(low[v], order[to])
			}
		}
		if low[v] == order[v] {
			for {
				u := path[len(path)-1]
				path = path[:len(path)-1]
				order[u] = n
				belong[u] = count
				if u == v {
					break
				}
			}
			count++
		}
	}

	for i := int32(0); i < n; i++ {
		if order[i] == -1 {
			dfs(i)
		}
	}
	for i := int32(0); i < n; i++ {
		belong[i] = count - 1 - belong[i]
	}
	return
}

// 有向图的强连通分量缩点.
func SCCDag(graph [][]int32, count int32, belong []int32) (dag [][]int32) {
	dag = make([][]int32, count)
	adjSet := make([]map[int32]struct{}, count)
	for i := int32(0); i < count; i++ {
		adjSet[i] = make(map[int32]struct{})
	}
	for cur, nexts := range graph {
		for _, next := range nexts {
			if bid1, bid2 := belong[cur], belong[next]; bid1 != bid2 {
				adjSet[bid1][bid2] = struct{}{}
			}
		}
	}
	for i := int32(0); i < count; i++ {
		for next := range adjSet[i] {
			dag[i] = append(dag[i], next)
		}
	}
	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
