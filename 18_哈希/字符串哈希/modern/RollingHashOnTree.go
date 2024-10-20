// !树上路径哈希(TreePathHash)

package main

import (
	"fmt"
	"math/rand"
)

func main() {
	//   0
	//  / \
	// 1   2
	//    / \
	//   3   4

	tree := NewTree32(5)
	tree.AddEdge(0, 1, 1)
	tree.AddEdge(0, 2, 1)
	tree.AddEdge(2, 3, 1)
	tree.AddEdge(2, 4, 1)
	tree.Build(0)
	r := NewRollingHashOnTree(tree, 0, false, func(viOrEi int32) uint64 { return uint64(viOrEi % 2) })
	fmt.Println(r.Get(0, 1))                    // 1
	fmt.Println(r.Get(2, 3))                    // 1
	fmt.Println(r.LcpAndCompareStr(0, 3, 0, 4)) // 2, 1
}

type RollingHashOnTree struct {
	n         int32
	tree      *Tree32
	base      uint64
	isEdge    bool
	data      []uint64
	pow, ipow []uint64
	dp1, dp2  []uint64
}

// !静态树上路径哈希(TreePathHash).
//
//	base: 0 表示随机生成.
//	isEdge: 边权还是点权.
//	f: 第 i 个点或边的权值.
func NewRollingHashOnTree(tree *Tree32, base uint64, isEdge bool, f func(viOrEi int32) uint64) *RollingHashOnTree {
	if base == 0 {
		base = uint64(37 + rand.Intn(1e9))
	}
	res := &RollingHashOnTree{n: tree.n, tree: tree, base: base, isEdge: isEdge}
	res.Build(f)
	return res
}

func (r *RollingHashOnTree) Build(f func(vidOrEid int32) uint64) {
	n, tree, base := r.n, r.tree, r.base
	data := make([]uint64, n)
	if r.isEdge {
		for i := int32(0); i < n-1; i++ {
			data[tree.EToV(i)] = f(i)
		}
	} else {
		for i := int32(0); i < n; i++ {
			data[i] = f(i)
		}
	}
	pow, ipow := make([]uint64, n+1), make([]uint64, n+1)
	pow[0], pow[1] = 1, base
	ipow[0], ipow[1] = 1, modInv(base)
	for i := int32(2); i <= n; i++ {
		pow[i] = modMul(pow[i-1], base)
		ipow[i] = modMul(ipow[i-1], ipow[1])
	}
	root := tree.IdToNode[0]
	dp1, dp2 := make([]uint64, n), make([]uint64, n)
	dp1[root], dp2[root] = data[0], data[0]
	for i := int32(1); i < n; i++ {
		v := tree.IdToNode[i]
		d, p := tree.Depth[v], tree.Parent[v]
		dp1[v] = modAdd(modMul(base, dp1[p]), data[v])
		dp2[v] = modAdd(dp2[p], modMul(pow[d], data[v]))
	}
	r.data, r.pow, r.ipow, r.dp1, r.dp2 = data, pow, ipow, dp1, dp2
}

// a到b路径的哈希值.
func (r *RollingHashOnTree) Get(a, b int32) uint64 {
	c := r.tree.Lca(a, b)
	x1, x2 := r.getDu(a, c), r.getUd(c, b)
	n2 := r.tree.Depth[b] - r.tree.Depth[c]
	if !r.isEdge {
		x1 = modAdd(modMul(x1, r.base), r.data[c])
	}
	return modAdd(modMul(x1, r.pow[n2]), x2)
}

func (r *RollingHashOnTree) Lcp(s1, t1, s2, t2 int32) int32 {
	a, _ := r.LcpAndCompareStr(s1, t1, s2, t2)
	return a
}

func (r *RollingHashOnTree) CompareStr(s1, t1, s2, t2 int32) int8 {
	_, b := r.LcpAndCompareStr(s1, t1, s2, t2)
	return b
}

