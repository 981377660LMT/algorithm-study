// https://ei1333.github.io/library/graph/tree/centroid-decomposition.hpp
// https://ei1333.github.io/library/test/verify/aoj-3139.test.cpp
// https://ei1333.github.io/library/test/verify/yosupo-frequency-table-of-tree-distance.test.cpp
// https://ei1333.github.io/library/test/verify/yukicoder-1002.test.cpp

// 重心互相连接形成的有根树, 可以想象把树拎起来, 重心在树的中心，连接着各个子树的重心...
//            3 (重)
//          / | \
//     (重)1  0  2 (重)
//        / \    |
//       4   5   6

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://atcoder.jp/contests/abc291/tasks/abc291_h
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	g := make([][]Edge, n)

	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		g[u] = append(g[u], Edge{to: v})
		g[v] = append(g[v], Edge{to: u})
	}

	centTree, root := CentroidDecomposition(g)

	parents := make([]int, n)
	for i := 0; i < n; i++ {
		parents[i] = -1
	}
	var dfs func(int)
	dfs = func(cur int) {
		for _, to := range centTree[cur] {
			parents[to] = cur
			dfs(to)
		}
	}
	dfs(root)

	for _, v := range parents {
		if v == -1 {
			fmt.Fprint(out, -1, " ")
		} else {
			fmt.Fprint(out, v+1, " ")
		}
	}
}

type Edge = struct{ to, weight int }

// 树的重心分解, 返回点分树和点分树的根
//
//	!tree: `无向`树的邻接表.
//	centTree: 重心互相连接形成的有根树, 可以想象把树拎起来, 重心在树的中心，连接着各个子树的重心...
//	root: 点分树的根
func CentroidDecomposition(tree [][]Edge) (centTree [][]int, root int) {
	n := len(tree)
	subSize := make([]int, n)
	removed := make([]bool, n)
	centTree = make([][]int, n)
	var getSize func(cur, parent int) int
	var getCentroid func(cur, parent, mid int) int
	var build func(cur int) int

	getSize = func(cur, parent int) int {
		subSize[cur] = 1
		for _, e := range tree[cur] {
			next := e.to
			if next == parent || removed[next] {
				continue
			}
			subSize[cur] += getSize(next, cur)
		}
		return subSize[cur]
	}
	getCentroid = func(cur, parent, mid int) int {
		for _, e := range tree[cur] {
			next := e.to
			if next == parent || removed[next] {
				continue
			}
			if subSize[next] > mid {
				return getCentroid(next, cur, mid)
			}
		}
		return cur
	}
	build = func(cur int) int {
		centroid := getCentroid(cur, -1, getSize(cur, -1)/2)
		removed[centroid] = true
		for _, e := range tree[centroid] {
			next := e.to
			if !removed[next] {
				centTree[centroid] = append(centTree[centroid], build(next))
			}
		}
		removed[centroid] = false
		return centroid
	}

	root = build(0)
	return
}
