package main

import "math/bits"

func cycleLengthQueries(n int, queries [][]int) []int {
	res := make([]int, len(queries))
	pf := PerfectBinaryTree{}
	for i, q := range queries {
		res[i] = pf.Dist(q[0], q[1]) + 1
	}
	return res
}

// 完全二叉树.
// 根节点编号为1,左右子节点编号为2*i和2*i+1.
type PerfectBinaryTree struct{}

func (pf *PerfectBinaryTree) Depth(u int) int {
	if u == 0 {
		return 0
	}
	return bits.Len(uint(u)) - 1
}

// 完全二叉树中两个节点的最近公共祖先(两个二进制数字的最长公共前缀).
func (pf *PerfectBinaryTree) Lca(u, v int) int {
	if u == v {
		return u
	}
	if u > v {
		u, v = v, u
	}
	depth1 := pf.Depth(u)
	depth2 := pf.Depth(v)
	diff := u ^ (v >> (depth2 - depth1))
	if diff == 0 {
		return u
	}
	len := bits.Len(uint(diff))
	return u >> len
}

func (pf *PerfectBinaryTree) Dist(u, v int) int {
	return pf.Depth(u) + pf.Depth(v) - 2*pf.Depth(pf.Lca(u, v))
}
