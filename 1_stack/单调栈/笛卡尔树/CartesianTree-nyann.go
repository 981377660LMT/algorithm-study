// https://nyaannyaan.github.io/library/tree/cartesian-tree.hpp
// 笛卡尔树

// 数列(an)が与えられたとき、
// 区間[0,n)に対して次の操作を再帰的に繰り返すことで得られる二分木をCartesian Treeと呼ぶ。
//  1. minarg: 区間[l,r)が与えられたとき、i∈[l,r)のうちa_iが最小のi(複数ある場合はiが最も小さい点)をmと置く。
//  2. rec: 頂点mを根、区間[l,m)から得られる二分木を左の子、区間[m+1,r)から得られる二分木を右の子とした木を返す。

// !Cartesian Treeの長所として「頂点(u,v)のLCAが区間[u,v]の最小値に対応する」という点があり、
// 前計算O(n)-クエリO(1)のRMQなどに利用されている。

// CartesianTree(a) : 配列aに対して{グラフ,根}を返す。

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	graph, root := CartesianTree(nums)
	parents := make([]int, n)
	parents[root] = root
	for i := 0; i < n; i++ {
		for _, j := range graph[i] {
			parents[j] = i
		}
	}
	for i := 0; i < n; i++ {
		fmt.Fprint(out, parents[i], " ")
	}
}

// 配列 nums に対して{グラフ,根}を返す。
//  グラフは根付き有向木.
func CartesianTree(nums []int) (tree [][]int, root int) {
	n := len(nums)
	tree = make([][]int, n)
	parent := make([]int, n)
	for i := range parent {
		parent[i] = -1
	}

	stack := make([]int, 0, n)
	for i, v := range nums {
		prv := -1
		for len(stack) > 0 && nums[stack[len(stack)-1]] > v {
			prv = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		}
		if prv != -1 {
			parent[prv] = i
		}
		if len(stack) > 0 {
			parent[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}

	root = -1
	for i := 0; i < n; i++ {
		if parent[i] != -1 {
			tree[parent[i]] = append(tree[parent[i]], i)
		} else {
			root = i
		}
	}
	return
}
