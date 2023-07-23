// https://leetcode.cn/problems/number-of-good-paths/
// 一条 好路径 需要满足以下条件：

// 开始节点和结束节点的值 相同 。
// 开始节点和结束节点中间的所有节点值都小于等于 开始节点的值（也就是说开始节点的值应该是路径上所有节点的最大值）。
// 请你返回不同好路径的数目。

package main

func numberOfGoodPaths(vals []int, edges [][]int) int {
	n := len(vals)
	g := make([][]Edge, n)

	for _, e := range edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], Edge{to: v, weight: 1})
		g[v] = append(g[v], Edge{to: u, weight: 1})
	}

	// 全局状态
	centTree, root := CentroidDecomposition(g)
	removed := make([]bool, n)
	res := 0

	var collect func(cur, pre, maxVal int, mp map[int]int)
	collect = func(cur, pre, maxVal int, mp map[int]int) {
		if vals[cur] >= maxVal {
			mp[vals[cur]]++
			maxVal = vals[cur]
		}
		for _, e := range g[cur] {
			next := e.to
			if next != pre && !removed[next] {
				collect(next, cur, maxVal, mp)
			}
		}
	}

	var decomposition func(cur, pre int)
	decomposition = func(cur, pre int) {
		removed[cur] = true
		for _, next := range centTree[cur] { // 点分树的子树中的答案(不经过重心)
			if !removed[next] {
				decomposition(next, cur)
			}
		}
		removed[cur] = false

		counter := map[int]int{vals[cur]: 1} // 经过重心的路径
		for _, e := range g[cur] {
			next := e.to
			if next == pre || removed[next] {
				continue
			}
			sub := map[int]int{} // value -> count
			collect(next, cur, vals[cur], sub)
			for k, v := range sub {
				res += v * counter[k]
				counter[k] += v
			}
		}
	}

	decomposition(root, -1)
	return res + n // 一个结点的路径
}

type Edge = struct{ to, weight int }

// 树的重心分解, 返回点分树和点分树的根
//  !tree: `无向`树的邻接表.
//  centTree: 重心互相连接形成的有根树, 可以想象把树拎起来, 重心在树的中心，连接着各个子树的重心...
//  root: 点分树的根
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
