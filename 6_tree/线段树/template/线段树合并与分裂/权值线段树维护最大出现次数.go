// https://zhuanlan.zhihu.com/p/575513452
// https://www.luogu.com/article/bsm4zrgr

package main

import (
	"bufio"
	"fmt"
	"os"
)

// Dominant Indices (对每个结点，查询子树中哪一层结点最多)
// https://www.luogu.com.cn/problem/CF1009F
// 对于树上每个节点node，求一个最小的k，使得其子树中到node距离为k的节点数最多。
//
// 维护(major, maxCount)，major 表示出现次数最多的距离，maxCount 表示最多的次数.
func CF1009F() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	tree := make([][]int32, n)
	for i := int32(0); i < n-1; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	seg := NewSegmentTreeMerger(0, n)
	roots := make([]*Node, n)
	for i := int32(0); i < n; i++ {
		roots[i] = seg.Alloc()
	}

	res := make([]int32, n)
	var dfs func(cur, pre, dep int32)
	dfs = func(cur, pre, dep int32) {
		seg.Set(roots[cur], dep, E{major: dep, maxCount: 1})
		for _, next := range tree[cur] {
			if next == pre {
				continue
			}
			dfs(next, cur, dep+1)
			roots[cur] = seg.MergeDestructively(roots[cur], roots[next])
		}
		res[cur] = seg.QueryAll(roots[cur]).major - dep
	}
	dfs(0, -1, 0)

	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// P4556 [Vani有约会] 雨天的尾巴 /【模板】线段树合并
// https://www.luogu.com.cn/problem/P4556
// 村落里的一共有n座房屋，并形成一个树状结构。
// 然后救济粮分m次发放，每次选择两个房屋x和y ，然后对于x到y的路径上(含x和y)每座房子里发放一袋z类型的救济粮。
// 深绘里想知道，当所有的救济粮发放完毕后，每座房子里存放的最多的是哪种救济粮。
// 如果有多种救济粮都是存放最多的，输出种类编号最小的一种。
// 如果某座房屋没有救济粮，则输出 0。
//
// 每个节点开一棵权值线段树，维护(major, maxCount)，major 表示出现次数最多的救济粮，maxCount 表示救济粮的最多次数.
func P4556() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MAX int32 = 1e5 + 10

	var n, q int32
	fmt.Fscan(in, &n, &q)
	tree := make([][]int32, n)
	for i := int32(0); i < n-1; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	seg := NewSegmentTreeMerger(0, MAX)
	roots := make([]*Node, n)
	for i := int32(0); i < n; i++ {
		roots[i] = seg.Alloc()
	}
	L := NewLCA(tree, []int{0})

	for i := int32(0); i < q; i++ {
		var u, v, z int32
		fmt.Fscan(in, &u, &v, &z)
		u, v = u-1, v-1
		lca := L.LCA(int(u), int(v))
		seg.Update(roots[u], z, E{major: z, maxCount: 1})

	}

}

type Node struct {
	Major                 int32 // 出现次数最多的权值.
	MaxCount              int32 // 出现次数最多的权值出现的次数.
	leftChild, rightChild *Node
}

func (n *Node) String() string {
	return fmt.Sprintf("%v", n.Major)
}

type SegmentTreeOnRange struct {
	min, max int32
}

// 指定闭区间[min,max]建立权值线段树.
func NewSegmentTreeOnRange(min, max int32) *SegmentTreeOnRange {
	return &SegmentTreeOnRange{min: min, max: max}
}

// NewRoot().
func (sm *SegmentTreeOnRange) Alloc() *Node {
	return &Node{}
}

// 权值线段树求第 k 小.
// 调用前需保证 1 <= k <= node.value.
func (sm *SegmentTreeOnRange) Kth(node *Node, k int32) (value int32, ok bool) {
	if k < 1 || k > sm._eval(node) {
		return 0, false
	}
	return sm._kth(k, node, sm.min, sm.max), true
}

func (sm *SegmentTreeOnRange) Get(node *Node, index int32) E {
	return sm._get(node, index, sm.min, sm.max)
}

