// VertexSetPathComposite
// https://judge.yosupo.jp/problem/vertex_add_path_sum
// 单点赋值/路径聚合
// 0 vertex mul add => 顶点变为 mul * x + add
// 1 root1 root2 x=> 路径和查询f(f(...(f(x)))) 模998244353 的值

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

	hld := NewHeavyLightDecomposition(n)
	values := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &values[i])
	}
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		hld.AddEdge(u, v)
	}
	hld.Build(0)

	// 线段树
	leaves := make([]E, n)
	for i := 0; i < n; i++ {
		dfn := hld.Get(i) // 注意这里
		leaves[dfn] = E{size: 1, sum: values[i]}
	}
	seg := NewLazySegTree(leaves)

	for i := 0; i < q; i++ {
		var op, vertex, add, root1, root2 int
		fmt.Fscan(in, &op)
		if op == 0 {
			fmt.Fscan(in, &vertex, &add)
			dfn := hld.Get(vertex)
			seg.Update(dfn, dfn+1, add)
		} else {
			fmt.Fscan(in, &root1, &root2)
			res := 0
			hld.ForEach(root1, root2, func(l, r int) {
				res += seg.Query(l, r).sum
			})
			fmt.Fprintln(out, res)
		}
	}
}

type HeavyLightDecomposition struct {
	tree                                                  [][]int
	dfn, dfnToNode, top, subSize, depth, parent, heavySon []int
	dfnId                                                 int // !从0开始
}

// !注意：
//  1. dfn 是 0-indexed 的.
//  2. 构建 HLD 需要调用 `Build` 方法.
//  3. 回调函数参数的 dfn区间 是左闭右开的, 即`[left, right)`.
func NewHeavyLightDecomposition(n int) *HeavyLightDecomposition {
	tree := make([][]int, n)
	dfn := make([]int, n)       // vertex => dfn
	dfnToNode := make([]int, n) // dfn => vertex
	top := make([]int, n)       // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	subSize := make([]int, n)   // 子树大小
	depth := make([]int, n)     // 深度
	parent := make([]int, n)    // 父结点
	heavySon := make([]int, n)  // 重儿子
	return &HeavyLightDecomposition{
		tree:      tree,
		dfn:       dfn,
		dfnToNode: dfnToNode,
		top:       top,
		subSize:   subSize,
		depth:     depth,
		parent:    parent,
		heavySon:  heavySon,
	}
}

// 添加无向边 u-v.
func (hld *HeavyLightDecomposition) AddEdge(u, v int) {
	hld.tree[u] = append(hld.tree[u], v)
	hld.tree[v] = append(hld.tree[v], u)
}

func (hld *HeavyLightDecomposition) Build(root int) {
	hld.build(root, -1, 0)
	hld.markTop(root, root)
}

// 返回树节点 u 对应的 dfs 序号.
//  0 <= u < n, 0 <= id < n.
func (hld *HeavyLightDecomposition) Get(u int) int {
	return hld.dfn[u]
}

// 处理树节点u到v的路径上的所有顶点.
//  回调函数内的参数是左闭右开的 dfn 区间, 即[left, right).
//   0<=left<=right<=n
func (hld *HeavyLightDecomposition) ForEach(u, v int, cb func(left, right int)) {
	for {
		if hld.dfn[u] > hld.dfn[v] {
			u, v = v, u
		}
		cb(max(hld.dfn[hld.top[v]], hld.dfn[u]), hld.dfn[v]+1)
		if hld.top[u] != hld.top[v] {
			v = hld.parent[hld.top[v]]
		} else {
			break
		}
	}
}

// 处理树节点u到v的路径上的所有边.
//  回调函数内的参数是左闭右开的 dfn 区间, 即[left, right)
//   0<=left<=right<=n
func (hld *HeavyLightDecomposition) ForEachEdge(u, v int, cb func(left, right int)) {
	for {
		if hld.dfn[u] > hld.dfn[v] {
			u, v = v, u
		}
		if hld.top[u] != hld.top[v] {
			cb(hld.dfn[hld.top[v]], hld.dfn[v]+1)
			v = hld.parent[hld.top[v]]
		} else {
			if u != v {
				cb(hld.dfn[u]+1, hld.dfn[v]+1)
			}
			break
		}
	}
}

func (hld *HeavyLightDecomposition) LCA(u, v int) int {
	for {
		if hld.dfn[u] > hld.dfn[v] {
			u, v = v, u
		}
		if hld.top[u] == hld.top[v] {
			return u
		}
		v = hld.parent[hld.top[v]]
	}
}

func (hld *HeavyLightDecomposition) Dist(u, v int) int {
	return hld.depth[u] + hld.depth[v] - 2*hld.depth[hld.LCA(u, v)]
}

func (hld *HeavyLightDecomposition) build(cur, pre, dep int) int {
	subSize, heavySize, heavySon := 1, 0, -1
	for _, next := range hld.tree[cur] {
		if next != pre {
			nextSize := hld.build(next, cur, dep+1)
			subSize += nextSize
			if nextSize > heavySize {
				heavySize, heavySon = nextSize, next
			}
		}
	}
	hld.depth[cur] = dep
	hld.subSize[cur] = subSize
	hld.heavySon[cur] = heavySon
	hld.parent[cur] = pre
	return subSize
}

