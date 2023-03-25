// !在原图的点双形成的树中将割点提出来，并连接相邻的VBCC，得到的树就是BlockCutTree。
// VBCC和割点交错连接.
// 用于处理`必经点`的问题/处理删点的连通性问题
// 从x到y的必经顶点 => BlockCutTree 上对应的两点路径上的割点 => LCA求距离

// https://oi-wiki.org/graph/block-forest/
// 例题:
// 有一张 n 个点 m 条边的无向连通图，还有 q 个点对，
// 你需要输出每个点是多少给定点对的`必经点`（即如果点对为 (u,v)，
// 那么如果 u 到 v 无论如何都要经过 x ，那么 x 是该点对的必经点）
// !直接建出BlockCutTree，发现(u,v) 在BlockCutTree路径上的圆点都是必经点，lca 树上差分一下就可以了。

// !注意特殊处理原图中某个连通分量只有一个点(孤立点).
// !在这里,不将孤立点当作点双,即孤立点不会出现在BlockCutTree中.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	yukicode1326()
}

func 从BlockCutTree还原点双连通分量() {
	// https://judge.yosupo.jp/problem/biconnected_components
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	graph := make([][]Edge, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		graph[u] = append(graph[u], Edge{u, v, 1, i})
		graph[v] = append(graph[v], Edge{v, u, 1, i})
	}
	B := NewBlockCutTree(graph)
	B.Build()

	tree := B.Tree
	vbccSize := len(tree) - B.CountCut()
	group := make([][]int, vbccSize, n)
	for i := 0; i < n; i++ {
		if B.IsRawIsolate(i) { // 孤立点特殊处理，作为一组
			group = append(group, []int{i})
			continue
		}
		if B.IsRawCut(i) {
			for _, v := range tree[B.Get(i)] {
				group[v] = append(group[v], i)
			}
		} else {
			group[B.Get(i)] = append(group[B.Get(i)], i)
		}
	}

	fmt.Fprintln(out, len(group))
	for i := 0; i < len(group); i++ {
		fmt.Fprint(out, len(group[i]))
		for _, v := range group[i] {
			fmt.Fprint(out, " ", v)
		}
		fmt.Fprintln(out)
	}
}

func 网络集群() {
	// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=3022
	// 给定一个`无向连通图`
	// 网络连接,对一个网络集群维修,每天维修一台机器,需要断开它和其他机器的连接
	// !对每个顶点,移除他连接的所有边后,求剩下的连通分量的权值和的最大值(総合性能値)
	// BlockCutTree上dp
	// 如果i不是割点的话,就是原图中的所有点的权值和减去这个点的权值
	// 如果i是割点的话,就是max(子树和,allSum-子树和-这个点的权值)

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	values := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &values[i])
	}

	rawGraph := make([][]Edge, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		rawGraph[u] = append(rawGraph[u], Edge{u, v, 1, i})
		rawGraph[v] = append(rawGraph[v], Edge{v, u, 1, i})
	}

	B := NewBlockCutTree(rawGraph)
	B.Build()

	tree := B.Tree
	weight := make([]int, len(tree)) // 每个BlockCutTree的顶点的权值和
	all := 0
	for i := 0; i < n; i++ {
		weight[B.Get(i)] += values[i]
		all += values[i]
	}

	res := make([]int, n)
	for i := 0; i < n; i++ {
		if !B.IsRawCut(i) {
			res[i] = all - values[i] // 如果i不是割点的话,就是原图中的所有点的权值和减去这个点的权值
		}
	}

	var dfs func(cur, parent int) int
	dfs = func(cur, parent int) int {
		max_, sum := 0, 0
		for _, next := range tree[cur] {
			if next == parent {
				continue
			}
			nextRes := dfs(next, cur)
			sum += nextRes
			max_ = max(max_, nextRes)
		}
		if B.IsNewCut(cur) {
			res[B.Group[cur][0]] = max(max_, all-sum-weight[cur])
		}
		return sum + weight[cur]
	}
	dfs(0, -1)

	for i := 0; i < n; i++ {
		fmt.Fprintln(out, res[i])
	}
}

