package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	tree := [][]int{{1, 2}, {0, 3}, {0, 4}, {1}, {2}}
	R := NewRangeAddPointGetTreeOnline(tree, 0, true)
	R.AddPoint(0, 1)
	R.AddPoint(2, 1)
	fmt.Println(R.GetPoint(0))
	fmt.Println(R.GetPoint(1))
	fmt.Println(R.GetPoint(2))
	fmt.Println(R.GetPoint(3))
}

// https://www.luogu.com.cn/problem/P3128
func luogu3128() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)

	tree := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	R := NewRangeAddPointGetTreeOnline(tree, 0, true)
	for i := 0; i < k; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		R.AddPath(u, v, 1)
	}

	res := 0
	for i := 0; i < n; i++ {
		res = max(res, R.GetPoint(i))
	}
	fmt.Fprintln(out, res)
}

type E = int

func e() E          { return 0 }
func op(e1, e2 E) E { return e1 + e2 }
func inv(e E) E     { return -e }

type RangeAddPointGetTreeOnline struct {
	tree [][]int
	root int
	lca  *LCAHLD
	bit  *BIT
}

// 树上差分在线版.区间加,单点查询.
func NewRangeAddPointGetTreeOnline(tree [][]int, root int, isVertex bool) *RangeAddPointGetTreeOnline {
	n := len(tree)
	lca := NewLCA(tree, root)
	bit := NewBIT(n, func(i int) E { return e() })
	return &RangeAddPointGetTreeOnline{
		tree: tree,
		root: root,
		lca:  lca,
		bit:  bit,
	}
}

func (r *RangeAddPointGetTreeOnline) AddPoint(node int, delta E) {
	if node == r.root {
		r.bit.Add(0, delta)
	} else {
		r.bit.Add(r.lca.LId[node], delta)
		r.bit.Add(r.lca.LId[r.lca.Parent[node]], inv(delta))
	}
}

// 路径上所有点加上delta.
func (r *RangeAddPointGetTreeOnline) AddPath(u, v int, delta E) {
	r.bit.Add(r.lca.LId[u], delta)
	r.bit.Add(r.lca.LId[v], delta)
	lca := r.lca.LCA(u, v)
	r.bit.Add(r.lca.LId[lca], inv(delta))
	if lca != r.root {
		r.bit.Add(r.lca.LId[r.lca.Parent[lca]], inv(delta))
	}
}

// 树链并.这里的树链为根节点到各个point的路径.
func (r *RangeAddPointGetTreeOnline) AddChains(chainEnds []int, weight E) {
	if len(chainEnds) == 0 {
		return
	}
	dfns := make([]int, len(chainEnds))
	lid := r.lca.LId
	idToNode := r.lca.IdToNode
	for i, end := range chainEnds {
		dfns[i] = lid[end]
	}
	sort.Ints(dfns)
	r.bit.Add(dfns[0], weight)
	for i := 1; i < len(dfns); i++ {
		r.bit.Add(dfns[i], weight)
		u, v := idToNode[dfns[i-1]], idToNode[dfns[i]]
		lca := r.lca.LCA(u, v)
		r.bit.Add(lid[lca], inv(weight))
	}
}

func (r *RangeAddPointGetTreeOnline) GetPoint(node int) E {
	start, end := r.lca.LId[node], r.lca.RId[node]
	return r.bit.QueryRange(start, end)
}

type LCAHLD struct {
	Depth, Parent []int
	LId, RId      []int
	IdToNode      []int
	tree          [][]int
	top, heavySon []int
	dfnId         int
}

func NewLCA(tree [][]int, root int) *LCAHLD {
	n := len(tree)
	lid := make([]int, n)
	rid := make([]int, n)
	idToNode := make([]int, n)
	top := make([]int, n)      // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	depth := make([]int, n)    // 深度
	parent := make([]int, n)   // 父结点
	heavySon := make([]int, n) // 重儿子
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

func (hld *LCAHLD) LCA(u, v int) int {
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

func (hld *LCAHLD) Dist(u, v int) int {
	return hld.Depth[u] + hld.Depth[v] - 2*hld.Depth[hld.LCA(u, v)]
}
func (hld *LCAHLD) build(cur, pre, dep int) int {
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
	hld.Depth[cur] = dep
	hld.heavySon[cur] = heavySon
	hld.Parent[cur] = pre
	return subSize
}

func (hld *LCAHLD) markTop(cur, top int) {
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

type BIT struct {
	n    int
	data []E
}

func NewBIT(n int, f func(i int) E) *BIT {
	data := make([]E, n)
	for i := 0; i < n; i++ {
		data[i] = f(i)
	}
	for i := 1; i <= n; i++ {
		j := i + (i & -i)
		if j <= n {
			data[j-1] = op(data[j-1], data[i-1])
		}
	}
	return &BIT{n: n, data: data}
}

func (b *BIT) Add(index int, v E) {
	for index++; index <= b.n; index += index & -index {
		b.data[index-1] = op(b.data[index-1], v)
	}
}

// [0, end).
func (b *BIT) QueryPrefix(end int) E {
	if end > b.n {
		end = b.n
	}
	res := e()
	for ; end > 0; end -= end & -end {
		res = op(res, b.data[end-1])
	}
	return res
}

// [start, end).
func (b *BIT) QueryRange(start, end int) E {
	if start < 0 {
		start = 0
	}
	if end > b.n {
		end = b.n
	}
	if start >= end {
		return e()
	}
	if start == 0 {
		return b.QueryPrefix(end)
	}
	pos, neg := e(), e()
	for end > start {
		pos = op(pos, b.data[end-1])
		end &= end - 1
	}
	for start > end {
		neg = op(neg, b.data[start-1])
		start &= start - 1
	}
	return op(pos, inv(neg))
}

func (b *BIT) String() string {
	sb := []string{}
	for i := 0; i < b.n; i++ {
		sb = append(sb, fmt.Sprintf("%d", b.QueryRange(i, i+1)))
	}
	return fmt.Sprintf("BIT: [%v]", strings.Join(sb, ", "))
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
