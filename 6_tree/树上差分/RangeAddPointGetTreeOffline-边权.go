// RangeAddPointGetTreeOffline-边权
// !这里的差分为边权的差分

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// demo()
	CF191C()
}

// Fools and Roads
// https://www.luogu.com.cn/problem/CF191C
// 有一颗 n 个节点的树，k 次旅行，问每一条边被走过的次数。
func CF191C() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	edgeId := make([]map[int32]int32, n)
	for i := int32(0); i < n; i++ {
		edgeId[i] = map[int32]int32{}
	}
	tree := make([][]int32, n)
	for i := int32(0); i < n-1; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u--
		v--
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
		if u > v {
			u, v = v, u
		}
		edgeId[u][v] = i
	}

	R := NewRangeAddPointGetTreeOfflineEdge(tree, 0)

	var k int32
	fmt.Fscan(in, &k)
	for i := int32(0); i < k; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		R.AddPath(u, v, 1)
	}

	weightToParent := make([]int32, n)
	for i := int32(0); i < n; i++ {
		weightToParent[i] = R.GetPoint(i)
	}
	res := make([]int32, n-1)
	for i, w := range weightToParent {
		parent := R.lca.Parent[i]
		cur := int32(i)
		if parent != -1 {
			if cur > parent {
				cur, parent = parent, cur
			}
			id := edgeId[cur][parent]
			res[id] = w
		}
	}

	for i := int32(0); i < n-1; i++ {
		fmt.Fprint(out, res[i], " ")
	}
}

func demo() {
	//    0
	//  / | \
	// 1  2  3
	//    |
	//    4
	tree := [][]int32{{1, 2, 3}, {0}, {0, 4}, {0}, {2}}
	R := NewRangeAddPointGetTreeOfflineEdge(tree, 0)
	R.AddPath(0, 1, 1)
	R.AddPath(0, 2, 2)
	R.AddPath(0, 3, 76)
	R.AddPath(2, 4, 9)
	fmt.Println(R.GetPoint(0)) // 0
	fmt.Println(R.GetPoint(1)) // 1
	fmt.Println(R.GetPoint(2)) // 2
	fmt.Println(R.GetPoint(3)) // 76
	fmt.Println(R.GetPoint(4)) // 9
}

type E = int32

func e() E          { return 0 }
func op(e1, e2 E) E { return e1 + e2 }
func inv(e E) E     { return -e }

type RangeAddPointGetTreeOfflineEdge struct {
	tree   [][]int32
	root   int32
	lca    *LCAHLD
	sum    []E
	preSum []E
	dirty  bool
}

// 树上差分离线版.区间加,单点查询.
func NewRangeAddPointGetTreeOfflineEdge(tree [][]int32, root int32) *RangeAddPointGetTreeOfflineEdge {
	n := len(tree)
	data := make([]E, n)
	for i := 0; i < n; i++ {
		data[i] = e()
	}
	lca := NewLCA(tree, root)
	return &RangeAddPointGetTreeOfflineEdge{
		tree:  tree,
		root:  root,
		lca:   lca,
		sum:   data,
		dirty: true,
	}
}

// 路径上所有点加上delta.
// sum[u]++, sum[v]++, sum[lca]-=2
func (r *RangeAddPointGetTreeOfflineEdge) AddPath(u, v int32, delta E) {
	lid := r.lca.LId
	r.sum[lid[u]] = op(r.sum[lid[u]], delta)
	r.sum[lid[v]] = op(r.sum[lid[v]], delta)
	lca := r.lca.LCA(u, v)
	r.sum[lid[lca]] = op(r.sum[lid[lca]], op(inv(delta), inv(delta)))
	r.dirty = true
}

// O(n) 构建.查询node与父亲结点的边权.
func (r *RangeAddPointGetTreeOfflineEdge) GetPoint(node int32) E {
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
