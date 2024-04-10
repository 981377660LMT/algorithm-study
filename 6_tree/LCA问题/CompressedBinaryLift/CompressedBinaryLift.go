package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	// yosupo()
	// P3398()
	// jump()

	CF519E()
}

/**
 * Your DistanceLimitedPathsExist object will be instantiated and called as such:
 * obj := Constructor(n, edgeList);
 * param_1 := obj.Query(p,q,limit);
 */
// https://judge.yosupo.jp/problem/lca
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	tree := make([][]int32, n)
	for i := 1; i < n; i++ {
		var parent int32
		fmt.Fscan(in, &parent)
		tree[parent] = append(tree[parent], int32(i))
	}
	bl := NewCompressedBinaryLiftFromTree(tree, 0)
	for i := 0; i < q; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		fmt.Fprintln(out, bl.Lca(u, v))
	}
}

func jump() {
	// https://judge.yosupo.jp/problem/jump_on_tree
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)

	tree := make([][]int32, n)
	for i := int32(0); i < n-1; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}
	D := NewCompressedBinaryLiftFromTree(tree, 0)

	for i := int32(0); i < q; i++ {
		var from, to, k int32
		fmt.Fscan(in, &from, &to, &k)
		fmt.Fprintln(out, D.Jump(from, to, k))
	}
}

// https://www.luogu.com.cn/problem/P3398
func P3398() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

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
	bl := NewCompressedBinaryLiftFromTree(tree, 0)
	query := func(a, b int32, c, d int32) bool {
		path1 := bl.CreatePath(a, b)
		path2 := bl.CreatePath(c, d)
		return path1.HasIntersection(path2)
	}
	for i := int32(0); i < q; i++ {
		var a, b, c, d int32
		fmt.Fscan(in, &a, &b, &c, &d)
		a, b, c, d = a-1, b-1, c-1, d-1
		if query(a, b, c, d) {
			fmt.Fprintln(out, "Y")
		} else {
			fmt.Fprintln(out, "N")
		}
	}
}

// A and B and Lecture Rooms (到树上两点距离相等的点的个数)
// https://www.luogu.com.cn/problem/CF519E
// 给定一棵n个点的无根树和q组询问，对于每一组询问，求出树上到某两点距离相等的点数（包含本身）。
func CF519E() {
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

	bl := NewCompressedBinaryLiftFromTree(tree, 0)
	subSize := make([]int32, n)
	var dfs func(v, p int32)
	dfs = func(v, p int32) {
		subSize[v] = 1
		for _, to := range tree[v] {
			if to == p {
				continue
			}
			dfs(to, v)
			subSize[v] += subSize[to]
		}
	}
	dfs(0, -1)

	query := func(a, b int32) int32 {
		if a == b {
			return n
		}
		lca := bl.Lca(a, b)
		len_ := bl.Depth[a] + bl.Depth[b] - 2*bl.Depth[lca] + 1
		if len_&1 == 0 {
			return 0
		}
		if bl.Depth[a] < bl.Depth[b] {
			a, b = b, a
		}
		// lca是否为中点
		halfLen := len_ / 2
		center := bl.KthAncestor(a, halfLen)
		if lca == center {
			p1, p2 := bl.KthAncestor(a, halfLen-1), bl.KthAncestor(b, halfLen-1)
			return n - subSize[p1] - subSize[p2]
		} else {
			p1 := bl.KthAncestor(a, halfLen-1)
			return subSize[center] - subSize[p1]
		}
	}

	var q int32
	fmt.Fscan(in, &q)
	for i := int32(0); i < q; i++ {
		var x, y int32
		fmt.Fscan(in, &x, &y)
		x, y = x-1, y-1
		fmt.Fprintln(out, query(x, y))
	}
}

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

// root:-1表示无根.
func NewCompressedBinaryLiftFromTree(tree [][]int32, root int32) *CompressedBinaryLift {
	n := int32(len(tree))
	res := &CompressedBinaryLift{
		Depth:  make([]int32, n),
		Parent: make([]int32, n),
		jump:   make([]int32, n),
	}
	if root != -1 {
		res.Parent[root] = -1
		res.jump[root] = root
		res._setUp(tree, root)
	} else {
		for i := int32(0); i < n; i++ {
			res.Parent[i] = -1
		}
		for i := int32(0); i < n; i++ {
			if res.Parent[i] == -1 {
				res.jump[i] = i
				res._setUp(tree, i)
			}
		}
	}
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

func (lca *CompressedBinaryLift) Jump(start, target, step int32) int32 {
	lca_ := lca.Lca(start, target)
	dep1, dep2, deplca := lca.Depth[start], lca.Depth[target], lca.Depth[lca_]
	dist := dep1 + dep2 - 2*deplca
	if step > dist {
		return -1
	}
	if step <= dep1-deplca {
		return lca.KthAncestor(start, step)
	}
	return lca.KthAncestor(target, dist-step)
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

func (bl *CompressedBinaryLift) _setUp(tree [][]int32, cur int32) {
	queue := []int32{cur}
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
