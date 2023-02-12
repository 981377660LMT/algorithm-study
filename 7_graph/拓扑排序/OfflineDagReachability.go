// Offline Dag Reachability(DAGの到達可能性クエリ)
// https://ei1333.github.io/library/graph/others/offline-dag-reachability.hpp

// 如果图上有环,先SCC分解成DAG
// 然后64个查询一组批处理

package main

// 在有向无环图graph上,对每个查询(u,v)判断从u是否能到达v.
//  O((E+V)*Q/64)
func offlineDagReachability(dag [][]int, queries [][]int) []bool {
	n, q := len(dag), len(queries)
	order := topoSort(dag)
	res := make([]bool, q)
	for i := 0; i < q; i += 64 {
		upper := min(q, i+64)
		dp := make([]uint64, n)
		for k := i; k < upper; k++ {
			dp[queries[k][0]] |= 1 << (k - i)
		}

		for _, cur := range order {
			for _, next := range dag[cur] {
				dp[next] |= dp[cur]
			}
		}

		for k := i; k < upper; k++ {
			res[k] = dp[queries[k][1]]&(1<<(k-i)) > 0
		}
	}

	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func topoSort(dag [][]int) []int {
	n := len(dag)
	deg := make([]int, n)
	for i := 0; i < n; i++ {
		for _, v := range dag[i] {
			deg[v]++
		}
	}

	queue := []int{}
	for i := 0; i < n; i++ {
		if deg[i] == 0 {
			queue = append(queue, i)
		}
	}

	order := []int{}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		order = append(order, cur)
		for _, next := range dag[cur] {
			deg[next]--
			if deg[next] == 0 {
				queue = append(queue, next)
			}
		}
	}

	return order
}
