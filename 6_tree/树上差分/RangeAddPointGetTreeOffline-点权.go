// RangeAddPointGetTreeOffline-点权
// !这里的差分为点权的差分

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// https://www.luogu.com.cn/problem/P3128
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int32
	fmt.Fscan(in, &n, &k)

	tree := make([][]int32, n)
	for i := int32(0); i < n-1; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u--
		v--
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	R := NewRangeAddPointGetTreeOffline(tree, 0)
	for i := int32(0); i < k; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u--
		v--
		R.AddPath(u, v, 1)
	}

	res := 0
	for i := int32(0); i < n; i++ {
		res = max(res, int(R.GetPoint(i)))
	}
	fmt.Fprintln(out, res)
}

type E = int

func e() E          { return 0 }
func op(e1, e2 E) E { return e1 + e2 }
func inv(e E) E     { return -e }

type RangeAddPointGetTreeOffline struct {
	tree   [][]int32
	root   int32
	lca    *LCAHLD
	sum    []E
	preSum []E
	dirty  bool
}

// 树上差分离线版.区间加,单点查询.
func NewRangeAddPointGetTreeOffline(tree [][]int32, root int32) *RangeAddPointGetTreeOffline {
	n := len(tree)
	sum := make([]E, n)
	for i := 0; i < n; i++ {
		sum[i] = e()
	}
	lca := NewLCA(tree, root)
	return &RangeAddPointGetTreeOffline{
		tree:  tree,
		root:  root,
		lca:   lca,
		sum:   sum,
		dirty: true,
	}
}

func NewRangeAddPointGetTreeOfflineWithValues(tree [][]int32, root int32, values []E) *RangeAddPointGetTreeOffline {
	n := len(tree)
	sum := make([]E, n)
	copy(sum, values)
	lca := NewLCA(tree, root)
	return &RangeAddPointGetTreeOffline{
		tree: tree,
		root: root,
		lca:  lca,
		sum:  sum,
	}
}

func (r *RangeAddPointGetTreeOffline) AddPoint(node int32, weight E) {
	if node == r.root {
		r.sum[0] = op(r.sum[0], weight)
	} else {
		cur := r.lca.LId[node]
		r.sum[cur] = op(r.sum[cur], weight)
		parent := r.lca.LId[r.lca.Parent[node]]
		r.sum[parent] = op(r.sum[parent], inv(weight))
	}
	r.dirty = true
}

// 路径上所有点加上delta.
func (r *RangeAddPointGetTreeOffline) AddPath(u, v int32, delta E) {
	lid := r.lca.LId
	r.sum[lid[u]] = op(r.sum[lid[u]], delta)
	r.sum[lid[v]] = op(r.sum[lid[v]], delta)
	lca := r.lca.LCA(u, v)
	r.sum[lid[lca]] = op(r.sum[lid[lca]], inv(delta))
	if lca != r.root {
		parent := lid[r.lca.Parent[lca]]
		r.sum[parent] = op(r.sum[parent], inv(delta))
	}
	r.dirty = true
}

// 树链并.这里的树链为根节点到各个point的路径.
func (r *RangeAddPointGetTreeOffline) AddChains(chainEnds []int32, weight E) {
	if len(chainEnds) == 0 {
		return
	}
	dfns := make([]int32, len(chainEnds))
	lid := r.lca.LId
	idToNode := r.lca.IdToNode
	for i, end := range chainEnds {
		dfns[i] = lid[end]
	}
	sort.Slice(dfns, func(i, j int) bool { return dfns[i] < dfns[j] })
	r.sum[dfns[0]] = op(r.sum[dfns[0]], weight)
	for i := 1; i < len(dfns); i++ {
		r.sum[dfns[i]] = op(r.sum[dfns[i]], weight)
		u, v := idToNode[dfns[i-1]], idToNode[dfns[i]]
		lca := r.lca.LCA(u, v)
		r.sum[lid[lca]] = op(r.sum[lid[lca]], inv(weight))
	}
	r.dirty = true
}

// O(n) 构建.
func (r *RangeAddPointGetTreeOffline) GetPoint(node int32) E {
	if r.dirty {
		r.preSum = make([]E, len(r.sum)+1)
		r.preSum[0] = e()
		for i := 1; i <= len(r.sum); i++ {
			r.preSum[i] = op(r.preSum[i-1], r.sum[i-1])
		}
		r.dirty = false
	}
	start, end := r.lca.LId[node], r.lca.RId[node]
	return op(r.preSum[end], inv(r.preSum[start]))
}

type LCAHLD struct {
	Depth, Parent []int32
	LId, RId      []int32
	IdToNode      []int32
	tree          [][]int32
	top, heavySon []int32
	dfnId         int32
}

func NewLCA(tree [][]int32, root int32) *LCAHLD {
	n := int32(len(tree))
	lid := make([]int32, n)
	rid := make([]int32, n)
	idToNode := make([]int32, n)
	top := make([]int32, n)      // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	depth := make([]int32, n)    // 深度
	parent := make([]int32, n)   // 父结点
	heavySon := make([]int32, n) // 重儿子
	for i := range parent {
		parent[i] = -1
	}

	res := &LCAHLD{
		tree:     tree,
		Depth:    depth,
		Parent:   parent,
		LId:      lid,
		RId:      rid,
		IdToNode: idToNode,
		top:      top,
		heavySon: heavySon,
	}
	res.build(root, -1, 0)
	res.markTop(root, root)
	return res
}

func (hld *LCAHLD) LCA(u, v int32) int32 {
	for {
		if hld.LId[u] > hld.LId[v] {
			u, v = v, u
		}
		if hld.top[u] == hld.top[v] {
			return u
		}
		v = hld.Parent[hld.top[v]]
	}
}

func (hld *LCAHLD) Dist(u, v int32) int32 {
	return hld.Depth[u] + hld.Depth[v] - 2*hld.Depth[hld.LCA(u, v)]
}

func (hld *LCAHLD) build(cur, pre, dep int32) int32 {
	subSize, heavySize, heavySon := int32(1), int32(0), int32(-1)
	for _, next := range hld.tree[cur] {
		if next != pre {
			nextSize := hld.build(next, cur, dep+1)
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

func (hld *LCAHLD) markTop(cur, top int32) {
	hld.top[cur] = top
	hld.LId[cur] = hld.dfnId
	hld.IdToNode[hld.dfnId] = cur
	hld.dfnId++
	heavySon := hld.heavySon[cur]
	if heavySon != -1 {
		hld.markTop(heavySon, top)
		for _, next := range hld.tree[cur] {
			if next != heavySon && next != hld.Parent[cur] {
				hld.markTop(next, next)
			}
		}
	}
	hld.RId[cur] = hld.dfnId
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
