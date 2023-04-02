// https://maspypy.github.io/library/graph/implicit_graph/toposort.hpp
// https://codeforces.com/contest/798/problem/E
// 不预先给出图的拓扑排序(在线拓扑排序)

package main

// 在线拓扑排序, 通过函数来寻找下一个可以合并(到达)的点.
//  setUsed(u) : 删除数据结构中的u, 用于标记u已经被使用了.
//  findUnused(u) : 返回u可以到达的下一个点, 如果没有返回-1.
//  check : 是否检查是否是DAG.
//  init : 重新初始化数据结构, check为true时需要调用.
func OnlineToposort(
	n int,
	setUsed func(u int), findUnused func(u int) (v int),
	check bool,
	init func(),
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