func (r *RollingHashOnTree) LcpAndCompareStr(s1, t1, s2, t2 int32) (int32, int8) {
	var lcp int32
	path1, path2 := r.tree.GetPathDecomposition(s1, t1, !r.isEdge), r.tree.GetPathDecomposition(s2, t2, !r.isEdge)
	for i, j := 0, len(path1)-1; i < j; i, j = i+1, j-1 {
		path1[i], path1[j] = path1[j], path1[i]
	}
	for i, j := 0, len(path2)-1; i < j; i, j = i+1, j-1 {
		path2[i], path2[j] = path2[j], path2[i]
	}
	for len(path1) > 0 && len(path2) > 0 {
		p1 := path1[len(path1)-1]
		path1 = path1[:len(path1)-1]
		p2 := path2[len(path2)-1]
		path2 = path2[:len(path2)-1]
		a, b := p1[0], p1[1]
		c, d := p2[0], p2[1]
		n1, n2 := abs32(a-b)+1, abs32(c-d)+1
		n := min32(n1, n2)
		if n < n1 {
			if a <= b {
				path1 = append(path1, [2]int32{a + n, b})
				b = a + n - 1
			} else {
				path1 = append(path1, [2]int32{a - n, b})
				b = a - n + 1
			}
		}
		if n < n2 {
			if c <= d {
				path2 = append(path2, [2]int32{c + n, d})
				d = c + n - 1
			} else {
				path2 = append(path2, [2]int32{c - n, d})
				d = c - n + 1
			}
		}
		x1, x2 := r.fromHldPair(a, b), r.fromHldPair(c, d)
		if x1 == x2 {
			lcp += n
			continue
		}
		check := func(n int32) bool {
			if n == 0 {
				return true
			}
			var x1 uint64
			if a <= b {
				x1 = r.fromHldPair(a, a+n-1)
			} else {
				x1 = r.fromHldPair(a, a-n+1)
			}
			var x2 uint64
			if c <= d {
				x2 = r.fromHldPair(c, c+n-1)
			} else {
				x2 = r.fromHldPair(c, c-n+1)
			}
			return x1 == x2
		}
		k := binarySearch32(check, 0, n)
		lcp += k
		if a <= b {
			a += k
		} else {
			a -= k
		}
		if c <= d {
			c += k
		} else {
			c -= k
		}
		a, c = r.tree.IdToNode[a], r.tree.IdToNode[c]
		if r.data[a] < r.data[c] {
			return lcp, -1
		}
		if r.data[a] == r.data[c] {
			return lcp, 0
		}
		if r.data[a] > r.data[c] {
			return lcp, 1
		}
	}
	if len(path1) > 0 {
		return lcp, 1
	}
	if len(path2) > 0 {
		return lcp, -1
	}
	return lcp, 0
}

func (r *RollingHashOnTree) getUd(a, b int32) uint64 {
	if a == -1 {
		return r.dp1[b]
	}
	return modSub(r.dp1[b], modMul(r.dp1[a], r.pow[r.tree.Depth[b]-r.tree.Depth[a]]))
}

func (r *RollingHashOnTree) getDu(a, b int32) uint64 {
	if b == -1 {
		return r.dp2[a]
	}
	return modMul(modSub(r.dp2[a], r.dp2[b]), r.ipow[r.tree.Depth[b]+1])
}

func (r *RollingHashOnTree) fromHldPair(a, b int32) uint64 {
	if a <= b {
		return r.getUd(r.tree.Parent[r.tree.IdToNode[a]], r.tree.IdToNode[b])
	}
	return r.getDu(r.tree.IdToNode[a], r.tree.Parent[r.tree.IdToNode[b]])
}

type neighbor = struct {
	to   int32
	cost int
	eid  int32
}

type Tree32 struct {
	Lid, Rid      []int32
	IdToNode      []int32
	Depth         []int32
	DepthWeighted []int
	Parent        []int32
	Head          []int32 // 重链头
	Tree          [][]neighbor
	Edges         [][2]int32
	vToE          []int32 // 节点v的父边的id
	n             int32
}

func NewTree32(n int32) *Tree32 {
	res := &Tree32{Tree: make([][]neighbor, n), Edges: make([][2]int32, 0, n-1), n: n}
	return res
}

func (t *Tree32) AddEdge(u, v int32, w int) {
	eid := int32(len(t.Edges))
	t.Tree[u] = append(t.Tree[u], neighbor{to: v, cost: w, eid: eid})
	t.Tree[v] = append(t.Tree[v], neighbor{to: u, cost: w, eid: eid})
	t.Edges = append(t.Edges, [2]int32{u, v})
}