func yukicode1326() {
	// https://yukicoder.me/problems/no/1326
	// No.1326 ふたりのDominator
	// 给定一个无向连通图
	// 对每组顶点(x,y), 你需要选定一个不为x,y的顶点z,删除z的所有邻接边,
	//                   使得x,y不再连通,问有多少种方案
	// n,m<=5e4 q<=1e5
	// !这样的点是在从x到y的必经顶点 => BlockCutTree 上对应的两点路径上的割点 => LCA求距离

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	graph := make([][]Edge, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		graph[u] = append(graph[u], Edge{u, v, 1, i})
		graph[v] = append(graph[v], Edge{v, u, 1, i})
	}

	B := NewBlockCutTree(graph)
	B.Build()
	tree := B.Tree
	hld := NewHeavyLightDecomposition(len(tree))
	for i := 0; i < len(tree); i++ {
		for _, v := range tree[i] {
			hld.AddDirectedEdge(i, v)
		}
	}
	hld.Build(0)

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		x--
		y--
		if x == y {
			fmt.Fprintln(out, 0)
			continue
		}

		id1, id2 := B.Get(x), B.Get(y)
		res := hld.Dist(id1, id2)
		if B.IsNewCut(id1) {
			res--
		}
		if B.IsNewCut(id2) {
			res--
		}
		fmt.Fprintln(out, res/2) // 从x到y的路径上的割点个数(除去x,y)
	}

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

type BlockCutTree struct {
	// !BlockCutTree邻接表, 新图的顶点个数为原图割点数+原图点双数.
	// ![0, len(点双)) 对应原图`去除割点`的点双.
	// ![len(点双), len(点双)+len(割点)) 对应原图的割点.
	Tree [][]int

	// BlockCutTree的每个顶点i对应的原图的顶点编号们
	// !割点所在组只有原图割点本身,即Group[i] = []int{rawCut}
	Group [][]int

	idar, idcc []int
	g          [][]Edge
	bcc        *_BCC
	cutCount   int
}

type Edge = struct{ from, to, cost, index int }
type Graph = [][]Edge

func NewBlockCutTree(g Graph) *BlockCutTree {
	return &BlockCutTree{g: g, bcc: _NewBCC(g)}
}

func (bct *BlockCutTree) Build() {
	bct.bcc.Build()
	cut := bct.bcc.lowLink.Articulation
	cutSize, bccSize := len(cut), len(bct.bcc.BCC)
	n := len(bct.g)
	idar, idcc := make([]int, n), make([]int, n)
	last := make([]int, n)
	for i := range idar {
		idar[i], idcc[i] = -1, -1
		last[i] = -1
	}
	for i := 0; i < len(cut); i++ {
		idar[cut[i]] = i + bccSize
	}
	bct.Tree = make([][]int, cutSize+bccSize)
	for i := 0; i < bccSize; i++ {
		st := make([]int, 0, len(bct.bcc.BCC[i])*2)
		for _, e := range bct.bcc.BCC[i] {
			st = append(st, e[0], e[1])
		}
		for _, u := range st {
			if idar[u] == -1 {
				idcc[u] = i
			} else if last[u] != i {
				bct.add(i, idar[u])
				last[u] = i
			}
		}
	}

	bct.idar, bct.idcc = idar, idcc
	bct.cutCount = cutSize

	group := make([][]int, len(bct.Tree))
	for i := 0; i < n; i++ {
		id := bct.Get(i)
		if id >= 0 {
			group[id] = append(group[id], i)
		}
	}
	bct.Group = group
}

// `原图`的顶点对应的BlockCutTree的顶点编号.
// !注意孤立点的编号为-1(不存在于BlockCutTree中).
//  0 <= rawV < len(bct.g)
func (bct *BlockCutTree) Get(rawV int) int {
	if bct.idar[rawV] >= 0 {
		return bct.idar[rawV]
	}
	return bct.idcc[rawV]
}

// `原图`的顶点是否是割点.
//  0 <= rawV < len(bct.g)
func (bct *BlockCutTree) IsRawCut(rawV int) bool { return bct.idar[rawV] >= 0 }

// `圆方树`中的顶点是否是(原图)的割点.
//  0 <= v < len(bct.Tree)
func (bct *BlockCutTree) IsNewCut(v int) bool {
	start := len(bct.bcc.BCC)
	return start <= v && v < start+bct.cutCount
}

// `原图`的顶点是否属于某个点双.
//  0 <= rawV < len(bct.g)
func (bct *BlockCutTree) IsRawVBCC(rawV int) bool { return bct.idar[rawV] < 0 && bct.idcc[rawV] >= 0 }

// `原图`的顶点是否是孤立点.
//  0 <= rawV < len(bct.g)
func (bct *BlockCutTree) IsRawIsolate(rawV int) bool { return bct.idar[rawV] < 0 && bct.idcc[rawV] < 0 }

// 原图中的割点个数.
func (bct *BlockCutTree) CountCut() int { return bct.cutCount }

