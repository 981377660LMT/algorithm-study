package main

// 空间复杂度`O(n)`的树上倍增.
//   - https://taodaling.github.io/blog/2020/03/18/binary-lifting/
//   - https://codeforces.com/blog/entry/74847
//   - https://codeforces.com/blog/entry/100826
type CompressedBinaryLift struct {
	Depth  []int32
	Parent []int32
	jump   []int32 // 指向当前节点的某个祖先节点.
}

func NewCompressedBinaryLift(n int32, depthOnTree, parentOnTree []int32) *CompressedBinaryLift {
	res := &CompressedBinaryLift{
		Depth:  depthOnTree,
		Parent: parentOnTree,
		jump:   make([]int32, n),
	}
	for i := int32(0); i < n; i++ {
		res.jump[i] = -1
	}
	for i := int32(0); i < n; i++ {
		res._consider(i)
	}
	return res
}

func NewCompressedBinaryLiftFromTree(tree [][]int32, root int32) *CompressedBinaryLift {
	n := int32(len(tree))
	res := &CompressedBinaryLift{
		Depth:  make([]int32, n),
		Parent: make([]int32, n),
		jump:   make([]int32, n),
	}
	res.Parent[root] = -1
	res.jump[root] = root
	res._setUp(tree, root)
	return res
}

func (bl *CompressedBinaryLift) FirstTrue(start int32, predicate func(end int32) bool) int32 {
	for !predicate(start) {
		if predicate(bl.jump[start]) {
			start = bl.Parent[start]
		} else {
			if start == bl.jump[start] {
				return -1
			}
			start = bl.jump[start]
		}
	}
	return start
}

func (bl *CompressedBinaryLift) LastTrue(start int32, predicate func(end int32) bool) int32 {
	if !predicate(start) {
		return -1
	}
	for {
		if predicate(bl.jump[start]) {
			if start == bl.jump[start] {
				return start
			}
			start = bl.jump[start]
		} else if predicate(bl.Parent[start]) {
			start = bl.Parent[start]
		} else {
			return start
		}
	}
}

func (bl *CompressedBinaryLift) UpToDepth(root int32, toDepth int32) int32 {
	if !(0 <= toDepth && toDepth <= bl.Depth[root]) {
		return -1
	}
	for bl.Depth[root] > toDepth {
		if bl.Depth[bl.jump[root]] < toDepth {
			root = bl.Parent[root]
		} else {
			root = bl.jump[root]
		}
	}
	return root
}

func (bl *CompressedBinaryLift) KthAncestor(node, k int32) int32 {
	targetDepth := bl.Depth[node] - k
	return bl.UpToDepth(node, targetDepth)
}

func (bl *CompressedBinaryLift) Lca(a, b int32) int32 {
	if bl.Depth[a] > bl.Depth[b] {
		a = bl.KthAncestor(a, bl.Depth[a]-bl.Depth[b])
	} else if bl.Depth[a] < bl.Depth[b] {
		b = bl.KthAncestor(b, bl.Depth[b]-bl.Depth[a])
	}
	for a != b {
		if bl.jump[a] == bl.jump[b] {
			a = bl.Parent[a]
			b = bl.Parent[b]
		} else {
			a = bl.jump[a]
			b = bl.jump[b]
		}
	}
	return a
}

func (bl *CompressedBinaryLift) Dist(a, b int32) int32 {
	return bl.Depth[a] + bl.Depth[b] - 2*bl.Depth[bl.Lca(a, b)]
}

func (bl *CompressedBinaryLift) CreatePath(from, to int32) *TreePath {
	return NewTreePath(from, to, bl.Depth, bl.KthAncestor, bl.Lca)
}

func (bl *CompressedBinaryLift) _consider(root int32) {
	if root == -1 || bl.jump[root] != -1 {
		return
	}
	p := bl.Parent[root]
	bl._consider(p)
	bl._addLeaf(root, p)
}

func (bl *CompressedBinaryLift) _addLeaf(leaf, parent int32) {
	if parent == -1 {
		bl.jump[leaf] = leaf
	} else if tmp := bl.jump[parent]; bl.Depth[parent]-bl.Depth[tmp] == bl.Depth[tmp]-bl.Depth[bl.jump[tmp]] {
		bl.jump[leaf] = bl.jump[tmp]
	} else {
		bl.jump[leaf] = parent
	}
}

func (bl *CompressedBinaryLift) _setUp(tree [][]int32, root int32) {
	queue := []int32{root}
	head := 0
	for head < len(queue) {
		cur := queue[head]
		head++
		nexts := tree[cur]
		for _, next := range nexts {
			if next == bl.Parent[cur] {
				continue
			}
			bl.Depth[next] = bl.Depth[cur] + 1
			bl.Parent[next] = cur
			queue = append(queue, next)
			bl._addLeaf(next, cur)
		}
	}
}

type TreePath struct {
	From, To      int32
	Lca           int32
	depth         []int32
	kthAncestorFn func(node, k int32) int32
	lcaFn         func(node1, node2 int32) int32
}

func NewTreePath(
	from, to int32,
	depth []int32, kthAncestorFn func(node, k int32) int32, lcaFn func(node1, node2 int32) int32,
) *TreePath {
	return &TreePath{
		From: from, To: to, Lca: lcaFn(from, to),
		depth: depth, kthAncestorFn: kthAncestorFn, lcaFn: lcaFn,
	}
}

// 从路径的起点开始，第k个节点(0-indexed).不存在则返回-1.
func (t *TreePath) KthNodeOnPath(k int32) int32 {
	if k <= t.depth[t.From]-t.depth[t.Lca] {
		return t.kthAncestorFn(t.From, k)
	}
	return t.kthAncestorFn(t.To, t.Len()-k)
}

func (t *TreePath) OnPath(node int32) bool {
	lcaFn := t.lcaFn
	return lcaFn(node, t.Lca) == t.Lca &&
		(lcaFn(node, t.From) == node || lcaFn(node, t.To) == node)
}

func (t *TreePath) HasIntersection(other *TreePath) bool {
	return t.OnPath(other.Lca) || other.OnPath(t.Lca)
}

// 求两条路径的交, 返回相交线段的两个端点.无交点则返回(-1, -1, false).
func (t *TreePath) GetIntersection(other *TreePath) (p1, p2 int32, ok bool) {
	a, b, c, d := t.From, t.To, other.From, other.To
	lcaFn, depth := t.lcaFn, t.depth
	x1, x2, x3, x4 := lcaFn(a, c), lcaFn(a, d), lcaFn(b, c), lcaFn(b, d)
	p1, p2 = x1, x2
	if depth[x2] > depth[p1] {
		p2 = p1
		p1 = x2
	}
	update := func(x int32) {
		curDepth := depth[x]
		if curDepth > depth[p1] {
			p2 = p1
			p1 = x
		} else if curDepth > depth[p2] {
			p2 = x
		}
	}
	update(x3)
	update(x4)
	lca1, lca2 := t.Lca, other.Lca
	if p1 != p2 {
		return p1, p2, true
	}
	if depth[p1] < depth[lca1] || depth[p1] < depth[lca2] {
		return -1, -1, false
	}
	return p1, p2, true
}

func (t *TreePath) CountIntersection(other *TreePath) int32 {
	p1, p2, ok := t.GetIntersection(other)
	if !ok {
		return 0
	}
	if p1 == p2 {
		return 1
	}
	return t.depth[p1] + t.depth[p2] - 2*t.depth[t.Lca] + 1
}

func (t *TreePath) Len() int32 {
	return t.depth[t.From] + t.depth[t.To] - 2*t.depth[t.Lca]
}
