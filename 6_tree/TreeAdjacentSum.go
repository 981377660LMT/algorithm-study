// 给定一棵树，每个节点有一个权值，支持以下操作：
// 1. 修改某个节点的权值;
// 2. 查询某个节点的邻接节点的权值和.

package main

import "fmt"

func main() {
	//     0
	//    / \
	//   1   2
	//  / \
	// 3   4

	n := int32(5)
	edges := [][]int32{{0, 1}, {0, 2}, {1, 3}, {1, 4}}
	values := []int{1, 2, 3, 4, 5}
	tree := make([][]int32, n)
	for _, e := range edges {
		tree[e[0]] = append(tree[e[0]], e[1])
		tree[e[1]] = append(tree[e[1]], e[0])
	}
	S := NewTreeAdjacentSum(tree, 0, values)

	fmt.Println(S.Query(0)) // 3
	fmt.Println(S.Query(1)) // 5
	fmt.Println(S.Query(3)) // 2
	S.Update(1, 10)
	fmt.Println(S.Query(0)) // 10
	fmt.Println(S.Query(1)) // 5
	fmt.Println(S.Query(3)) // 10
}

const INF int = 1e18

type E = int

func (*TreeAdjacentSum) e() E        { return -INF }
func (*TreeAdjacentSum) op(a, b E) E { return max(a, b) }

type TreeAdjacentSum struct {
	tree   [][]int32
	parent []int32
	values []E
	subSum []E
}

func NewTreeAdjacentSum(tree [][]int32, root int32, values []E) *TreeAdjacentSum {
	res := &TreeAdjacentSum{
		tree:   tree,
		parent: make([]int32, len(tree)),
		values: append(values[:0:0], values...),
		subSum: make([]E, len(tree)),
	}
	for i := range res.parent {
		res.parent[i] = -1
		res.subSum[i] = res.e()
	}
	res._dfs(root)
	return res
}

func (tas *TreeAdjacentSum) Update(u int, lazy E) {
	tas.values[u] = tas.op(tas.values[u], lazy)
	if p := tas.parent[u]; p != -1 {
		tas.subSum[p] = tas.op(tas.subSum[p], lazy)
	}
}

func (tas *TreeAdjacentSum) Query(u int) E {
	if p := tas.parent[u]; p != -1 {
		return tas.op(tas.subSum[u], tas.values[p])
	} else {
		return tas.subSum[u]
	}
}

func (tas *TreeAdjacentSum) _dfs(cur int32) {
	for _, next := range tas.tree[cur] {
		if next != tas.parent[cur] {
			tas.parent[next] = cur
			tas.subSum[cur] = tas.op(tas.subSum[cur], tas.values[next])
			tas._dfs(next)
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
