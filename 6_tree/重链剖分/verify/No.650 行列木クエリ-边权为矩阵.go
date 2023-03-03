// https://yukicoder.me/problems/no/650

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	hld := NewHeavyLightDecomposition(n)
	edges := make([][2]int, n-1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		hld.AddEdge(u, v)
		edges[i] = [2]int{u, v}
	}
	hld.Build(0)

	leaves := make([]E, n)
	for i := range leaves {
		leaves[i] = E{{1, 0}, {0, 1}}
	}
	seg := NewSegmentTree(leaves)

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var op string
		fmt.Fscan(in, &op)
		if op == "x" {
			// 编号为ei的边上的矩阵变为x00 x01 x10 x11
			var ei, x00, x01, x10, x11 int
			fmt.Fscan(in, &ei, &x00, &x01, &x10, &x11)
			u, v := edges[ei][0], edges[ei][1]
			// 边的序号用较深的那个顶点的欧拉序的起点编号表示.
			in1, _ := hld.Id(u)
			in2, _ := hld.Id(v)
			if in1 < in2 {
				in1, in2 = in2, in1
			}
			seg.Set(in1, E{{x00, x01}, {x10, x11}})
		} else if op == "g" {
			// 查询从anscestor到cur的路径上的矩阵的乘积(从根到叶子方向)
			var ancestor, cur int
			fmt.Fscan(in, &ancestor, &cur)
			res := seg.e()
			hld.QueryNonCommutativePath(ancestor, cur, false, func(start, end int) {
				res = seg.op(res, seg.Query(start, end))
			})
			fmt.Fprintln(out, res[0][0], res[0][1], res[1][0], res[1][1])
		}
	}
}

const MOD int = 1e9 + 7

type E = [2][2]int

func (*SegmentTree) e() E { return E{{1, 0}, {0, 1}} }
func (*SegmentTree) op(a, b E) E {
	return E{
		{(a[0][0]*b[0][0] + a[0][1]*b[1][0]) % MOD, (a[0][0]*b[0][1] + a[0][1]*b[1][1]) % MOD},
		{(a[1][0]*b[0][0] + a[1][1]*b[1][0]) % MOD, (a[1][0]*b[0][1] + a[1][1]*b[1][1]) % MOD},
	}
}

type SegmentTree struct {
	n, size int
	seg     []E
}

func NewSegmentTree(leaves []E) *SegmentTree {
	res := &SegmentTree{}
	n := len(leaves)
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, 2*size)
	for i := 0; i < n; i++ {
		seg[i+size] = leaves[i]
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[2*i], seg[2*i+1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}

func (st *SegmentTree) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}

func (st *SegmentTree) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[2*index], st.seg[2*index+1])
	}
}