func (bct *BlockCutTree) add(i, j int) {
	if i == -1 || j == -1 {
		return
	}
	bct.Tree[i] = append(bct.Tree[i], j)
	bct.Tree[j] = append(bct.Tree[j], i)
}

type _BCC struct {
	BCC     [][][2]int // 边:(from,to)
	g       [][]Edge
	lowLink *_lowlink
	used    []bool
	tmp     []Edge
}

func _NewBCC(g [][]Edge) *_BCC {
	return &_BCC{
		g:       g,
		lowLink: _NEewLowlink(g),
	}
}

func (bcc *_BCC) Build() {
	bcc.lowLink.Build()
	bcc.used = make([]bool, len(bcc.g))
	for i := 0; i < len(bcc.used); i++ {
		if !bcc.used[i] {
			bcc.dfs(i, -1)
		}
	}
}

func (bcc *_BCC) dfs(idx, par int) {
	bcc.used[idx] = true
	beet := false
	for _, next := range bcc.g[idx] {
		if next.to == par {
			b := beet
			beet = true
			if !b {
				continue
			}
		}

		if !bcc.used[next.to] || bcc.lowLink.ord[next.to] < bcc.lowLink.ord[idx] {
			bcc.tmp = append(bcc.tmp, next)
		}

		if !bcc.used[next.to] {
			bcc.dfs(next.to, idx)
			if bcc.lowLink.low[next.to] >= bcc.lowLink.ord[idx] {
				bcc.BCC = append(bcc.BCC, [][2]int{})
				for {
					e := bcc.tmp[len(bcc.tmp)-1]
					bcc.BCC[len(bcc.BCC)-1] = append(bcc.BCC[len(bcc.BCC)-1], [2]int{e.from, e.to})
					bcc.tmp = bcc.tmp[:len(bcc.tmp)-1]
					if e.index == next.index {
						break
					}
				}
			}
		}
	}
}

type _lowlink struct {
	Articulation []int // 関節点
	g            [][]Edge
	ord, low     []int
	used         []bool
}

func _NEewLowlink(g [][]Edge) *_lowlink {
	return &_lowlink{g: g}
}

func (ll *_lowlink) Build() {
	ll.used = make([]bool, len(ll.g))
	ll.ord = make([]int, len(ll.g))
	ll.low = make([]int, len(ll.g))
	k := 0
	for i := 0; i < len(ll.g); i++ {
		if !ll.used[i] {
			k = ll.dfs(i, k, -1)
		}
	}
}

func (ll *_lowlink) dfs(idx, k, par int) int {
	ll.used[idx] = true
	ll.ord[idx] = k
	k++
	ll.low[idx] = ll.ord[idx]
	isArticulation := false
	beet := false
	cnt := 0
	for _, e := range ll.g[idx] {
		if e.to == par {
			tmp := beet
			beet = true
			if !tmp {
				continue
			}
		}
		if !ll.used[e.to] {
			cnt++
			k = ll.dfs(e.to, k, idx)
			ll.low[idx] = min(ll.low[idx], ll.low[e.to])
			if par >= 0 && ll.low[e.to] >= ll.ord[idx] {
				isArticulation = true
			}
		} else {
			ll.low[idx] = min(ll.low[idx], ll.ord[e.to])
		}
	}

	if par == -1 && cnt > 1 {
		isArticulation = true
	}
	if isArticulation {
		ll.Articulation = append(ll.Articulation, idx)
	}
	return k
}

// 有时BlockCutTree 需要配合 HLD/Tree 等
//
//
//
type HeavyLightDecomposition struct {
	Tree                          [][]int
	SubSize, Depth, Parent        []int
	dfn, dfnToNode, top, heavySon []int
	dfnId                         int
}

func (hld *HeavyLightDecomposition) Build(root int) {
	hld.build(root, -1, 0)
	hld.markTop(root, root)
}

func NewHeavyLightDecomposition(n int) *HeavyLightDecomposition {
	tree := make([][]int, n)
	dfn := make([]int, n)       // vertex => dfn
	dfnToNode := make([]int, n) // dfn => vertex
	top := make([]int, n)       // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	subSize := make([]int, n)   // 子树大小
	depth := make([]int, n)     // 深度
	parent := make([]int, n)    // 父结点
	heavySon := make([]int, n)  // 重儿子
	return &HeavyLightDecomposition{
		Tree:      tree,
		dfn:       dfn,
		dfnToNode: dfnToNode,
		top:       top,
		SubSize:   subSize,
		Depth:     depth,
		Parent:    parent,
		heavySon:  heavySon,
	}
}

