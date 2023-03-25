// https://ei1333.github.io/library/graph/others/namori-graph.hpp
// Namori Graph 无向图基环树

// !一个有 n 个顶点和 n 条边的“连通”无向图。图中只有一个环。
// 这个图被称为“Namori图”，这是根据一位漫画家的名字而命名的，
// 但在学术界中，正确的称呼是“单环图(Unicyclic)”或“伪森林(Pseudoforest)”。
// 在这里，将该图分解为一个环和附加在各个环上顶点的树。
// !将包含在环中的顶点重新编号为 [0, k)，其中 k 是环中的顶点数，对应各个treeRoot.
// !同样，附加的树也被重新编号为 [0, l)，其中 l 是树的顶点数, 0 是树的根。

// NewNamoriGraph(g Graph)。
// Build()：构建基环树和各个环上的子树。
// GetId(rawV int) (rootId, idInTree int)：
//   给定原图的顶点rawV,返回rawV所在的树的根节点和rawV在树中的编号.
// GetInvId(rootId, idInTree int) (rawV int)：
//   给定树的顶点编号root和某个点在树中的顶点编号idInTree, 返回这个点在原图中的顶点编号.
//   GetInvId(root,0) 表示在环上的顶点root在原图中对应的顶点.

// Trees []Graph : !以环上各个顶点i为根的无向树
// CycleEdges []Edge : !基环树中在环上的边,边i连接着树的 root i和i+1 (i>=0)
// HLDs []*HeavyLightDecomposition : !各个树的HLD

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	namoriCut()
}

func yuki1254() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	g := make([][]Edge, n)
	for i := 0; i < n; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		g[a] = append(g[a], Edge{a, b, 1, i})
		g[b] = append(g[b], Edge{b, a, 1, i})
	}

	G := NewNamoriGraph(g)
	G.Build(false)
	res := []int{}
	for _, e := range G.CycleEdges {
		res = append(res, e.id+1)
	}
	sort.Ints(res)
	fmt.Fprintln(out, len(res))
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

func abc266_f() {
	// https://atcoder.jp/contests/abc266/tasks/abc266_f
	// 给定一个基环树,问从x到y的路径是否唯一
	// 等价于不能走环上=>在同一个子树中
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	graph := make([][]Edge, n)
	for i := 0; i < n; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		graph[a] = append(graph[a], Edge{a, b, 1, i})
		graph[b] = append(graph[b], Edge{b, a, 1, i})
	}

	G := NewNamoriGraph(graph)
	G.Build(false) // !without HLD

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		x--
		y--
		root1, _ := G.GetId(x)
		root2, _ := G.GetId(y)
		if root1 == root2 {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}

func namoriCut() {
	// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=2891
	// 给定一个基环树 q个询问(x,y)
	// 求使得x和y不连通最少需要切掉多少条边
	// 如果两个点都在环上,则答案为2
	// 否则答案为1(两点之间只有唯一的一条路径)

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	graph := make([][]Edge, n)
	for i := 0; i < n; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a, b = a-1, b-1
		graph[a] = append(graph[a], Edge{a, b, 1, i})
		graph[b] = append(graph[b], Edge{b, a, 1, i})
	}
	G := NewNamoriGraph(graph)
	G.Build(false)

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		x, y = x-1, y-1
		_, treeId1 := G.GetId(x)
		_, treeId2 := G.GetId(y)
		if treeId1 != 0 || treeId2 != 0 {
			fmt.Fprintln(out, 1)
		} else {
			fmt.Fprintln(out, 2)
		}
	}
}

type NamoriGraph struct {
	// !以环上各个顶点i为根的无向树
	Trees []Graph
	// !基环树中在环上的边,边i连接着树的 root i和i+1 (i>=0)
	CycleEdges []Edge

	// 每个树的重链剖分(需要在Build(needHLD=true)后才能使用)
	HLDs []*_HLD

	g          Graph
	iv         [][]int
	markId, id []int
}

type Edge = struct{ from, to, cost, id int }
type Graph = [][]Edge

