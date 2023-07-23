package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://www.luogu.com.cn/problem/P3806
	// 给定一棵有 n 个点的树，q次询问树上距离为 k 的点对是否存在。(优化:FreqTable in Tree 卷积求出所有距离对)
	// !n<=1e4 q<=100

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	g := make([][]Edge, n)
	for i := 0; i < n-1; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		u, v = u-1, v-1
		g[u] = append(g[u], Edge{to: v, weight: w})
		g[v] = append(g[v], Edge{to: u, weight: w})
	}

	// 全局状态
	centTree, root := CentroidDecomposition(g)
	removed := make([]bool, n)
	res := 0

	var collectDist func(cur, pre, dist int, dists *[]int)
	collectDist = func(cur, pre, dist int, dists *[]int) {
		*dists = append(*dists, dist)
		for _, e := range g[cur] {
			next, cost := e.to, e.weight
			if next != pre && !removed[next] {
				collectDist(next, cur, dist+cost, dists)
			}
		}
	}

	var decomposition func(cur, pre, limit int)
	decomposition = func(cur, pre, limit int) {
		removed[cur] = true
		for _, next := range centTree[cur] { // 点分树的子树中的答案(不经过重心)
			if !removed[next] {
				decomposition(next, cur, limit)
			}
		}
		removed[cur] = false

		counter := map[int]int{0: 1} // 经过重心的路径
		for _, e := range g[cur] {
			next, cost := e.to, e.weight
			if next == pre || removed[next] {
				continue
			}
			dist := []int{}
			collectDist(next, cur, cost, &dist)
			for _, d := range dist { // 经过重心的答案
				need := limit - d
				res += counter[need]
			}
			for _, d := range dist { // 更新状态
				counter[d]++
			}
		}

	}

	for i := 0; i < q; i++ {
		var k int
		fmt.Fscan(in, &k)
		res = 0
		decomposition(root, -1, k)
		if res > 0 {
			fmt.Fprintln(out, "AYE")
		} else {
			fmt.Fprintln(out, "NAY")
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
