// BlockCutTree-圆方树点双
// 处理点双、割点相关逻辑直接使用圆方树即可
// https://oi-wiki.org/graph/block-forest/
// https://oi-wiki.org/graph/images/block-forest2.svg
// 如果原图连通，则「圆方树」才是一棵树，如果原图有 k 个连通分量，则它的圆方树也会形成 k 棵树形成的森林。
//
// 例如:
//
//	     原图        圆方树
//			0  —  1     0     1
//			 \	 /       \	 /
//			  \ /   =>     5 (block)
//			   2           |
//			   |           2 (原图割点,与block相连)
//			   3           |
//		                 4 (block)
//		                 |
//		                 3
//
// !割点：到达两个区域的必经之点.
// 1. 原图的割点`至少`在两个不同的 v-BCC 中
// 2. 原图不是割点的点都`只存在`于一个 v-BCC 中
// 3. v-BCC 形成的子图内没有割点
// !4. 对于圆方树上一条从 u 到 v 的简单路径，设为 u -> s1 -> c1 -> s2 -> c2 -> ... -> sk -> v，其中s为方点，c为圆点。
// !   那么有：u->v 所有路径上点的并集恰好是 s1, s2, ..., sk 的并集。

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	abc318g()
	// yosupo()
}

// G - Typical Path Problem
// https://atcoder.jp/contests/abc318/tasks/abc318_g
// 给定一张无向图以及a,b,c三个顶点.
// 问是由存在一条从a到b，且经过c的简单路径.
//
// https://www.cnblogs.com/xrkforces/p/ABC318G.html
// 因为 a->c 经过的点是a->c路径上所有方点的并集.
// !因此，只需满足：b所在的任意一个方点在a->c路径上即可.
func abc318g() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int32
	fmt.Fscan(in, &n, &m)
	var a, b, c int32
	fmt.Fscan(in, &a, &b, &c)
	a, b, c = a-1, b-1, c-1
	graph := make([][]int32, n)
	for i := int32(0); i < m; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		graph[u] = append(graph[u], v)
		graph[v] = append(graph[v], u)
	}

	blockCutTree := BlockCutTree32(graph)
	newTree := NewCompressedBinaryLiftFromTree(blockCutTree, 0)
	pathAC := newTree.CreatePath(a, c)
	belong := GetBelong(n, blockCutTree)
	for _, v := range belong[b] {
		if pathAC.Contains(v) {
			fmt.Fprintln(out, "Yes")
			return
		}
	}
	fmt.Fprintln(out, "No")
}

// https://judge.yosupo.jp/problem/biconnected_components
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int32
	fmt.Fscan(in, &n, &m)
	graph := make([][]int32, n)
	for i := int32(0); i < m; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		graph[u] = append(graph[u], v)
		graph[v] = append(graph[v], u)
	}

	tree := BlockCutTree32(graph)
	fmt.Fprintln(out, int32(len(tree))-n) // block数量=点双连通分量数量
	for i := n; i < int32(len(tree)); i++ {
		fmt.Fprint(out, len(tree[i]))
		for _, v := range tree[i] {
			fmt.Fprint(out, " ", v)
		}
		fmt.Fprintln(out)
	}
}

// 原图的每个点所属的点双连通分量(每个点所连接的方点).
// 如果所属>=2个方点,则说明是原图的割点.
func GetBelong(rawN int32, blockTree [][]int32) [][]int32 {
	belong := make([][]int32, rawN)
	for i := rawN; i < int32(len(blockTree)); i++ {
		for _, v := range blockTree[i] {
			belong[v] = append(belong[v], i)
		}
	}
	return belong
}