// [start, end)
func (st *SegmentTree) Query(start, end int) E {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return st.e()
	}
	leftRes, rightRes := st.e(), st.e()
	start += st.size
	end += st.size
	for start < end {
		if start&1 == 1 {
			leftRes = st.op(leftRes, st.seg[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = st.op(st.seg[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return st.op(leftRes, rightRes)
}

func (st *SegmentTree) QueryAll() E { return st.seg[1] }

// maxRight returns the maximum r such that [start, r) satisfies the predicate.
func (st *SegmentTree) MaxRight(start int, predicate func(E) bool) int {
	if start == st.n {
		return st.n
	}

	start += st.size
	res := st.e()
	for {
		for start&1 == 0 {
			start >>= 1
		}
		if !predicate(st.op(res, st.seg[start])) {
			for start < st.size {
				start = 2 * start
				if predicate(st.op(res, st.seg[start])) {
					res = st.op(res, st.seg[start])
					start++
				}
			}

			return start - st.size
		}
		res = st.op(res, st.seg[start])
		start++
		if (start & -start) == start {
			break
		}
	}
	return st.n
}

// minLeft returns the minimum l such that [l, end) satisfies the predicate.
func (st *SegmentTree) MinLeft(end int, predicate func(E) bool) int {
	if end == 0 {
		return 0
	}
	end += st.size
	sm := st.e()
	for {
		end--
		for end > 1 && end&1 == 1 {
			end >>= 1
		}
		if !predicate(st.op(st.seg[end], sm)) {
			for end < st.size {
				end = 2*end + 1
				if predicate(st.op(st.seg[end], sm)) {
					sm = st.op(st.seg[end], sm)
					end--
				}
			}
			return end + 1 - st.size
		}
		sm = st.op(st.seg[end], sm)
		if end&-end == end {
			break
		}
	}
	return 0
}

type HeavyLightDecomposition struct {
	Parent   []int
	Depth    []int
	Size     []int
	g        [][]int
	id       int
	down, up []int
	nxt      []int // heavy pathの先頭
}

func NewHeavyLightDecomposition(n int) *HeavyLightDecomposition {
	return &HeavyLightDecomposition{g: make([][]int, n)}
}

// 無向辺 u <-> v を追加する.
func (hld *HeavyLightDecomposition) AddEdge(u, v int) {
	hld.g[u] = append(hld.g[u], v)
	hld.g[v] = append(hld.g[v], u)
}

// 有向辺 u -> v を追加する.
func (hld *HeavyLightDecomposition) AddDirectedEdge(u, v int) {
	hld.g[u] = append(hld.g[u], v)
}

// rootを根とした重軽分解を構築する.
func (hld *HeavyLightDecomposition) Build(root int) {
	n := len(hld.g)
	hld.Size = make([]int, n)
	hld.Depth = make([]int, n)
	hld.down = make([]int, n)
	hld.up = make([]int, n)
	hld.nxt = make([]int, n)
	hld.Parent = make([]int, n)
	for i := 0; i < n; i++ {
		hld.down[i] = -1
		hld.up[i] = -1
		hld.nxt[i] = root
		hld.Parent[i] = root
	}

	hld.dfsSize(root, -1)
	hld.dfsHld(root, -1)
}

// 頂点 i のオイラーツアー順を [down,up) の形で返す.
//  0 <= down < up <= n.
func (hld *HeavyLightDecomposition) Id(u int) (down, up int) {
	down, up = hld.down[u], hld.up[u]
	return
}

// 可換なパスクエリを処理する.
//   0 <= start <= end <= n, [start,end).
func (hld *HeavyLightDecomposition) QueryPath(u, v int, vertex bool, f func(start, end int)) {
	lca_ := hld.LCA(u, v)
	for _, p := range hld.ascend(u, lca_) {
		s, t := p[0]+1, p[1]
		if s > t {
			f(t, s)
		} else {
			f(s, t)
		}
	}
	if vertex {
		f(hld.down[lca_], hld.down[lca_]+1)
	}
	for _, p := range hld.descend(lca_, v) {
		s, t := p[0], p[1]+1
		if s > t {
			f(t, s)
		} else {
			f(s, t)
		}
	}
}

// 非可換なパスクエリを処理する.
//   0 <= start <= end <= n, [start,end).
//   https://nyaannyaan.github.io/library/verify/verify-yosupo-ds/yosupo-vertex-set-path-composite.test.cpp
func (hld *HeavyLightDecomposition) QueryNonCommutativePath(u, v int, vertex bool, f func(start, end int)) {
	lca_ := hld.LCA(u, v)
	for _, p := range hld.ascend(u, lca_) {
		f(p[0]+1, p[1])
	}
	if vertex {
		f(hld.down[lca_], hld.down[lca_]+1)
	}
	for _, p := range hld.descend(lca_, v) {
		f(p[0], p[1]+1)
	}
}

// 部分木クエリを処理する.
//   0 <= start <= end <= n, [start,end).
func (hld *HeavyLightDecomposition) QuerySubTree(u int, vertex bool, f func(start, end int)) {
	if vertex {
		f(hld.down[u], hld.up[u])
	} else {
		f(hld.down[u]+1, hld.up[u])
	}
}

func (hld *HeavyLightDecomposition) LCA(u, v int) int {
	for hld.nxt[u] != hld.nxt[v] {
		if hld.down[u] < hld.down[v] {
			u, v = v, u
		}
		u = hld.Parent[hld.nxt[u]]
	}
	if hld.Depth[u] < hld.Depth[v] {
		return u
	}
	return v
}

func (hld *HeavyLightDecomposition) Dist(u, v int) int {
	return hld.Depth[u] + hld.Depth[v] - hld.Depth[hld.LCA(u, v)]*2
}

func (hld *HeavyLightDecomposition) dfsSize(cur, pre int) {
	hld.Size[cur] = 1
	for i, to := range hld.g[cur] {
		if to == pre {
			continue
		}
		// if to == hld.Parent[cur] {
		// 	if len(hld.g[cur]) >= 2 && hld.g[cur][0] == to {
		// 		hld.g[cur][0], hld.g[cur][1] = hld.g[cur][1], hld.g[cur][0]
		// 	} else {
		// 		continue
		// 	}
		// }

		hld.Depth[to] = hld.Depth[cur] + 1
		hld.Parent[to] = cur
		hld.dfsSize(to, cur)
		hld.Size[cur] += hld.Size[to]
		if hld.Size[to] > hld.Size[hld.g[cur][0]] {
			hld.g[cur][0], hld.g[cur][i] = hld.g[cur][i], hld.g[cur][0]
		}
	}

}

func (hld *HeavyLightDecomposition) dfsHld(cur, pre int) {
	hld.down[cur] = hld.id
	hld.id++
	for _, to := range hld.g[cur] {
		if to == pre {
			continue
		}
		if to == hld.g[cur][0] {
			hld.nxt[to] = hld.nxt[cur]
		} else {
			hld.nxt[to] = to
		}
		hld.dfsHld(to, cur)
	}
	hld.up[cur] = hld.id
}

// [u, v)
func (hld *HeavyLightDecomposition) ascend(u, v int) [][2]int {
	var res [][2]int
	for hld.nxt[u] != hld.nxt[v] {
		res = append(res, [2]int{hld.down[u], hld.down[hld.nxt[u]]})
		u = hld.Parent[hld.nxt[u]]
	}
	if u != v {
		res = append(res, [2]int{hld.down[u], hld.down[v] + 1})
	}
	return res
}

// (u, v]
func (hld *HeavyLightDecomposition) descend(u, v int) [][2]int {
	if u == v {
		return nil
	}
	if hld.nxt[u] == hld.nxt[v] {
		return [][2]int{{hld.down[u] + 1, hld.down[v]}}
	}
	res := hld.descend(u, hld.Parent[hld.nxt[v]])
	res = append(res, [2]int{hld.down[hld.nxt[v]], hld.down[v]})
	return res
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

// !Range Add Range Sum, 0-based.
type BITArray struct {
	n     int
	tree1 []int
	tree2 []int
}

func NewBITArray(n int) *BITArray {
	return &BITArray{
		n:     n,
		tree1: make([]int, n+1),
		tree2: make([]int, n+1),
	}
}

// 切片内[start, end)的每个元素加上delta.
//  0<=start<=end<=n
func (b *BITArray) Add(start, end, delta int) {
	end--
	b.add(start, delta)
	b.add(end+1, -delta)
}

// 求切片内[start, end)的和.
//  0<=start<=end<=n
func (b *BITArray) Query(start, end int) int {
	end--
	return b.query(end) - b.query(start-1)
}

func (b *BITArray) add(index, delta int) {
	index++
	rawIndex := index
	for index <= b.n {
		b.tree1[index] += delta
		b.tree2[index] += (rawIndex - 1) * delta
		index += index & -index
	}
}

func (b *BITArray) query(index int) (res int) {
	index++
	if index > b.n {
		index = b.n
	}
	rawIndex := index
	for index > 0 {
		res += rawIndex*b.tree1[index] - b.tree2[index]
		index -= index & -index
	}
	return
}
