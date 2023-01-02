// VertexSetPathComposite
// https://judge.yosupo.jp/problem/vertex_add_path_sum
// 单点赋值/路径聚合
// 0 vertex mul add => 顶点变为 mul * x + add
// 1 root1 root2 x=> 路径和查询f(f(...(f(x)))) 模998244353 的值

// 使用线段树维护

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	values := make([]S, n)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		values[i] = S{sum: x, size: 1}
	}

	tree := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	hld := HeavyLightDecomposition(n, 0, tree, values)
	for i := 0; i < q; i++ {
		var op, vertex, add, root1, root2 int
		fmt.Fscan(in, &op)
		if op == 0 {
			fmt.Fscan(in, &vertex, &add)
			hld.UpdatePath(vertex, vertex, add)
		} else {
			fmt.Fscan(in, &root1, &root2)
			fmt.Fprintln(out, hld.QueryPath(root1, root2).sum)
		}
	}
}

type HLD struct {
	UpdatePath    func(root1, root2 int, lazy F)
	QueryPath     func(root1, root2 int) S
	UpdateSubtree func(root int, lazy F)
	QuerySubtree  func(root int) S
	GetHeavyPath  func(start int) []int
}

// HeavyLightDecomposition returns a HLD struct.
//  n: the number of nodes in the tree.
//  root: the root of the tree.
//  adjList: the adjacency list of the tree.
//  vals: the values of the nodes.
// https://github.dev/EndlessCheng/codeforces-go/blob/016834c19c4289ae5999988585474174224f47e2/copypasta/graph_tree.go#L1008
func HeavyLightDecomposition(n, root int, tree [][]int, vals []S) *HLD {
	type node struct {
		// !重儿子，不存在时为 -1
		heavySon int
		// !所处轻/重链的顶点（深度最小），轻链的顶点为自身
		top int
		// 深度
		depth int
		// 子树大小
		subSize int
		// 父节点
		parent int
		// DFS 序（作为线段树中的编号，从 1 开始）
		//  !子树区间为 [dfn, dfn+size-1]
		dfn int
	}
	nodes := make([]node, n)
	dfnToNode := make([]int, n+1) // dfs序在树中所对应的节点  dfnToNode[nodes[v].dfn] == v

	// 寻找重儿子，拆分出重链与轻链
	var build func(cur, parent, depth int) int
	build = func(cur, parent, depth int) int {
		subSize, heavySize, heavySon := 1, 0, -1
		for _, next := range tree[cur] {
			if next != parent {
				nextSize := build(next, cur, depth+1)
				subSize += nextSize
				if nextSize > heavySize {
					heavySize, heavySon = nextSize, next
				}
			}
		}
		nodes[cur] = node{depth: depth, subSize: subSize, heavySon: heavySon, parent: parent}
		return subSize
	}
	build(root, -1, 0)

	// 对这些链进行维护，就要确保每个链上的节点都是连续的
	// 注意在进行重新编号的时候先访问重链，这样可以保证重链内的节点编号连续
	dfn := 0
	var markTop func(cur, top int)
	markTop = func(cur, top int) {
		node := &nodes[cur]
		node.top = top
		dfn++
		node.dfn = dfn
		dfnToNode[dfn] = cur
		if node.heavySon != -1 {
			// 优先遍历重儿子，保证在同一条重链上的点的 DFS 序是连续的
			markTop(node.heavySon, top)
			for _, next := range tree[cur] {
				if next != node.parent && next != node.heavySon {
					markTop(next, next)
				}
			}
		}
	}
	markTop(root, root)

	// 按照 DFS 序对应的点权初始化线段树 维护重链信息
	initNums := make([]S, n+1)
	for i, v := range vals {
		dfsId := nodes[i].dfn
		initNums[dfsId] = v
	}
	seg := NewLazySegTree(initNums)

	// 利用重链上跳到top来加速路径处理
	walk := func(root1, root2 int, cb func(left, right int)) {
		node1, node2 := nodes[root1], nodes[root2]
		for ; node1.top != node2.top; node1, node2 = nodes[root1], nodes[root2] {
			top1, top2 := nodes[node1.top], nodes[node2.top]
			// root1 所处的重链顶点必须比 root2 的深
			if top1.depth < top2.depth {
				root1, root2 = root2, root1
				node1, node2 = node2, node1
				top1, top2 = top2, top1
			}
			cb(top1.dfn, node1.dfn)
			// TODO: 边权下，处理轻边的情况
			root1 = top1.parent
		}
		if node1.depth > node2.depth {
			root1, root2 = root2, root1
			node1, node2 = node2, node1
		}
		cb(node1.dfn, node2.dfn)
		// TODO: 边权下，处理轻边的情况
	}

	// 寻找以 start 为 top 的重链 ,heavyPath[-1] 即为重链末端节点
	getHeavyPath := func(start int) []int {
		heavyPath := []int{start}
		for node := nodes[start]; node.heavySon != -1; node = nodes[node.heavySon] {
			heavyPath = append(heavyPath, node.heavySon)
		}
		return heavyPath
	}

	// 将节点 root1 和节点 root2 之间路径上的所有节点（包括这两个节点）更新为 lazy 的值。
	updatePath := func(root1, root2 int, lazy F) {
		walk(root1, root2, func(left, right int) { seg.Update(left, right+1, lazy) })
	}

	// 询问节点 root1 和节点 root2 之间路径上的所有节点（包括这两个节点）的聚合值.
	// % mod
	queryPath := func(root1, root2 int) S {
		res := seg.dataUnit()
		walk(root1, root2, func(left, right int) { res = seg.mergeChildren(res, seg.Query(left, right+1)) })
		return res
	}

	// 将以节点 root 为根的子树上的所有节点的权值变为 lazy 的值。
	updateSubtree := func(root int, lazy F) {
		node := nodes[root]
		seg.Update(node.dfn, node.dfn+node.subSize-1+1, lazy)
	}

	// 询问以节点 root 为根的子树上的所有节点的聚合值.
	querySubtree := func(root int) S {
		node := nodes[root]
		return seg.Query(node.dfn, node.dfn+node.subSize-1+1)
	}

	return &HLD{
		QueryPath:     queryPath,
		UpdatePath:    updatePath,
		QuerySubtree:  querySubtree,
		UpdateSubtree: updateSubtree,
		GetHeavyPath:  getHeavyPath,
	}
}