func (t *Tree32) AddDirectedEdge(from, to int32, cost int) {
	eid := int32(len(t.Edges))
	t.Tree[from] = append(t.Tree[from], neighbor{to: to, cost: cost, eid: eid})
	t.Edges = append(t.Edges, [2]int32{from, to})
}

func (t *Tree32) Build(root int32) {
	if root != -1 && int32(len(t.Edges)) != t.n-1 {
		panic("edges count != n-1")
	}
	n := t.n
	t.Lid = make([]int32, n)
	t.Rid = make([]int32, n)
	t.IdToNode = make([]int32, n)
	t.Depth = make([]int32, n)
	t.DepthWeighted = make([]int, n)
	t.Parent = make([]int32, n)
	t.Head = make([]int32, n)
	t.vToE = make([]int32, n)
	for i := int32(0); i < n; i++ {
		t.Depth[i] = -1
		t.Head[i] = root
		t.vToE[i] = -1
	}
	if root != -1 {
		t._dfsSize(root, -1)
		time := int32(0)
		t._dfsHld(root, &time)
	} else {
		time := int32(0)
		for i := int32(0); i < n; i++ {
			if t.Depth[i] == -1 {
				t._dfsSize(i, -1)
				t._dfsHld(i, &time)
			}
		}
	}
}

// 从v开始沿着重链向下收集节点.
func (t *Tree32) HeavyPathAt(v int32) []int32 {
	path := []int32{v}
	for {
		a := path[len(path)-1]
		for _, e := range t.Tree[a] {
			if e.to != t.Parent[a] && t.Head[e.to] == v {
				path = append(path, e.to)
				break
			}
		}
		if path[len(path)-1] == a {
			break
		}
	}
	return path
}

// 返回重儿子，如果没有返回 -1.
func (t *Tree32) HeavyChild(v int32) int32 {
	k := t.Lid[v] + 1
	if k == t.n {
		return -1
	}
	w := t.IdToNode[k]
	if t.Parent[w] == v {
		return w
	}
	return -1
}

// 从v开始向上走k步.
func (t *Tree32) KthAncestor(v, k int32) int32 {
	if k > t.Depth[v] {
		return -1
	}
	for {
		u := t.Head[v]
		if t.Lid[v]-k >= t.Lid[u] {
			return t.IdToNode[t.Lid[v]-k]
		}
		k -= t.Lid[v] - t.Lid[u] + 1
		v = t.Parent[u]
	}
}

func (t *Tree32) Lca(u, v int32) int32 {
	for {
		if t.Lid[u] > t.Lid[v] {
			u, v = v, u
		}
		if t.Head[u] == t.Head[v] {
			return u
		}
		v = t.Parent[t.Head[v]]
	}
}

func (t *Tree32) LcaRooted(u, v, root int32) int32 {
	return t.Lca(u, v) ^ t.Lca(u, root) ^ t.Lca(v, root)
}

func (t *Tree32) Dist(a, b int32) int32 {
	c := t.Lca(a, b)
	return t.Depth[a] + t.Depth[b] - 2*t.Depth[c]
}

func (t *Tree32) DistWeighted(a, b int32) int {
	c := t.Lca(a, b)
	return t.DepthWeighted[a] + t.DepthWeighted[b] - 2*t.DepthWeighted[c]
}

// c 是否在 p 的子树中.c和p不能相等.
func (t *Tree32) InSubtree(c, p int32) bool {
	return t.Lid[p] <= t.Lid[c] && t.Lid[c] < t.Rid[p]
}

// 从 a 开始走 k 步到 b.
func (t *Tree32) Jump(a, b, k int32) int32 {
	if k == 1 {
		if a == b {
			return -1
		}
		if t.InSubtree(b, a) {
			return t.KthAncestor(b, t.Depth[b]-t.Depth[a]-1)
		}
		return t.Parent[a]
	}
	c := t.Lca(a, b)
	dac := t.Depth[a] - t.Depth[c]
	dbc := t.Depth[b] - t.Depth[c]
	if k > dac+dbc {
		return -1
	}
	if k <= dac {
		return t.KthAncestor(a, k)
	}
	return t.KthAncestor(b, dac+dbc-k)
}

