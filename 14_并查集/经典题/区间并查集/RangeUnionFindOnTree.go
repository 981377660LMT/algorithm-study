// 树上区间并查集 RangeUnionFindOnTree/UnionFindRangeOnTree/UnionFindTree
// https://atcoder.jp/contests/abc295/tasks/abc295_g
// 给定一棵外向树，每次修改为连一条有向边 (u,v) 保证 v 可以到达 u，
// 查询某一个节点能到达的编号最小的点。
// 发现每次连边是缩点的过程，而且一定是将 u,v 路径上强连通分量合并，显然可以用并查集维护。

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// G - Minimum Reachable City
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	parents := make([]int, n)
	for i := 1; i < n; i++ {
		fmt.Fscan(in, &parents[i])
		parents[i]--
	}
	var q int
	fmt.Fscan(in, &q)
	queries := make([][3]int, q)
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var child, ancestor int
			fmt.Fscan(in, &child, &ancestor)
			queries[i] = [3]int{op, child - 1, ancestor - 1}
		} else {
			var u int
			fmt.Fscan(in, &u)
			queries[i] = [3]int{op, u - 1, -1}
		}
	}

	uf := NewUnionFindRangeOnTree(n, parents)
	groupMin := make([]int, n)
	for i := range groupMin {
		groupMin[i] = i
	}
	for i := 0; i < q; i++ {
		op, child, ancestor := queries[i][0], queries[i][1], queries[i][2]
		if op == 1 {

			uf.UnionRange(ancestor, child, func(big, small int) {
				groupMin[big] = min(groupMin[big], groupMin[small])
			})
		} else {
			root := uf.Find(child)
			fmt.Fprintln(out, groupMin[root]+1)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type UnionFindRangeOnTree struct {
	Part        int
	n           int
	data        []int
	treeParents []int
}

func NewUnionFindRangeOnTree(n int, treeParents []int) *UnionFindRangeOnTree {
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = -1
	}
	return &UnionFindRangeOnTree{
		Part:        n,
		n:           n,
		data:        data,
		treeParents: treeParents,
	}
}

// 将child结点合并到parent结点上,返回是否合并成功.
func (ufa *UnionFindRangeOnTree) Union(parent, child int, f func(parentRoot, childRoot int)) bool {
	parent, child = ufa.Find(parent), ufa.Find(child)
	if parent == child {
		return false
	}
	ufa.data[parent] += ufa.data[child]
	ufa.data[child] = parent
	ufa.Part--
	if f != nil {
		f(parent, child)
	}
	return true
}

// 定向合并从祖先ancestor到子孙child路径上的所有节点,返回合并次数.
func (ufa *UnionFindRangeOnTree) UnionRange(ancestor, child int, f func(ancestorRoot, childRoot int)) (mergeCount int) {
	target := ufa.Find(ancestor)
	for {
		child = ufa.Find(child)
		if child == target {
			break
		}
		ufa.Union(ufa.treeParents[child], child, f)
		mergeCount++
	}
	return
}

func (ufa *UnionFindRangeOnTree) Find(key int) int {
	if ufa.data[key] < 0 {
		return key
	}
	ufa.data[key] = ufa.Find(ufa.data[key])
	return ufa.data[key]
}

func (ufa *UnionFindRangeOnTree) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func (ufa *UnionFindRangeOnTree) GetSize(key int) int {
	return -ufa.data[ufa.Find(key)]
}

func (ufa *UnionFindRangeOnTree) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *UnionFindRangeOnTree) String() string {
	sb := []string{"UnionFindRangeOnTree:"}
	groups := ufa.GetGroups()
	keys := make([]int, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, root := range keys {
		member := groups[root]
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}
