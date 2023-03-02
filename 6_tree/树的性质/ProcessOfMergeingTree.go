// https://nyaannyaan.github.io/library/tree/process-of-merging-tree.hpp
// 表示合并过程的树,按照edges中边的顺序合并顶点,edges中可以包含不合并的边.
// process-of-merging-tree

// マージ過程を表す木
// edges の昇順に木をマージしていく。(edges にはマージに使わない辺が入っていてよい)
//
// 返り値 : (graph, 補助的な頂点に対応する辺の重み, root)

// 例如 0和1合并,边权为2；0和2合并,边权为3,那么返回值为:
// tree: [[] [] [] [{3 1 2} {3 0 2}] [{4 3 3} {4 2 3}]]  (树的有向图邻接表)
// nodes: [2 3]  (辅助结点(虚拟合并的顶点)的边权,每次成功的合并都会产生一个辅助结点)
// root: 4  (根结点)

//      4(3)
//     /    \
//    3(2)   2
//   / \
//  0   1

package main

import (
	"fmt"
	"strings"
)

func main() {
	edges := []Edge{
		{0, 1, 2},
		{0, 2, 3},
		{2, 3, 9},
		{0, 4, 9},
	}
	tree, nodes, root := ProcessOfMergingTree(edges)
	fmt.Println(tree)
	fmt.Println(nodes)
	fmt.Println(root)
}

type Edge struct{ from, to, weight int }

func ProcessOfMergingTree(edges []Edge) (tree [][]Edge, nodes []int, root int) {
	n := 1
	for _, e := range edges {
		from, to := e.from, e.to
		n = max(n, max(from, to)+1)
	}

	tree = make([][]Edge, n*2-1)
	nodes = make([]int, n-1)
	uf := NewUnionFindArray(n)
	roots := make([]int, n)
	for i := 0; i < n; i++ {
		roots[i] = i
	}
	aux := n

	for _, e := range edges {
		from, to := e.from, e.to
		if uf.Find(from) == uf.Find(to) {
			continue
		}
		weight := e.weight
		f := func(big, small int) {
			tree[aux] = append(tree[aux], Edge{aux, roots[big], weight})
			tree[aux] = append(tree[aux], Edge{aux, roots[small], weight})
			roots[big], roots[small] = aux, aux
		}
		uf.UnionWithCallback(from, to, f)
		nodes[aux-n] = weight
		aux++
	}

	// aux == 2*n - 1

	root = 2*n - 2
	return
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// NewUnionFindWithCallback ...
func NewUnionFindArray(n int) *UnionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &UnionFindArray{
		Part:   n,
		rank:   rank,
		n:      n,
		parent: parent,
	}
}

type UnionFindArray struct {
	// 连通分量的个数
	Part int

	rank   []int
	n      int
	parent []int
}

func (ufa *UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}

	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	return true
}

func (ufa *UnionFindArray) UnionWithCallback(key1, key2 int, cb func(big, small int)) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	cb(root2, root1)
	return true
}

func (ufa *UnionFindArray) Find(key int) int {
	for ufa.parent[key] != key {
		ufa.parent[key] = ufa.parent[ufa.parent[key]]
		key = ufa.parent[key]
	}
	return key
}

func (ufa *UnionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func (ufa *UnionFindArray) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *UnionFindArray) Size(key int) int {
	return ufa.rank[ufa.Find(key)]
}

func (ufa *UnionFindArray) String() string {
	sb := []string{"UnionFindArray:"}
	for root, member := range ufa.GetGroups() {
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}