func NewNamoriGraph(g Graph) *NamoriGraph {
	return &NamoriGraph{g: g}
}

// needHLD :是否需要对各个子树进行重链剖分.
func (ng *NamoriGraph) Build(needHLD bool) {
	n := len(ng.g)
	deg := make([]int, n)
	used := make([]bool, n)
	que := []int{}
	for i := 0; i < n; i++ {
		deg[i] = len(ng.g[i])
		if deg[i] == 1 {
			que = append(que, i)
			used[i] = true
		}
	}

	for len(que) > 0 {
		idx := que[0]
		que = que[1:]
		for _, e := range ng.g[idx] {
			if used[e.to] {
				continue
			}
			deg[e.to]--
			if deg[e.to] == 1 {
				que = append(que, e.to)
				used[e.to] = true
			}
		}
	}

	mx := 0
	for _, edges := range ng.g {
		for _, e := range edges {
			mx = max(mx, e.id)
		}
	}

	edgeUsed := make([]bool, mx+1)
	loop := []int{}
	for i := 0; i < n; i++ {
		if used[i] {
			continue
		}
		for update := true; update; {
			update = false
			loop = append(loop, i)
			for _, e := range ng.g[i] {
				if used[e.to] || edgeUsed[e.id] {
					continue
				}
				edgeUsed[e.id] = true
				ng.CycleEdges = append(ng.CycleEdges, Edge{i, e.to, e.cost, e.id})
				i = e.to
				update = true
				break
			}
		}
		break
	}

	loop = loop[:len(loop)-1]
	ng.markId = make([]int, n)
	ng.id = make([]int, n)
	for i := 0; i < len(loop); i++ {
		pre := loop[(i+len(loop)-1)%len(loop)]
		nxt := loop[(i+1)%len(loop)]
		sz := 0
		ng.markId[loop[i]] = i
		ng.iv = append(ng.iv, []int{})
		ng.id[loop[i]] = sz
		sz++
		ng.iv[len(ng.iv)-1] = append(ng.iv[len(ng.iv)-1], loop[i])
		for _, e := range ng.g[loop[i]] {
			if e.to != pre && e.to != nxt {
				ng.markDfs(e.to, loop[i], i, &sz)
			}
		}
		tree := make(Graph, sz)
		for _, e := range ng.g[loop[i]] {
			if e.to != pre && e.to != nxt {
				tree[ng.id[loop[i]]] = append(tree[ng.id[loop[i]]], Edge{ng.id[loop[i]], ng.id[e.to], e.cost, e.id})
				tree[ng.id[e.to]] = append(tree[ng.id[e.to]], Edge{ng.id[e.to], ng.id[loop[i]], e.cost, e.id})
				ng.buildDfs(e.to, loop[i], tree)
			}
		}
		ng.Trees = append(ng.Trees, tree)
	}

	// HLD
	if !needHLD {
		return
	}

	t := len(ng.Trees)
	ng.HLDs = make([]*_HLD, 0, t)
	for _, tree := range ng.Trees {
		hld := _NewHLD(tree)
		hld.Build(0)
		ng.HLDs = append(ng.HLDs, hld)
	}
}

// 给定原图的顶点rawV,返回rawV所在的树的根节点和rawV在树中的编号.
func (ng *NamoriGraph) GetId(rawV int) (rootId, idInTree int) {
	return ng.markId[rawV], ng.id[rawV]
}

// 给定树的顶点编号root和某个点在树中的顶点编号idInTree,返回这个点在原图中的顶点编号.
//  GetInvId(root,0) 表示在环上的顶点root在原图中对应的顶点.
func (ng *NamoriGraph) GetInvId(rootId, idInTree int) (rawV int) {
	return ng.iv[rootId][idInTree]
}

func (ng *NamoriGraph) markDfs(idx, par, k int, l *int) {
	ng.markId[idx] = k
	ng.id[idx] = *l
	*l++
	ng.iv[len(ng.iv)-1] = append(ng.iv[len(ng.iv)-1], idx)
	for _, e := range ng.g[idx] {
		if e.to != par {
			ng.markDfs(e.to, idx, k, l)
		}
	}
}