func (sm *SegmentTreeOnRange) Set(node *Node, index int32, value E) {
	sm._set(node, index, value, sm.min, sm.max)
}

func (sm *SegmentTreeOnRange) Query(node *Node, left, right int32) E {
	return sm._query(node, left, right, sm.min, sm.max)
}

func (sm *SegmentTreeOnRange) QueryAll(node *Node) E {
	return sm._eval(node)
}

func (sm *SegmentTreeOnRange) Update(node *Node, index int32, value E) {
	sm._update(node, index, value, sm.min, sm.max)
}

// 用一个新的节点存合并的结果，会生成重合节点数量的新节点.
func (sm *SegmentTreeOnRange) Merge(a, b *Node) *Node {
	return sm._merge(a, b, sm.min, sm.max)
}

// 把第二棵树直接合并到第一棵树上，比较省空间，缺点是会丢失合并前树的信息.
func (sm *SegmentTreeOnRange) MergeDestructively(a, b *Node) *Node {
	return sm._mergeDestructively(a, b, sm.min, sm.max)
}

// 线段树分裂，将区间 [left,right] 从原树分离到 other 上, this 为原树的剩余部分.
func (sm *SegmentTreeOnRange) Split(node *Node, left, right int32) (this, other *Node) {
	this, other = sm._split(node, nil, left, right, sm.min, sm.max)
	return
}

func (sm *SegmentTreeOnRange) _kth(k int32, node *Node, left, right int32) int32 {
	if left == right {
		return left
	}
	mid := (left + right) >> 1
	if leftCount := sm._eval(node.leftChild); leftCount >= k {
		return sm._kth(k, node.leftChild, left, mid)
	} else {
		return sm._kth(k-leftCount, node.rightChild, mid+1, right)
	}
}

func (sm *SegmentTreeOnRange) _get(node *Node, index int32, left, right int32) E {
	if node == nil {
		return 0
	}
	if left == right {
		return node.Major
	}
	mid := (left + right) >> 1
	if index <= mid {
		return sm._get(node.leftChild, index, left, mid)
	} else {
		return sm._get(node.rightChild, index, mid+1, right)
	}
}

func (sm *SegmentTreeOnRange) _query(node *Node, L, R int32, left, right int32) E {
	if node == nil {
		return 0
	}
	if L <= left && right <= R {
		return node.Major
	}
	mid := (left + right) >> 1
	if R <= mid {
		return sm._query(node.leftChild, L, R, left, mid)
	}
	if L > mid {
		return sm._query(node.rightChild, L, R, mid+1, right)
	}
	return sm._query(node.leftChild, L, R, left, mid) + sm._query(node.rightChild, L, R, mid+1, right)
}

func (sm *SegmentTreeOnRange) _set(node *Node, index int32, value E, left, right int32) {
	if left == right {
		node.Major = value
		return
	}
	mid := (left + right) >> 1
	if index <= mid {
		if node.leftChild == nil {
			node.leftChild = sm.Alloc()
		}
		sm._set(node.leftChild, index, value, left, mid)
	} else {
		if node.leftChild == nil {
			node.leftChild = sm.Alloc()
		}
		sm._set(node.rightChild, index, value, mid+1, right)
	}
	node.Major = sm._eval(node.leftChild) + sm._eval(node.rightChild)
}

func (sm *SegmentTreeOnRange) _update(node *Node, index int32, value E, left, right int32) {
	if left == right {
		node.Major += value
		return
	}
	mid := (left + right) >> 1
	if index <= mid {
		if node.leftChild == nil {
			node.leftChild = sm.Alloc()
		}
		sm._update(node.leftChild, index, value, left, mid)
	} else {
		if node.rightChild == nil {
			node.rightChild = sm.Alloc()
		}
		sm._update(node.rightChild, index, value, mid+1, right)
	}
	node.Major = sm._eval(node.leftChild) + sm._eval(node.rightChild)
}

