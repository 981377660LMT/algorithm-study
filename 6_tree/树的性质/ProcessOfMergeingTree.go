// https://nyaannyaan.github.io/library/tree/process-of-merging-tree.hpp
// 表示合并过程的树,按照edges中边的顺序合并顶点.
// process-of-merging-tree

// 例如 0和1合并,边权为2；0和2合并,边权为3,那么返回值为:
// tree: [[] [] [] [{3 1 2} {3 0 2}] [{4 3 3} {4 2 3}]]  (树的有向图邻接表)
// nodes: [2 3]  (辅助结点(虚拟合并的顶点)的边权,每次成功的合并都会产生一个辅助结点)
// root: 4  (根结点)
//
//      4(3)
//     /    \
//    3(2)   2
//   / \
//  0   1
//
// 0-n-1 为原始顶点, n-2n-2 为辅助顶点

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const MOD int = 1e9 + 7

var INV2 = Pow(2, MOD-2, MOD)

func main() {
	// https://yukicoder.me/problems/no/1451
	// 初始时有n个人
	// 给定m个操作,每次将i和j所在的班级合并,大小相等时随机选取班长,否则选取较大的班级的班长作为新班级的班长
	// 对每个人,问最后成为班长的概率

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([]Edge, 0, m)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		edges = append(edges, Edge{u, v, 1})
	}

	tree, roots := ProcessOfMergingTree(n, edges)
	subSize := make([]int, len(tree))
	var getSubSize func(int) int
	getSubSize = func(cur int) int {
		if cur < n { // 原始顶点
			subSize[cur] = 1
		}
		for _, e := range tree[cur] {
			subSize[cur] += getSubSize(e.to)
		}
		return subSize[cur]
	}
	for _, root := range roots {
		getSubSize(root)
	}

	res := make([]int, n)
	var run func(int, int)
	run = func(cur, p int) {
		if cur < n { // 原始顶点
			res[cur] = p
			return
		}

		if len(tree[cur]) == 1 { // 只有一个子节点
			run(tree[cur][0].to, p)
			return
		}

		left, right := tree[cur][0].to, tree[cur][1].to // 两个子节点
		if subSize[left] > subSize[right] {
			run(left, p)
		} else if subSize[left] < subSize[right] {
			run(right, p)
		} else {
			run(left, p*INV2%MOD)
			run(right, p*INV2%MOD)
		}
	}

	for _, root := range roots {
		run(root, 1)
	}

	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type Edge struct{ from, to, weight int }

// 表示合并过程的树,按照edges中边的顺序合并顶点.
//
//	返回: 树的有向图邻接表, 新图中的各个根节点.
func ProcessOfMergingTree(n int, edges []Edge) (tree [][]Edge, roots []int) {
	parent := make([]int, 2*n-1)
	for i := range parent {
		parent[i] = i
	}

	tree = make([][]Edge, 2*n-1)
	uf := NewUnionFindArray(n)
	aux := n
	for _, e := range edges {
		from, to := e.from, e.to
		f := func(big, small int) {
			w, p1, p2 := e.weight, parent[big], parent[small]
			tree[aux] = append(tree[aux], Edge{aux, p1, w})
			tree[aux] = append(tree[aux], Edge{aux, p2, w})
			parent[p1], parent[p2] = aux, aux
			parent[big], parent[small] = aux, aux
			aux++
		}
		uf.UnionWithCallback(from, to, f)
	}

	tree = tree[:aux]
	for i := 0; i < aux; i++ {
		if parent[i] == i {
			roots = append(roots, i)
		}
	}
	return
}

func Pow(base, exp, mod int) int {
	base %= mod
	res := 1 % mod
	for ; exp > 0; exp >>= 1 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
	}
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// NewUnionFindWithCallback ...
func NewUnionFindArray(n int) *_UnionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &_UnionFindArray{
		Part:   n,
		rank:   rank,
		n:      n,
		parent: parent,
	}
}

type _UnionFindArray struct {
	// 连通分量的个数
	Part int

	rank   []int
	n      int
	parent []int
}

func (ufa *_UnionFindArray) Union(key1, key2 int) bool {
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

func (ufa *_UnionFindArray) UnionWithCallback(key1, key2 int, cb func(big, small int)) bool {
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

func (ufa *_UnionFindArray) Find(key int) int {
	for ufa.parent[key] != key {
		ufa.parent[key] = ufa.parent[ufa.parent[key]]
		key = ufa.parent[key]
	}
	return key
}

func (ufa *_UnionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func (ufa *_UnionFindArray) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *_UnionFindArray) Size(key int) int {
	return ufa.rank[ufa.Find(key)]
}

func (ufa *_UnionFindArray) String() string {
	sb := []string{"UnionFindArray:"}
	for root, member := range ufa.GetGroups() {
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}
