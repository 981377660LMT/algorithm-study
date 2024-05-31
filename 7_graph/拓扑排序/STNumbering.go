// 无向图的双极性方向(无向边定向)
// STNumbering/BipoarOrientation
// https://en.wikipedia.org/wiki/Bipolar_orientation
// https://maspypy.github.io/library/graph/st_numbering.hpp
//
// 给无向图的每条边定向，使该图成为起点为s、终点为t的有向无环图.
// 图的 st 编号（st-numbering）是得到的有向无环图的拓扑序.

package main

import "fmt"

func main() {
	//   0
	//  / \
	// 1 - 2
	fmt.Println(STNumbering([][]int32{{1, 2}, {0, 2}, {0, 1}}, 0, 2)) // [0 1 2], true
}

// 返回无向图定位后的st编号res.
// res[s]=0, res[t]=n-1.
// res[a] < res[b] 时, 有向边a->b，且图中存在s->b->t的路径.
func STNumbering(graph [][]int32, s, t int32) (res []int32, ok bool) {
	n := int32(len(graph))
	if n == 1 {
		return []int32{0}, true
	}
	if s == t {
		return nil, false
	}
	parent, pre, low := make([]int32, n), make([]int32, n), make([]int32, n)
	for i := int32(0); i < n; i++ {
		parent[i], pre[i], low[i] = -1, -1, -1
	}
	var path []int32

	var dfs func(int32)
	dfs = func(v int32) {
		pre[v] = int32(len(path))
		path = append(path, v)
		low[v] = v
		for _, to := range graph[v] {
			if v == to {
				continue
			}
			if pre[to] == -1 {
				dfs(to)
				parent[to] = v
				if pre[low[to]] < pre[low[v]] {
					low[v] = low[to]
				}
			} else if pre[to] < pre[low[v]] {
				low[v] = to
			}
		}
	}

	pre[s] = 0
	path = append(path, s)
	dfs(t)
	if int32(len(path)) != n {
		return nil, false
	}
	next, prev := make([]int32, n), make([]int32, n)
	for i := int32(0); i < n; i++ {
		next[i] = -1
	}
	next[s], prev[t] = t, s
	sign := make([]int32, n)
	sign[s] = -1
	for i := int32(2); i < n; i++ {
		v := path[i]
		p := parent[v]
		if sign[low[v]] == -1 {
			q := prev[p]
			if q == -1 {
				return nil, false
			}
			next[q], next[v] = v, p
			prev[v], prev[p] = q, v
			sign[p] = 1
		} else {
			q := next[p]
			if q == -1 {
				return nil, false
			}
			next[p], next[v] = v, q
			prev[v], prev[q] = p, v
			sign[p] = -1
		}
	}

	A := []int32{s}
	for A[len(A)-1] != t {
		A = append(A, next[A[len(A)-1]])
	}
	if int32(len(A)) < n {
		return nil, false
	}
	rank := make([]int32, n)
	for i := range rank {
		rank[i] = -1
	}
	for i, v := range A {
		rank[v] = int32(i)
	}
	for i := int32(0); i < n; i++ {
		var l, r bool
		v := A[i]
		for _, to := range graph[v] {
			if rank[to] < rank[v] {
				l = true
			}
			if rank[v] < rank[to] {
				r = true
			}
		}
		if i > 0 && !l {
			return nil, false
		}
		if i < n-1 && !r {
			return nil, false
		}
	}
	res = make([]int32, n)
	for i, v := range A {
		res[v] = int32(i)
	}
	return res, true
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a <= b {
		return a
	}
	return b
}