func (sm *SegmentTreeOnRange) _merge(a, b *Node, left, right int32) *Node {
	if a == nil || b == nil {
		if a == nil {
			return b
		}
		return a
	}
	newNode := sm.Alloc()
	if left == right {
		newNode.Major = a.Major + b.Major
		return newNode
	}
	mid := (left + right) >> 1
	newNode.leftChild = sm._merge(a.leftChild, b.leftChild, left, mid)
	newNode.rightChild = sm._merge(a.rightChild, b.rightChild, mid+1, right)
	newNode.Major = sm._eval(newNode.leftChild) + sm._eval(newNode.rightChild)
	return newNode
}

func (sm *SegmentTreeOnRange) _mergeDestructively(a, b *Node, left, right int32) *Node {
	if a == nil || b == nil {
		if a == nil {
			return b
		}
		return a
	}
	if left == right {
		a.Major += b.Major
		return a
	}
	mid := (left + right) >> 1
	a.leftChild = sm._mergeDestructively(a.leftChild, b.leftChild, left, mid)
	a.rightChild = sm._mergeDestructively(a.rightChild, b.rightChild, mid+1, right)
	a.Major = sm._eval(a.leftChild) + sm._eval(a.rightChild)
	return a
}

func (sm *SegmentTreeOnRange) _split(a, b *Node, L, R int32, left, right int32) (*Node, *Node) {
	if a == nil || L > right || R < left {
		return a, nil
	}
	if L <= left && right <= R {
		return nil, a
	}
	if b == nil {
		b = sm.Alloc()
	}
	mid := (left + right) >> 1
	a.leftChild, b.leftChild = sm._split(a.leftChild, b.leftChild, L, R, left, mid)
	a.rightChild, b.rightChild = sm._split(a.rightChild, b.rightChild, L, R, mid+1, right)
	a.Major = sm._eval(a.leftChild) + sm._eval(a.rightChild)
	b.Major = sm._eval(b.leftChild) + sm._eval(b.rightChild)
	return a, b
}

func (sm *SegmentTreeOnRange) _eval(node *Node) E {
	if node == nil {
		return 0
	}
	return node.Major
}

type LCAFast struct {
	Depth, Parent           []int32
	Tree                    [][]int32
	lid, rid, top, heavySon []int32
	idToNode                []int32
	dfnId                   int32
}

func NewLCA(tree [][]int32, roots []int) *LCAFast {
	n := len(tree)
	lid := make([]int32, n) // vertex => dfn
	rid := make([]int32, n)
	top := make([]int32, n)      // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	depth := make([]int32, n)    // 深度
	parent := make([]int32, n)   // 父结点
	heavySon := make([]int32, n) // 重儿子
	idToNode := make([]int32, n)
	res := &LCAFast{
		Tree:     tree,
		lid:      lid,
		rid:      rid,
		top:      top,
		Depth:    depth,
		Parent:   parent,
		heavySon: heavySon,
		idToNode: idToNode,
	}
	for _, root := range roots {
		root32 := int32(root)
		res._build(root32, -1, 0)
		res._markTop(root32, root32)
	}
	return res
}

// 树上多个点的 LCA，就是 DFS 序最小的和 DFS 序最大的这两个点的 LCA.
func (hld *LCAFast) LCAMultiPoint(nodes []int) int {
	if len(nodes) == 1 {
		return nodes[0]
	}
	if len(nodes) == 2 {
		return hld.LCA(nodes[0], nodes[1])
	}
	minDfn, maxDfn := int32(1<<31-1), int32(-1)
	for _, root := range nodes {
		root32 := int32(root)
		if hld.lid[root32] < minDfn {
			minDfn = hld.lid[root32]
		}
		if hld.lid[root32] > maxDfn {
			maxDfn = hld.lid[root32]
		}
	}
	u, v := hld.idToNode[minDfn], hld.idToNode[maxDfn]
	return hld.LCA(int(u), int(v))
}

func (hld *LCAFast) LCA(u, v int) int {
	u32, v32 := int32(u), int32(v)
	for {
		if hld.lid[u32] > hld.lid[v32] {
			u32, v32 = v32, u32
		}
		if hld.top[u32] == hld.top[v32] {
			return int(u32)
		}
		v32 = hld.Parent[hld.top[v32]]
	}
}

func (hld *LCAFast) Dist(u, v int) int {
	return int(hld.Depth[u] + hld.Depth[v] - 2*hld.Depth[hld.LCA(u, v)])
}

