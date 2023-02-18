package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://judge.yosupo.jp/problem/lca

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	hld := NewHeavyLightDecomposition(n)
	for i := 1; i < n; i++ {
		var parent int
		fmt.Fscan(in, &parent)
		hld.AddDirectedEdge(parent, i)
	}
	hld.Build(0)

	for i := 0; i < q; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		fmt.Fprintln(out, hld.LCA(u, v))
	}
}

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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
