// https://leetcode.cn/problems/count-paths-that-can-form-a-palindrome-in-a-tree/
// 给你一棵 树（即，一个连通、无向且无环的图），根 节点为 0 ，由编号从 0 到 n - 1 的 n 个节点组成。
// 这棵树用一个长度为 n 、下标从 0 开始的数组 parent 表示，其中 parent[i] 为节点 i 的父节点，由于节点 0 为根节点，所以 parent[0] == -1 。
// 另给你一个长度为 n 的字符串 s ，其中 s[i] 是分配给 i 和 parent[i] 之间的边的字符。s[0] 可以忽略。
// !找出满足 u < v ，且从 u 到 v 的路径上分配的字符可以 重新排列 形成 回文 的所有节点对 (u, v) ，并返回节点对的数目。

package main

import "fmt"

func main() {
	// parent = [-1,0,0,1,1,2], s = "acaabc"
	// 8
	fmt.Println(countPalindromePaths([]int{-1, 0, 0, 1, 1, 2}, "acaabc"))

	// 	[-1,5,0,5,5,2]
	// "xsbcqq"
	// 7
	fmt.Println(countPalindromePaths([]int{-1, 5, 0, 5, 5, 2}, "xsbcqq"))
}

func countPalindromePaths(parent []int, s string) int64 {
	n := len(parent)
	tree := make([][]Edge, n)

	for cur := 1; cur < n; cur++ {
		p := parent[cur]
		cost := 1 << (s[cur] - 'a')
		tree[p] = append(tree[p], Edge{to: cur, cost: cost})
		tree[cur] = append(tree[cur], Edge{to: p, cost: cost})
	}

	// 全局状态
	centTree, root := CentroidDecomposition(tree)
	removed := make([]bool, n)
	res := 0

	var collect func(cur, pre int, state int, sub map[int]int)
	collect = func(cur, pre int, state int, sub map[int]int) {
		sub[state]++
		for _, e := range tree[cur] {
			next, cost := e.to, e.cost
			if next != pre && !removed[next] {
				collect(next, cur, state^cost, sub)
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

		counter := map[int]int{0: 1} // 经过重心的路径
		for _, e := range tree[cur] {
			next, cost := e.to, e.cost
			if next == pre || removed[next] {
				continue
			}
			sub := map[int]int{} // state -> count，统计子树内(不含cur)
			collect(next, cur, cost, sub)
			for s, v := range sub {
				res += v * counter[s]
				for i := 0; i < 26; i++ {
					res += v * counter[s^(1<<i)]
				}
			}
			for s, v := range sub {
				counter[s] += v
			}
		}
	}

	decomposition(root, -1)
	return int64(res)
}

type Edge = struct{ to, cost int }

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