func (hld *HeavyLightDecomposition) markTop(cur, top int) {
	hld.top[cur] = top
	hld.dfn[cur] = hld.dfnId
	hld.dfnId++
	hld.dfnToNode[hld.dfn[cur]] = cur
	if hld.heavySon[cur] != -1 {
		hld.markTop(hld.heavySon[cur], top)
		for _, next := range hld.tree[cur] {
			if next != hld.heavySon[cur] && next != hld.parent[cur] {
				hld.markTop(next, next)
			}
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type E = struct{ size, sum int }
type Id = int // delta

func (*LazySegTree) e() E   { return E{size: 1} }
func (*LazySegTree) id() Id { return 0 }
func (*LazySegTree) op(left, right E) E {
	return E{
		size: left.size + right.size,
		sum:  left.sum + right.sum,
	}
}
func (*LazySegTree) mapping(lazy Id, data E) E {
	if lazy == 0 {
		return data
	}
	return E{
		size: data.size,
		sum:  data.sum + data.size*lazy,
	}
}
func (*LazySegTree) composition(parentLazy, childLazy Id) Id {
	if parentLazy == 0 {
		return childLazy
	}
	return parentLazy + childLazy
}

//
//
//
//
// !template
type LazySegTree struct {
	n    int
	log  int
	size int
	data []E
	lazy []Id
}

func NewLazySegTree(
	leaves []E,
) *LazySegTree {
	tree := &LazySegTree{}
	n := int(len(leaves))
	tree.n = n
	tree.log = int(bits.Len(uint(n - 1)))
	tree.size = 1 << tree.log
	tree.data = make([]E, 2*tree.size)
	tree.lazy = make([]Id, tree.size)
	for i := range tree.data {
		tree.data[i] = tree.e()
	}
	for i := range tree.lazy {
		tree.lazy[i] = tree.id()
	}
	for i := 0; i < n; i++ {
		tree.data[tree.size+i] = leaves[i]
	}
	for i := tree.size - 1; i >= 1; i-- {
		tree.pushUp(i)
	}
	return tree
}

// 查询切片[left:right]的值
//   0<=left<=right<=len(tree.data)
func (tree *LazySegTree) Query(left, right int) E {
	if left < 0 {
		left = 0
	}
	if right > tree.n {
		right = tree.n
	}
	if left >= right {
		return tree.e()
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
	sml, smr := tree.e(), tree.e()
	for left < right {
		if left&1 != 0 {
			sml = tree.op(sml, tree.data[left])
			left++
		}
		if right&1 != 0 {
			right--
			smr = tree.op(tree.data[right], smr)
		}
		left >>= 1
		right >>= 1
	}
	return tree.op(sml, smr)
}

func (tree *LazySegTree) QueryAll() E {
	return tree.data[1]
}

// 更新切片[left:right]的值
//   0<=left<=right<=len(tree.data)
func (tree *LazySegTree) Update(left, right int, f Id) {
	if left < 0 {
		left = 0
	}
	if right > tree.n {
		right = tree.n
	}
	if left >= right {
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

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 predicate
func (tree *LazySegTree) MinLeft(right int, predicate func(data E) bool) int {
	if right == 0 {
		return 0
	}

	right += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown((right - 1) >> i)
	}

	res := tree.e()
	for {
		right--
		for right > 1 && right&1 != 0 {
			right >>= 1
		}

		if !predicate(tree.op(tree.data[right], res)) {
			for right < tree.size {
				tree.pushDown(right)
				right = right*2 + 1
				if predicate(tree.op(tree.data[right], res)) {
					res = tree.op(tree.data[right], res)
					right--
				}
			}

			return right + 1 - tree.size
		}

		res = tree.op(tree.data[right], res)
		if (right & -right) == right {
			break
		}
	}

	return 0
}

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (tree *LazySegTree) MaxRight(left int, predicate func(data E) bool) int {
	if left == tree.n {
		return tree.n
	}

	left += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(left >> i)
	}

	res := tree.e()
	for {
		for left%2 == 0 {
			left >>= 1
		}
		if !predicate(tree.op(res, tree.data[left])) {
			for left < tree.size {
				tree.pushDown(left)
				left *= 2
				if predicate(tree.op(res, tree.data[left])) {
					res = tree.op(res, tree.data[left])
					left++
				}
			}

			return left - tree.size
		}

		res = tree.op(res, tree.data[left])
		left++
		if (left & -left) == left {
			break
		}
	}

	return tree.n
}

// 单点查询(不需要 pushUp/op 操作时使用)
func (tree *LazySegTree) Get(index int) E {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	return tree.data[index]
}

// 单点赋值
func (tree *LazySegTree) Set(index int, e E) {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	tree.data[index] = e
	for i := 1; i <= tree.log; i++ {
		tree.pushUp(index >> i)
	}
}

func (tree *LazySegTree) pushUp(root int) {
	tree.data[root] = tree.op(tree.data[2*root], tree.data[2*root+1])
}

func (tree *LazySegTree) pushDown(root int) {
	if tree.lazy[root] != tree.id() {
		tree.propagate(2*root, tree.lazy[root])
		tree.propagate(2*root+1, tree.lazy[root])
		tree.lazy[root] = tree.id()
	}
}

func (tree *LazySegTree) propagate(root int, f Id) {
	tree.data[root] = tree.mapping(f, tree.data[root])
	// !叶子结点不需要更新lazy
	if root < tree.size {
		tree.lazy[root] = tree.composition(f, tree.lazy[root])
	}
}
