// 给定一个有向图 $G$ 和它的反图，算法通过深度优先搜索（DFS）遍历这两个图，
// 得到它们的拓扑排序，然后按照拓扑排序的逆序进行 DFS，得到每个强连通分量。

package main

import "fmt"

func main() {
	n := 5
	g := make([][]int, n)
	rg := make([][]int, n)
	addEdge := func(from, to int) {
		g[from] = append(g[from], to)
		rg[to] = append(rg[to], from)
	}
	addEdge(0, 1)
	addEdge(1, 2)
	addEdge(2, 0)
	addEdge(1, 3)
	addEdge(3, 4)

	used := make([]bool, n)
	rUsed := make([]bool, n)
	setUsed := func(v int, rev bool) {
		if rev {
			rUsed[v] = true
		} else {
			used[v] = true
		}
	}

	// brute force
	findUnused := func(v int, rev bool) int {
		if rev {
			for _, from := range rg[v] {
				if !rUsed[from] {
					return from
				}
			}
		} else {
			for _, to := range g[v] {
				if !used[to] {
					return to
				}
			}
		}

		return -1
	}

	groupCount, groupId := OnlineSCC(n, setUsed, findUnused)
	fmt.Println(groupCount, groupId) // 3 [0 0 0 1 2]
}

// 在线求有向图的强连通分量.
//
//	setUsed(v, rev)：将 v 设置为已使用, rev 表示是否是反图
//	findUnused(v, rev)：返回未使用过的点中与 v 相邻的点, rev 表示是否是反图.不存在时返回 -1.
func OnlineSCC(n int, setUsed func(v int, rev bool), findUnused func(v int, rev bool) int) (groupCount int, groupId []int) {
	order := make([]int, n)

	next := n
	visited := make([]bool, n)
	var dfs1 func(v int)
	dfs1 = func(v int) {
		visited[v] = true
		setUsed(v, false)
		for {
			to := findUnused(v, false)
			if to == -1 {
				break
			}
			dfs1(to)
		}
		next--
		order[next] = v
	}
	for v := 0; v < n; v++ {
		if !visited[v] {
			dfs1(v)
		}
	}

	groupId = make([]int, n)
	visited = make([]bool, n)
	var dfs2 func(v int)
	dfs2 = func(v int) {
		visited[v] = true
		groupId[v] = groupCount
		setUsed(v, true)
		for {
			to := findUnused(v, true)
			if to == -1 {
				break
			}
			dfs2(to)
		}
	}
	for _, v := range order {
		if !visited[v] {
			dfs2(v)
			groupCount++
		}
	}

	return
}