func (t *Tree32) SubtreeSize(v int32) int32 {
	return t.Rid[v] - t.Lid[v]
}

func (t *Tree32) SubtreeSizeRooted(v, root int32) int32 {
	if v == root {
		return t.n
	}
	x := t.Jump(v, root, 1)
	if t.InSubtree(v, x) {
		return t.Rid[v] - t.Lid[v]
	}
	return t.n - t.Rid[x] + t.Lid[x]
}

func (t *Tree32) CollectChild(v int32) []int32 {
	var res []int32
	for _, e := range t.Tree[v] {
		if e.to != t.Parent[v] {
			res = append(res, e.to)
		}
	}
	return res
}

// 收集与 v 相邻的轻边.
func (t *Tree32) CollectLight(v int32) []int32 {
	var res []int32
	skip := true
	for _, e := range t.Tree[v] {
		if e.to != t.Parent[v] {
			if !skip {
				res = append(res, e.to)
			}
			skip = false
		}
	}
	return res
}

func (tree *Tree32) RestorePath(from, to int32) []int32 {
	res := []int32{}
	composition := tree.GetPathDecomposition(from, to, true)
	for _, e := range composition {
		a, b := e[0], e[1]
		if a <= b {
			for i := a; i <= b; i++ {
				res = append(res, tree.IdToNode[i])
			}
		} else {
			for i := a; i >= b; i-- {
				res = append(res, tree.IdToNode[i])
			}
		}
	}
	return res
}

// 返回沿着`路径顺序`的 [起点,终点] 的 欧拉序 `左闭右闭` 数组.
//
//	!eg:[[2 0] [4 4]] 沿着路径顺序但不一定沿着欧拉序.
func (tree *Tree32) GetPathDecomposition(u, v int32, vertex bool) [][2]int32 {
	up, down := [][2]int32{}, [][2]int32{}
	for {
		if tree.Head[u] == tree.Head[v] {
			break
		}
		if tree.Lid[u] < tree.Lid[v] {
			down = append(down, [2]int32{tree.Lid[tree.Head[v]], tree.Lid[v]})
			v = tree.Parent[tree.Head[v]]
		} else {
			up = append(up, [2]int32{tree.Lid[u], tree.Lid[tree.Head[u]]})
			u = tree.Parent[tree.Head[u]]
		}
	}
	edgeInt := int32(1)
	if vertex {
		edgeInt = 0
	}
	if tree.Lid[u] < tree.Lid[v] {
		down = append(down, [2]int32{tree.Lid[u] + edgeInt, tree.Lid[v]})
	} else if tree.Lid[v]+edgeInt <= tree.Lid[u] {
		up = append(up, [2]int32{tree.Lid[u], tree.Lid[v] + edgeInt})
	}
	for i := 0; i < len(down)/2; i++ {
		down[i], down[len(down)-1-i] = down[len(down)-1-i], down[i]
	}
	return append(up, down...)
}

// 遍历路径上的 `[起点,终点)` 欧拉序 `左闭右开` 区间.
func (tree *Tree32) EnumeratePathDecomposition(u, v int32, vertex bool, f func(start, end int32)) {
	for {
		if tree.Head[u] == tree.Head[v] {
			break
		}
		if tree.Lid[u] < tree.Lid[v] {
			a, b := tree.Lid[tree.Head[v]], tree.Lid[v]
			if a > b {
				a, b = b, a
			}
			f(a, b+1)
			v = tree.Parent[tree.Head[v]]
		} else {
			a, b := tree.Lid[u], tree.Lid[tree.Head[u]]
			if a > b {
				a, b = b, a
			}
			f(a, b+1)
			u = tree.Parent[tree.Head[u]]
		}
	}

	edgeInt := int32(1)
	if vertex {
		edgeInt = 0
	}

	if tree.Lid[u] < tree.Lid[v] {
		a, b := tree.Lid[u]+edgeInt, tree.Lid[v]
		if a > b {
			a, b = b, a
		}
		f(a, b+1)
	} else if tree.Lid[v]+edgeInt <= tree.Lid[u] {
		a, b := tree.Lid[u], tree.Lid[v]+edgeInt
		if a > b {
			a, b = b, a
		}
		f(a, b+1)
	}
}

