package main

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