// !线段树维护的值的类型
type S = struct{ sum, size int }

// !更新操作的值的类型/懒标记的值的类型
type F = int

// !线段树维护的值的幺元
//  alias: e
func (tree *LazySegTree) dataUnit() S { return S{size: 1} }

// !更新操作/懒标记的幺元
//  alias: id
func (tree *LazySegTree) lazyUnit() F { return 0 }

// !合并左右区间的值
//  alias: op
func (tree *LazySegTree) mergeChildren(left, right S) S {
	return S{sum: (left.sum + right.sum), size: left.size + right.size}
}

// !父结点的懒标记更新子结点的值
//  alias: mapping
func (tree *LazySegTree) updateData(lazy F, data S) S {
	return S{sum: (data.sum + lazy*data.size), size: data.size}
}

// !合并父结点的懒标记和子结点的懒标记
//  alias: composition
func (tree *LazySegTree) updateLazy(parentLazy, childLazy F) F {
	return (parentLazy + childLazy)
}

func NewLazySegTree(
	leaves []S,
) *LazySegTree {
	tree := &LazySegTree{}

	n := int(len(leaves))
	tree.n = n
	tree.log = int(bits.Len(uint(n - 1)))
	tree.size = 1 << tree.log
	tree.data = make([]S, 2*tree.size)
	tree.lazy = make([]F, tree.size)
	for i := range tree.data {
		tree.data[i] = tree.dataUnit()
	}
	for i := range tree.lazy {
		tree.lazy[i] = tree.lazyUnit()
	}
	for i := 0; i < n; i++ {
		tree.data[tree.size+i] = leaves[i]
	}
	for i := tree.size - 1; i >= 1; i-- {
		tree.pushUp(i)
	}
	return tree
}

// !template
type LazySegTree struct {
	n    int
	log  int
	size int
	data []S
	lazy []F
}

// 查询切片[left:right]的值
//   0<=left<=right<=len(tree.data)
func (tree *LazySegTree) Query(left, right int) S {
	if left == right {
		return tree.dataUnit()
	}
	left += tree.size
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		if ((left >> i) << i) != left {
			tree.pushDown(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushDown((right - 1) >> i)
		}
	}
	sml, smr := tree.dataUnit(), tree.dataUnit()
	for left < right {
		if left&1 != 0 {
			sml = tree.mergeChildren(sml, tree.data[left])
			left++
		}
		if right&1 != 0 {
			right--
			smr = tree.mergeChildren(tree.data[right], smr)
		}
		left >>= 1
		right >>= 1
	}
	return tree.mergeChildren(sml, smr)
}

func (tree *LazySegTree) QueryAll() S {
	return tree.data[1]
}

// 更新切片[left:right]的值
//   0<=left<=right<=len(tree.data)
func (tree *LazySegTree) Update(left, right int, f F) {
	if left == right {
		return
	}
	left += tree.size
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		if ((left >> i) << i) != left {
			tree.pushDown(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushDown((right - 1) >> i)
		}
	}
	l2, r2 := left, right
	for left < right {
		if left&1 != 0 {
			tree.propagate(left, f)
			left++
		}
		if right&1 != 0 {
			right--
			tree.propagate(right, f)
		}
		left >>= 1
		right >>= 1
	}
	left = l2
	right = r2
	for i := 1; i <= tree.log; i++ {
		if ((left >> i) << i) != left {
			tree.pushUp(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushUp((right - 1) >> i)
		}
	}
}

func (tree *LazySegTree) pushUp(root int) {
	tree.data[root] = tree.mergeChildren(tree.data[2*root], tree.data[2*root+1])
}

func (tree *LazySegTree) pushDown(root int) {
	tree.propagate(2*root, tree.lazy[root])
	tree.propagate(2*root+1, tree.lazy[root])
	tree.lazy[root] = tree.lazyUnit()
}

func (tree *LazySegTree) propagate(root int, f F) {
	tree.data[root] = tree.updateData(f, tree.data[root])
	// !叶子结点不需要更新lazy
	if root < tree.size {
		tree.lazy[root] = tree.updateLazy(f, tree.lazy[root])
	}
}