// 返回 root 的欧拉序区间, 左闭右开, 0-indexed.
func (tree *Tree32) Id(root int32) (int32, int32) {
	return tree.Lid[root], tree.Rid[root]
}

// 点v对应的父边的边id.如果v是根节点则返回-1.
func (tre *Tree32) VToE(v int32) int32 {
	return tre.vToE[v]
}

// 第i条边对应的深度更深的那个节点.
func (tree *Tree32) EToV(i int32) int32 {
	u, v := tree.Edges[i][0], tree.Edges[i][1]
	if tree.Parent[u] == v {
		return u
	}
	return v
}

func (t *Tree32) _dfsSize(cur, pre int32) {
	size := t.Rid
	t.Parent[cur] = pre
	if pre != -1 {
		t.Depth[cur] = t.Depth[pre] + 1
	} else {
		t.Depth[cur] = 0
	}
	size[cur] = 1
	nexts := t.Tree[cur]
	for i := int32(len(nexts)) - 2; i >= 0; i-- {
		e := nexts[i+1]
		if t.Depth[e.to] == -1 {
			nexts[i], nexts[i+1] = nexts[i+1], nexts[i]
		}
	}
	hldSize := int32(0)
	for i, e := range nexts {
		to := e.to
		if t.Depth[to] == -1 {
			t.DepthWeighted[to] = t.DepthWeighted[cur] + e.cost
			t.vToE[to] = e.eid
			t._dfsSize(to, cur)
			size[cur] += size[to]
			if size[to] > hldSize {
				hldSize = size[to]
				if i != 0 {
					nexts[0], nexts[i] = nexts[i], nexts[0]
				}
			}
		}
	}
}

func (t *Tree32) _dfsHld(cur int32, times *int32) {
	t.Lid[cur] = *times
	*times++
	t.Rid[cur] += t.Lid[cur]
	t.IdToNode[t.Lid[cur]] = cur
	heavy := true
	for _, e := range t.Tree[cur] {
		to := e.to
		if t.Depth[to] > t.Depth[cur] {
			if heavy {
				t.Head[to] = t.Head[cur]
			} else {
				t.Head[to] = to
			}
			heavy = false
			t._dfsHld(to, times)
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
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

func min32(a, b int32) int32 {
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

func binarySearch32(check func(x int32) bool, ok, ng int32) int32 {
	for abs32(ok-ng) > 1 {
		x := (ng + ok) / 2
		if check(x) {
			ok = x
		} else {
			ng = x
		}
	}
	return ok
}

// ModFast/fastMod/mod61
const (
	hashStringMod    uint64 = (1 << 61) - 1
	hashStringModi64 int64  = (1 << 61) - 1
	hashStringMask30 uint64 = (1 << 30) - 1
	hashStringMask31 uint64 = (1 << 31) - 1
	hashStringMASK61 uint64 = hashStringMod
)

// a*b % (2^61-1)
func modMul(a, b uint64) uint64 {
	au := a >> 31
	ad := a & hashStringMask31
	bu := b >> 31
	bd := b & hashStringMask31
	mid := ad*bu + au*bd
	midu := mid >> 30
	midd := mid & hashStringMask30
	return mod(au*bu<<1 + midu + (midd << 31) + ad*bd)
}

// x % (2^61-1)
func mod(x uint64) uint64 {
	xu := x >> 61
	xd := x & hashStringMASK61
	res := xu + xd
	if res >= hashStringMod {
		res -= hashStringMod
	}
	return res
}

func modInv(x uint64) uint64 {
	a, b, u, v, t := int64(x), hashStringModi64, int64(1), int64(0), int64(0)
	for b > 0 {
		t = a / b
		a -= t * b
		a, b = b, a
		u -= t * v
		u, v = v, u
	}
	u %= hashStringModi64
	if u < 0 {
		u += hashStringModi64
	}
	return uint64(u)
}

func modAdd(a, b uint64) uint64 {
	res := a + b
	if res >= hashStringMod {
		res -= hashStringMod
	}
	return res
}

func modSub(a, b uint64) uint64 {
	tmp := a - b
	if tmp >= hashStringMod {
		return modAdd(tmp, hashStringMod)
	}
	return tmp
}