func (hld *LCAFast) EnumeratePathDecomposition(u, v int, vertex bool, f func(start, end int)) {
	u32, v32 := int32(u), int32(v)
	for {
		if hld.top[u32] == hld.top[v32] {
			break
		}
		if hld.lid[u32] < hld.lid[v32] {
			a, b := hld.lid[hld.top[v32]], hld.lid[v32]
			if a > b {
				a, b = b, a
			}
			f(int(a), int(b+1))
			v32 = hld.Parent[hld.top[v32]]
		} else {
			a, b := hld.lid[u32], hld.lid[hld.top[u32]]
			if a > b {
				a, b = b, a
			}
			f(int(a), int(b+1))
			u32 = hld.Parent[hld.top[u32]]
		}
	}

	edgeInt := int32(1)
	if vertex {
		edgeInt = 0
	}

	if hld.lid[u32] < hld.lid[v32] {
		a, b := hld.lid[u32]+edgeInt, hld.lid[v32]
		if a > b {
			a, b = b, a
		}
		f(int(a), int(b+1))
	} else if hld.lid[v32]+edgeInt <= hld.lid[u32] {
		a, b := hld.lid[u32], hld.lid[v32]+edgeInt
		if a > b {
			a, b = b, a
		}
		f(int(a), int(b+1))
	}
}

// k: 0-based
//
//	如果不存在第k个祖先，返回-1
func (hld *LCAFast) KthAncestor(root, k int) int {
	root32 := int32(root)
	k32 := int32(k)
	if k32 > hld.Depth[root32] {
		return -1
	}
	for {
		u := hld.top[root32]
		if hld.lid[root32]-k32 >= hld.lid[u] {
			return int(hld.idToNode[hld.lid[root32]-k32])
		}
		k32 -= hld.lid[root32] - hld.lid[u] + 1
		root32 = hld.Parent[u]
	}
}

// 从 from 节点跳向 to 节点,跳过 step 个节点(0-indexed)
//
//	返回跳到的节点,如果不存在这样的节点,返回-1
func (hld *LCAFast) Jump(from, to, step int) int {
	if step == 1 {
		if from == to {
			return -1
		}
		if hld.IsInSubtree(to, from) {
			return hld.KthAncestor(to, int(hld.Depth[to]-hld.Depth[from]-1))
		}
		return int(hld.Parent[from])
	}
	c := hld.LCA(from, to)
	dac := hld.Depth[from] - hld.Depth[c]
	dbc := hld.Depth[to] - hld.Depth[c]
	if step > int(dac+dbc) {
		return -1
	}
	if step <= int(dac) {
		return hld.KthAncestor(from, step)
	}
	return hld.KthAncestor(to, int(dac+dbc-int32(step)))
}

// child 是否在 root 的子树中 (child和root不能相等)
func (hld *LCAFast) IsInSubtree(child, root int) bool {
	return hld.lid[root] <= hld.lid[child] && hld.lid[child] < hld.rid[root]
}

func (hld *LCAFast) _build(cur, pre, dep int32) int {
	subSize, heavySize, heavySon := 1, 0, int32(-1)
	for _, next := range hld.Tree[cur] {
		if next != pre {
			nextSize := hld._build(next, cur, dep+1)
			subSize += nextSize
			if nextSize > heavySize {
				heavySize, heavySon = nextSize, next
			}
		}
	}
	hld.Depth[cur] = dep
	hld.heavySon[cur] = heavySon
	hld.Parent[cur] = pre
	return subSize
}

func (hld *LCAFast) _markTop(cur, top int32) {
	hld.top[cur] = top
	hld.lid[cur] = hld.dfnId
	hld.idToNode[hld.dfnId] = cur
	hld.dfnId++
	if hld.heavySon[cur] != -1 {
		hld._markTop(hld.heavySon[cur], top)
		for _, next := range hld.Tree[cur] {
			if next != hld.heavySon[cur] && next != hld.Parent[cur] {
				hld._markTop(next, next)
			}
		}
	}
	hld.rid[cur] = hld.dfnId
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