// 求出无向图的blockCutTree, 用于解决点双连通分量相关问题.
// 在blockCutTree中, 满足性质:
// !1.[0, n)为原图中的点, [n, n+n_block)为block.每一个点双连通分量连接、对应一个block(方点).
// !2.割点 <=> [0, n)中满足degree>=2的点.
func BlockCutTree32(graph [][]int32) (tree [][]int32) {
	n := int32(len(graph))
	low := make([]int32, n)
	order := make([]int32, n)
	stack := make([]int32, 0, n)
	used := make([]bool, n)
	id := n
	now := int32(0)
	edges := [][2]int32{}

	var dfs func(int32, int32)
	dfs = func(cur, pre int32) {
		stack = append(stack, cur)
		used[cur] = true
		low[cur] = now
		order[cur] = now
		now++
		child := 0
		for _, to := range graph[cur] {
			if to == pre {
				continue
			}
			if !used[to] {
				child++
				s := len(stack)
				dfs(int32(to), cur)
				low[cur] = min32(low[cur], low[to])
				if (pre == -1 && child > 1) || (pre != -1 && low[to] >= order[cur]) {
					edges = append(edges, [2]int32{id, cur})
					for len(stack) > s {
						edges = append(edges, [2]int32{id, stack[len(stack)-1]})
						stack = stack[:len(stack)-1]
					}
					id++
				}
			} else {
				low[cur] = min32(low[cur], order[to])
			}
		}
	}

	for i := int32(0); i < n; i++ {
		if !used[i] {
			dfs(i, -1)
			for _, v := range stack {
				edges = append(edges, [2]int32{id, v})
			}
			id++
			stack = stack[:0]
		}
	}

	tree = make([][]int32, id)
	for _, e := range edges {
		u, v := e[0], e[1]
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}
	return
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a <= b {
		return a
	}
	return b
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

func (lca *CompressedBinaryLift) InSubtree(maybeChild, maybeAncestor int32) bool {
	return lca.Depth[maybeChild] >= lca.Depth[maybeAncestor] &&
		lca.KthAncestor(maybeChild, lca.Depth[maybeChild]-lca.Depth[maybeAncestor]) == maybeAncestor
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

func (t *TreePath) Contains(node int32) bool {
	lcaFn := t.lcaFn
	return lcaFn(node, t.Lca) == t.Lca &&
		(lcaFn(node, t.From) == node || lcaFn(node, t.To) == node)
}

func (t *TreePath) HasIntersection(other *TreePath) bool {
	return t.Contains(other.Lca) || other.Contains(t.Lca)
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

// 将路径以separator为分隔，按顺序分成两段.separtor必须在路径上.
func (t *TreePath) Split(separator int32) (path1, path2 *TreePath) {
	down, top := t.From, t.To
	if down == top {
		return nil, nil
	}

	swapped := false
	if t.depth[down] < t.depth[top] {
		down, top = top, down
		swapped = true
	}

	from1, to1, from2, to2 := int32(-1), int32(-1), int32(-1), int32(-1)
	if t.Lca == top {
		// down和top在一条链上.
		if separator == down {
			from2 = t.kthAncestorFn(separator, 1)
			to2 = top
		} else if separator == top {
			from1 = down
			to1 = t.kthAncestorFn(down, t.depth[down]-t.depth[separator]-1)
		} else {
			from1 = down
			to1 = t.kthAncestorFn(down, t.depth[down]-t.depth[separator]-1)
			from2 = t.kthAncestorFn(separator, 1)
			to2 = top
		}
	} else {
		// down和top在lca两个子树上.
		if separator == down {
			from2 = t.kthAncestorFn(separator, 1)
			to2 = top
		} else if separator == top {
			from1 = down
			to1 = t.kthAncestorFn(separator, 1)
		} else {
			var jump1, jump2 int32
			if separator == t.Lca {
				jump1 = t.kthAncestorFn(down, t.depth[down]-t.depth[separator]-1)
				jump2 = t.kthAncestorFn(top, t.depth[top]-t.depth[separator]-1)
			} else if t.lcaFn(separator, down) == separator {
				jump1 = t.kthAncestorFn(down, t.depth[down]-t.depth[separator]-1)
				jump2 = t.kthAncestorFn(separator, 1)
			} else {
				jump1 = t.kthAncestorFn(separator, 1)
				jump2 = t.kthAncestorFn(top, t.depth[top]-t.depth[separator]-1)
			}
			from1 = down
			to1 = jump1
			from2 = jump2
			to2 = top
		}
	}

	if swapped {
		from1, to1, from2, to2 = to2, from2, to1, from1
	}
	if from1 != -1 && to1 != -1 {
		path1 = NewTreePath(from1, to1, t.depth, t.kthAncestorFn, t.lcaFn)
	}
	if from2 != -1 && to2 != -1 {
		path2 = NewTreePath(from2, to2, t.depth, t.kthAncestorFn, t.lcaFn)
	}
	return
}

func (t *TreePath) Len() int32 {
	return t.depth[t.From] + t.depth[t.To] - 2*t.depth[t.Lca]
}