func (ng *NamoriGraph) buildDfs(idx, par int, tree Graph) {
	for _, e := range ng.g[idx] {
		if e.to != par {
			tree[ng.id[idx]] = append(tree[ng.id[idx]], Edge{ng.id[idx], ng.id[e.to], e.cost, e.id})
			tree[ng.id[e.to]] = append(tree[ng.id[e.to]], Edge{ng.id[e.to], ng.id[idx], e.cost, e.id})
			ng.buildDfs(e.to, idx, tree)
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type _HLD struct {
	Tree                          Graph
	SubSize, Depth, Parent        []int
	dfn, dfnToNode, top, heavySon []int
	dfnId                         int
}

func (hld *_HLD) Build(root int) {
	hld.build(root, -1, 0)
	hld.markTop(root, root)
}

func _NewHLD(tree Graph) *_HLD {
	n := len(tree)
	dfn := make([]int, n)       // vertex => dfn
	dfnToNode := make([]int, n) // dfn => vertex
	top := make([]int, n)       // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	subSize := make([]int, n)   // 子树大小
	depth := make([]int, n)     // 深度
	parent := make([]int, n)    // 父结点
	heavySon := make([]int, n)  // 重儿子
	return &_HLD{
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

// 返回树节点 u 对应的 欧拉序区间 [down, up).
//  0 <= down < up <= n.
func (hld *_HLD) Id(u int) (down, up int) {
	down, up = hld.dfn[u], hld.dfn[u]+hld.SubSize[u]
	return
}

// 返回边 u-v 对应的 欧拉序起点编号.
func (hld *_HLD) Eid(u, v int) int {
	id1, _ := hld.Id(u)
	id2, _ := hld.Id(v)
	if id1 < id2 {
		return id2
	}
	return id1
}

// 处理路径上的可换操作.
//   0 <= start <= end <= n, [start,end).
func (hld *_HLD) QueryPath(u, v int, vertex bool, f func(start, end int)) {
	if vertex {
		hld.forEach(u, v, f)
	} else {
		hld.forEachEdge(u, v, f)
	}
}

// 处理以 root 为根的子树上的查询.
//   0 <= start <= end <= n, [start,end).
func (hld *_HLD) QuerySubTree(u int, vertex bool, f func(start, end int)) {
	if vertex {
		f(hld.dfn[u], hld.dfn[u]+hld.SubSize[u])
	} else {
		f(hld.dfn[u]+1, hld.dfn[u]+hld.SubSize[u])
	}
}

func (hld *_HLD) forEach(u, v int, cb func(start, end int)) {
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

func (hld *_HLD) LCA(u, v int) int {
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

func (hld *_HLD) Dist(u, v int) int {
	return hld.Depth[u] + hld.Depth[v] - 2*hld.Depth[hld.LCA(u, v)]
}

// 寻找以 start 为 top 的重链 ,heavyPath[-1] 即为重链末端节点.
func (hld *_HLD) GetHeavyPath(start int) []int {
	heavyPath := []int{start}
	cur := start
	for hld.heavySon[cur] != -1 {
		cur = hld.heavySon[cur]
		heavyPath = append(heavyPath, cur)
	}
	return heavyPath
}

func (hld *_HLD) forEachEdge(u, v int, cb func(start, end int)) {
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

func (hld *_HLD) build(cur, pre, dep int) int {
	subSize, heavySize, heavySon := 1, 0, -1
	for _, e := range hld.Tree[cur] {
		next := e.to
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

func (hld *_HLD) markTop(cur, top int) {
	hld.top[cur] = top
	hld.dfn[cur] = hld.dfnId
	hld.dfnId++
	hld.dfnToNode[hld.dfn[cur]] = cur
	if hld.heavySon[cur] != -1 {
		hld.markTop(hld.heavySon[cur], top)
		for _, e := range hld.Tree[cur] {
			next := e.to
			if next != hld.heavySon[cur] && next != hld.Parent[cur] {
				hld.markTop(next, next)
			}
		}
	}
}