// 添加无向边 u-v.
func (hld *HeavyLightDecomposition) AddEdge(u, v int) {
	hld.Tree[u] = append(hld.Tree[u], v)
	hld.Tree[v] = append(hld.Tree[v], u)
}

// 添加有向边 u->v.
func (hld *HeavyLightDecomposition) AddDirectedEdge(u, v int) {
	hld.Tree[u] = append(hld.Tree[u], v)
}

// 返回树节点 u 对应的 欧拉序区间 [down, up).
//  0 <= down < up <= n.
func (hld *HeavyLightDecomposition) Id(u int) (down, up int) {
	down, up = hld.dfn[u], hld.dfn[u]+hld.SubSize[u]
	return
}

// 返回边 u-v 对应的 欧拉序起点编号.
func (hld *HeavyLightDecomposition) Eid(u, v int) int {
	id1, _ := hld.Id(u)
	id2, _ := hld.Id(v)
	if id1 < id2 {
		return id2
	}
	return id1
}

// 处理路径上的可换操作.
//   0 <= start <= end <= n, [start,end).
func (hld *HeavyLightDecomposition) QueryPath(u, v int, vertex bool, f func(start, end int)) {
	if vertex {
		hld.forEach(u, v, f)
	} else {
		hld.forEachEdge(u, v, f)
	}
}

// 处理以 root 为根的子树上的查询.
//   0 <= start <= end <= n, [start,end).
func (hld *HeavyLightDecomposition) QuerySubTree(u int, vertex bool, f func(start, end int)) {
	if vertex {
		f(hld.dfn[u], hld.dfn[u]+hld.SubSize[u])
	} else {
		f(hld.dfn[u]+1, hld.dfn[u]+hld.SubSize[u])
	}
}

func (hld *HeavyLightDecomposition) forEach(u, v int, cb func(start, end int)) {
	for {
		if hld.dfn[u] > hld.dfn[v] {
			u, v = v, u
		}
		cb(max(hld.dfn[hld.top[v]], hld.dfn[u]), hld.dfn[v]+1)
		if hld.top[u] != hld.top[v] {
			v = hld.Parent[hld.top[v]]
		} else {
			break
		}
	}
}
func (hld *HeavyLightDecomposition) LCA(u, v int) int {
	for {
		if hld.dfn[u] > hld.dfn[v] {
			u, v = v, u
		}
		if hld.top[u] == hld.top[v] {
			return u
		}
		v = hld.Parent[hld.top[v]]
	}
}

func (hld *HeavyLightDecomposition) Dist(u, v int) int {
	return hld.Depth[u] + hld.Depth[v] - 2*hld.Depth[hld.LCA(u, v)]
}

// 寻找以 start 为 top 的重链 ,heavyPath[-1] 即为重链末端节点.
func (hld *HeavyLightDecomposition) GetHeavyPath(start int) []int {
	heavyPath := []int{start}
	cur := start
	for hld.heavySon[cur] != -1 {
		cur = hld.heavySon[cur]
		heavyPath = append(heavyPath, cur)
	}
	return heavyPath
}

func (hld *HeavyLightDecomposition) forEachEdge(u, v int, cb func(start, end int)) {
	for {
		if hld.dfn[u] > hld.dfn[v] {
			u, v = v, u
		}
		if hld.top[u] != hld.top[v] {
			cb(hld.dfn[hld.top[v]], hld.dfn[v]+1)
			v = hld.Parent[hld.top[v]]
		} else {
			if u != v {
				cb(hld.dfn[u]+1, hld.dfn[v]+1)
			}
			break
		}
	}
}

func (hld *HeavyLightDecomposition) build(cur, pre, dep int) int {
	subSize, heavySize, heavySon := 1, 0, -1
	for _, next := range hld.Tree[cur] {
		if next != pre {
			nextSize := hld.build(next, cur, dep+1)
			subSize += nextSize
			if nextSize > heavySize {
				heavySize, heavySon = nextSize, next
			}
		}
	}
	hld.Depth[cur] = dep
	hld.SubSize[cur] = subSize
	hld.heavySon[cur] = heavySon
	hld.Parent[cur] = pre
	return subSize
}

func (hld *HeavyLightDecomposition) markTop(cur, top int) {
	hld.top[cur] = top
	hld.dfn[cur] = hld.dfnId
	hld.dfnId++
	hld.dfnToNode[hld.dfn[cur]] = cur
	if hld.heavySon[cur] != -1 {
		hld.markTop(hld.heavySon[cur], top)
		for _, next := range hld.Tree[cur] {
			if next != hld.heavySon[cur] && next != hld.Parent[cur] {
				hld.markTop(next, next)
			}
		}
	}
}
