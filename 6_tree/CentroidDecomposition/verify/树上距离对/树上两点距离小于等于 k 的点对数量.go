package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	// https://www.luogu.com.cn/problem/P4178
	// 给定一棵 n 个节点的树，每条边有边权，求出树上两点距离小于等于 k 的点对数量。
	// !n<=4e4 O(n*logn*logn)

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	g := make([][]Edge, n)
	for i := 0; i < n-1; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		u, v = u-1, v-1
		g[u] = append(g[u], Edge{to: v, cost: w})
		g[v] = append(g[v], Edge{to: u, cost: w})
	}

	// 全局状态
	centTree, root := CentroidDecomposition(g)
	removed := make([]bool, n)
	res := 0

	var collectDist func(cur, pre, dist int, dists *[]int)
	collectDist = func(cur, pre, dist int, dists *[]int) {
		*dists = append(*dists, dist)
		for _, e := range g[cur] {
			next, cost := e.to, e.cost
			if next != pre && !removed[next] {
				collectDist(next, cur, dist+cost, dists)
			}
		}
	}

	// 双指针求有序数组中 (nums[i]+nums[j]<=k,i<j) 的对数
	countPair := func(nums []int, k int) int {
		res, left, right := 0, 0, len(nums)-1
		for left < right {
			if nums[left]+nums[right] <= k {
				res += right - left
				left++
			} else {
				right--
			}
		}
		return res
	}

	var decomposition func(cur, pre, limit int)
	decomposition = func(cur, pre, limit int) {
		removed[cur] = true

		// 统计`经过重心`的长度等于k的路径数(点对数)
		// = 所有点对数 - 同一个子树中的点对数(不合法)
		allDist := []int{0}
		for _, e := range g[cur] {
			next, cost := e.to, e.cost
			if next == pre || removed[next] {
				continue
			}

			dist := []int{}
			collectDist(next, cur, cost, &dist)
			sort.Ints(dist)
			res -= countPair(dist, limit)
			allDist = append(allDist, dist...)
		}
		sort.Ints(allDist)
		res += countPair(allDist, limit)

		for _, next := range centTree[cur] { // 点分树的子树中的答案(不经过重心)
			if !removed[next] {
				decomposition(next, cur, limit)
			}
		}
		removed[cur] = false
	}

	var k int
	fmt.Fscan(in, &k)
	decomposition(root, -1, k)
	fmt.Fprintln(out, res)
}

type Edge = struct{ to, cost int }

// 树的重心分解, 返回点分树和点分树的根
//  g: 原图
//  centTree: 重心互相连接形成的有根树, 可以想象把树拎起来, 重心在树的中心，连接着各个子树的重心...
//  root: 点分树的根
func CentroidDecomposition(g [][]Edge) (centTree [][]int, root int) {
	n := len(g)
	subSize := make([]int, n)
	removed := make([]bool, n)
	centTree = make([][]int, n)

	var getSize func(cur, parent int) int
	var getCentroid func(cur, parent, mid int) int
	var build func(cur int) int
	getSize = func(cur, parent int) int {
		subSize[cur] = 1
		for _, e := range g[cur] {
			next := e.to
			if next == parent || removed[next] {
				continue
			}
			subSize[cur] += getSize(next, cur)
		}
		return subSize[cur]
	}
	getCentroid = func(cur, parent, mid int) int {
		for _, e := range g[cur] {
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
		for _, e := range g[centroid] {
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
