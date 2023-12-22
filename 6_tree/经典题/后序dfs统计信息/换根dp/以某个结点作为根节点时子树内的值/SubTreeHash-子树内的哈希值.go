// SubTreeHash-子树内的哈希值

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"time"
)

func main() {
	// https://judge.yosupo.jp/problem/rooted_tree_isomorphism_classification
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	parents := make([]int, n-1)
	for i := range parents {
		fmt.Fscan(in, &parents[i])
	}
	H := NewSubTreeHash(n)
	for i := 1; i < n; i++ {
		H.AddEdge(parents[i-1], i, 0)
	}
	H.Build()

	res := make([]int, n)
	for i := range res {
		res[i] = int(H.GetSubTree(0, i))
	}

	sorted, rank := sortedSet(res)
	for i := range res {
		res[i] = rank[res[i]]
	}
	fmt.Fprintln(out, len(sorted))
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

func sortedSet(xs []int) (sorted []int, rank map[int]int) {
	set := make(map[int]struct{}, len(xs))
	for _, v := range xs {
		set[v] = struct{}{}
	}
	sorted = make([]int, 0, len(set))
	for k := range set {
		sorted = append(sorted, k)
	}
	sort.Ints(sorted)
	rank = make(map[int]int, len(sorted))
	for i, v := range sorted {
		rank[v] = i
	}
	return
}

var SEED uint = uint(time.Now().UnixNano()/2 + 1)

func nextRand(seed *uint) uint {
	*seed ^= *seed << 13
	*seed ^= *seed >> 17
	*seed ^= *seed << 5
	return *seed
}

type SubTreeHash struct {
	rerooting *_ReRootingEdge
	dp        []uint
	dp1       []uint
	dp2       []uint
	n         int
}

func NewSubTreeHash(n int) *SubTreeHash {
	return &SubTreeHash{rerooting: _NewReRootingEdge(n), n: n}
}

func (sd *SubTreeHash) AddEdge(from, to, weight int) {
	sd.rerooting.tree.AddEdge(from, to, weight)
}

func (sd *SubTreeHash) Build() {
	n := sd.n
	hashBases := make([]uint, n)
	for i := range hashBases {
		hashBases[i] = nextRand(&SEED)
	}

	unit := func() E { return E{0, 1} }
	op := func(dp1, dp2 E) E {
		return E{maxUint(dp1[0], dp2[0]), dp1[1] * dp2[1]}
	}
	composition := func(dp E, from int) E {
		return E{dp[0] + 1, dp[1]}
	}
	compositionEdge := func(dp E, edge Edge) E {
		return E{dp[0], dp[1] + hashBases[dp[0]]}
	}

	sd.rerooting.ReRooting(unit, op, composition, compositionEdge)
	sd.dp = make([]uint, n)
	sd.dp1 = make([]uint, n)
	sd.dp2 = make([]uint, n)
	for v := 0; v < n; v++ {
		sd.dp[v] = sd.rerooting.dp[v][1]
		sd.dp1[v] = sd.rerooting.dp1[v][1]
		sd.dp2[v] = sd.rerooting.dp2[v][1]
	}
}

// 以 root 为根的子树的哈希值.
func (sd *SubTreeHash) Get(root int) uint {
	return sd.dp[root]
}

// 以 root 为根时, 子树 v 的哈希值.
func (sd *SubTreeHash) GetSubTree(root, v int) uint {
	if root == v {
		return sd.dp[v]
	}
	if !sd.rerooting.tree.IsInSubtree(root, v) {
		return sd.dp1[v]
	}
	w := sd.rerooting.tree.Jump(v, root, 1)
	return sd.dp2[w]
}

type E = [2]uint

type Edge = struct{ from, to, weight, eid int }

type _ReRootingEdge struct {
	tree *_T
	dp1  []E // 边 parent-root 的子树 root 的 dp值
	dp2  []E // 边 parent-root 的子树 parent 的 dp值
	dp   []E // 顶点 v 的子树的 dp值
}

func _NewReRootingEdge(n int) *_ReRootingEdge {
	return &_ReRootingEdge{tree: _NT(n)}
}

func (rr *_ReRootingEdge) AddEdge(from, to, weight int) {
	rr.tree.AddEdge(from, to, weight)
}

// !root 作为根节点时, 子树 v 的 dp 值
func (rr *_ReRootingEdge) Get(root, v int) E {
	if root == v {
		return rr.dp[v]
	}
	if !rr.tree.IsInSubtree(root, v) {
		return rr.dp1[v]
	}
	w := rr.tree.Jump(v, root, 1)
	return rr.dp2[w]
}

func (rr *_ReRootingEdge) ReRooting(
	e func() E,
	op func(dp1, dp2 E) E,
	composition func(dp E, from int) E,
	compositionEdge func(dp E, edge Edge) E,
) []E {
	rr.tree.Build(-1)
	unit := e()
	N := len(rr.tree.Tree)
	dp1, dp2, dp := make([]E, N), make([]E, N), make([]E, N)
	for i := 0; i < N; i++ {
		dp1[i] = unit
		dp2[i] = unit
		dp[i] = unit
	}

	V := rr.tree.IdToNode
	par := rr.tree.Parent
	for i := N - 1; i >= 0; i-- {
		v := V[i]
		ch := rr.tree.CollectChild(v)
		n := len(ch)
		x1, x2 := make([]E, n+1), make([]E, n+1)
		for i := range x1 {
			x1[i] = unit
			x2[i] = unit
		}
		for i := 0; i < n; i++ {
			x1[i+1] = op(x1[i], dp2[ch[i]])
		}
		for i := n - 1; i >= 0; i-- {
			x2[i] = op(dp2[ch[i]], x2[i+1])
		}
		for i := 0; i < n; i++ {
			dp2[ch[i]] = op(x1[i], x2[i+1])
		}
		dp[v] = x2[0]
		dp1[v] = composition(dp[v], v)
		for _, e := range rr.tree.Tree[v] {
			to := e.to
			if to == par[v] {
				dp2[v] = compositionEdge(dp1[v], e)
			}
		}
	}

	v := V[0]
	dp[v] = composition(dp[v], v)
	for _, e := range rr.tree.Tree[v] {
		to := e.to
		dp2[to] = composition(dp2[to], v)
	}

	for i := 0; i < N; i++ {
		v := V[i]
		for _, e := range rr.tree.Tree[v] {
			if e.to == par[v] {
				continue
			}
			x := compositionEdge(dp2[e.to], e)
			for _, f := range rr.tree.Tree[e.to] {
				if f.to == par[e.to] {
					continue
				}
				dp2[f.to] = op(dp2[f.to], x)
				dp2[f.to] = composition(dp2[f.to], e.to)
			}
			x = op(dp[e.to], x)
			dp[e.to] = composition(x, e.to)
		}
	}

	rr.dp1, rr.dp2, rr.dp = dp1, dp2, dp
	return dp
}

type _T struct {
	Tree          [][]Edge
	Depth         []int
	Parent        []int
	LID, RID      []int // 欧拉序[in,out)
	IdToNode      []int
	top, heavySon []int
	timer         int
	eid           int
}

func _NT(n int) *_T {
	tree := make([][]Edge, n)
	lid := make([]int, n)
	rid := make([]int, n)
	IdToNode := make([]int, n)
	top := make([]int, n)      // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	depth := make([]int, n)    // 深度
	parent := make([]int, n)   // 父结点
	heavySon := make([]int, n) // 重儿子
	for i := range parent {
		parent[i] = -1
	}

	return &_T{
		Tree:     tree,
		Depth:    depth,
		Parent:   parent,
		LID:      lid,
		RID:      rid,
		IdToNode: IdToNode,
		top:      top,
		heavySon: heavySon,
	}
}

// 添加无向边 u-v, 边权为w.
func (tree *_T) AddEdge(u, v, w int) {
	tree.Tree[u] = append(tree.Tree[u], Edge{from: u, to: v, weight: w, eid: tree.eid})
	tree.Tree[v] = append(tree.Tree[v], Edge{from: v, to: u, weight: w, eid: tree.eid})
	tree.eid++
}

// root:0-based
//
//	当root设为-1时，会从0开始遍历未访问过的连通分量
func (tree *_T) Build(root int) {
	if root != -1 {
		tree.build(root, -1, 0)
		tree.markTop(root, root)
	} else {
		for i := 0; i < len(tree.Tree); i++ {
			if tree.Parent[i] == -1 {
				tree.build(i, -1, 0)
				tree.markTop(i, i)
			}
		}
	}
}

func (tree *_T) LCA(u, v int) int {
	for {
		if tree.LID[u] > tree.LID[v] {
			u, v = v, u
		}
		if tree.top[u] == tree.top[v] {
			return u
		}
		v = tree.Parent[tree.top[v]]
	}
}

// k: 0-based
//
//	如果不存在第k个祖先，返回-1
func (tree *_T) KthAncestor(root, k int) int {
	if k > tree.Depth[root] {
		return -1
	}
	for {
		u := tree.top[root]
		if tree.LID[root]-k >= tree.LID[u] {
			return tree.IdToNode[tree.LID[root]-k]
		}
		k -= tree.LID[root] - tree.LID[u] + 1
		root = tree.Parent[u]
	}
}

// 从 from 节点跳向 to 节点,跳过 step 个节点(0-indexed)
//
//	返回跳到的节点,如果不存在这样的节点,返回-1
func (tree *_T) Jump(from, to, step int) int {
	if step == 1 {
		if from == to {
			return -1
		}
		if tree.IsInSubtree(to, from) {
			return tree.KthAncestor(to, tree.Depth[to]-tree.Depth[from]-1)
		}
		return tree.Parent[from]
	}
	c := tree.LCA(from, to)
	dac := tree.Depth[from] - tree.Depth[c]
	dbc := tree.Depth[to] - tree.Depth[c]
	if step > dac+dbc {
		return -1
	}
	if step <= dac {
		return tree.KthAncestor(from, step)
	}
	return tree.KthAncestor(to, dac+dbc-step)
}

func (tree *_T) CollectChild(root int) []int {
	res := []int{}
	for _, e := range tree.Tree[root] {
		next := e.to
		if next != tree.Parent[root] {
			res = append(res, next)
		}
	}
	return res
}

// child 是否在 root 的子树中 (child和root不能相等)
func (tree *_T) IsInSubtree(child, root int) bool {
	return tree.LID[root] <= tree.LID[child] && tree.LID[child] < tree.RID[root]
}

func (tree *_T) build(cur, pre, dep int) int {
	subSize, heavySize, heavySon := 1, 0, -1
	for _, e := range tree.Tree[cur] {
		next := e.to
		if next != pre {
			nextSize := tree.build(next, cur, dep+1)
			subSize += nextSize
			if nextSize > heavySize {
				heavySize, heavySon = nextSize, next
			}
		}
	}
	tree.Depth[cur] = dep
	tree.heavySon[cur] = heavySon
	tree.Parent[cur] = pre
	return subSize
}

func (tree *_T) markTop(cur, top int) {
	tree.top[cur] = top
	tree.LID[cur] = tree.timer
	tree.IdToNode[tree.timer] = cur
	tree.timer++
	if tree.heavySon[cur] != -1 {
		tree.markTop(tree.heavySon[cur], top)
		for _, e := range tree.Tree[cur] {
			next := e.to
			if next != tree.heavySon[cur] && next != tree.Parent[cur] {
				tree.markTop(next, next)
			}
		}
	}
	tree.RID[cur] = tree.timer
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

func maxUint(a, b uint) uint {
	if a > b {
		return a
	}
	return b
}
