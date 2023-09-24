// https://leetcode.cn/problems/count-valid-paths-in-a-tree/description/
// 100047. 统计树中的合法路径数目
//
// 给你一棵 n 个节点的无向树，节点编号为 1 到 n 。给你一个整数 n 和一个长度为 n - 1 的二维整数数组 edges ，
// 其中 edges[i] = [ui, vi] 表示节点 ui 和 vi 在树中有一条边。
//
// 请你返回树中的 合法路径数目 。
//
// 如果在节点 a 到节点 b 之间 恰好有一个 节点的编号是质数，那么我们称路径 (a, b) 是 合法的 。
//
// 注意：
//
// 路径 (a, b) 指的是一条从节点 a 开始到节点 b 结束的一个节点序列，序列中的节点 互不相同 ，且相邻节点之间在树上有一条边。
// 路径 (a, b) 和路径 (b, a) 视为 同一条 路径，且只计入答案 一次 。

package main

import (
	"math"
)

func countPaths(n int, edges [][]int) int64 {
	tree := make([][]Edge, n)
	for _, e := range edges {
		u, v := e[0]-1, e[1]-1
		tree[u] = append(tree[u], Edge{to: v, cost: 1})
		tree[v] = append(tree[v], Edge{to: u, cost: 1})
	}

	// 全局状态
	centTree, root := CentroidDecomposition(tree)
	removed := make([]bool, n)
	res := 0

	var collect func(cur, pre int, primeCount int, mp map[int]int)
	collect = func(cur, pre int, primeCount int, mp map[int]int) {
		nextCount := primeCount
		if E.IsPrime(cur + 1) {
			nextCount++
		}
		if nextCount <= 1 {
			mp[nextCount]++
		} else {
			return
		}

		for _, e := range tree[cur] {
			next := e.to
			if next != pre && !removed[next] {
				collect(next, cur, nextCount, mp)
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

		init := 0
		if E.IsPrime(cur + 1) {
			init = 1
		}
		counter := map[int]int{init: 1} // 经过重心的路径
		for _, e := range tree[cur] {
			next := e.to
			if next == pre || removed[next] {
				continue
			}

			sub := map[int]int{} // 统计子树内(不含cur)
			collect(next, cur, init, sub)
			if init == 0 {
				res += sub[0] * counter[1]
				res += sub[1] * counter[0]
				counter[0] += sub[0]
				counter[1] += sub[1]
			} else {
				res += sub[1] * counter[1]
				counter[1] += sub[1]
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

var E *eratosthenesSieve

func init() {
	E = newEratosthenesSieve(1e5 + 10)
}

// 埃氏筛
type eratosthenesSieve struct {
	minPrime []int
}

func newEratosthenesSieve(maxN int) *eratosthenesSieve {
	minPrime := make([]int, maxN+1)
	for i := range minPrime {
		minPrime[i] = i
	}
	upper := int(math.Sqrt(float64(maxN))) + 1
	for i := 2; i < upper; i++ {
		if minPrime[i] < i {
			continue
		}
		for j := i * i; j <= maxN; j += i {
			if minPrime[j] == j {
				minPrime[j] = i
			}
		}
	}
	return &eratosthenesSieve{minPrime}
}

func (es *eratosthenesSieve) IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	return es.minPrime[n] == n
}

func (es *eratosthenesSieve) GetPrimeFactors(n int) map[int]int {
	res := make(map[int]int)
	for n > 1 {
		m := es.minPrime[n]
		res[m]++
		n /= m
	}
	return res
}

func (es *eratosthenesSieve) GetPrimes() []int {
	res := []int{}
	for i, x := range es.minPrime {
		if i >= 2 && i == x {
			res = append(res, x)
		}
	}
	return res
}
