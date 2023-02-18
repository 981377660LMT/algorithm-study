// https://maspypy.github.io/library/graph/online_toposort.hpp
// https://codeforces.com/contest/798/problem/E
// 不预先给出图的拓扑排序

package main

// 陽にグラフを作らず、何らかのデータ構造で未訪問の行き先を探す想定。
//  set_used(v)：v を使用済に変更する
//  find_unused(v)：v の行き先を探す。なければ -1 を返すこと。
func OnlineToposort(
	n int, init func(),
	setUsed func(u int), findUnused func(u int) (v int),
	check bool,
) (order []int, ok bool) {
	init()
	done := make([]bool, n)
	var dfs func(v int)
	dfs = func(v int) {
		setUsed(v)
		done[v] = true
		for {
			to := findUnused(v)
			if to == -1 {
				break
			}
			dfs(to)
		}
		order = append(order, v)
	}

	for v := 0; v < n; v++ {
		if !done[v] {
			dfs(v)
		}
	}

	if check {
		init()
		for _, v := range order {
			setUsed(v)
			to := findUnused(v)
			if to != -1 {
				// not DAG
				return
			}
		}
	}

	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}

	ok = true
	return
}
