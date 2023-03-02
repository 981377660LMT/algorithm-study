// 重链剖分求LCA 速度快

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
	var q int
	fmt.Fscan(in, &q)
	lca := NewLCA(n)
	for i := 1; i < n; i++ {
		var parent int
		fmt.Fscan(in, &parent)
		lca.AddEdge(i, parent)
	}
	lca.Build(0)

	for i := 0; i < q; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		fmt.Fprintln(out, lca.LCA(u, v))
	}
}

type LCAHLD struct {
	Depth, Parent      []int
	Tree               [][]int
	dfn, top, heavySon []int
	dfnId              int
}

func NewLCA(n int) *LCAHLD {
	tree := make([][]int, n)
	dfn := make([]int, n)      // vertex => dfn
	top := make([]int, n)      // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	depth := make([]int, n)    // 深度
	parent := make([]int, n)   // 父结点
	heavySon := make([]int, n) // 重儿子
	return &LCAHLD{
		Tree:     tree,
		dfn:      dfn,
		top:      top,
		Depth:    depth,
		Parent:   parent,
		heavySon: heavySon,
	}
}

// 添加无向边 u-v.
func (hld *LCAHLD) AddEdge(u, v int) {
	hld.Tree[u] = append(hld.Tree[u], v)
	hld.Tree[v] = append(hld.Tree[v], u)
}

func (hld *LCAHLD) Build(root int) {
	hld.build(root, -1, 0)
	hld.markTop(root, root)
}

func (hld *LCAHLD) LCA(u, v int) int {
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

func (hld *LCAHLD) Dist(u, v int) int {
	return hld.Depth[u] + hld.Depth[v] - 2*hld.Depth[hld.LCA(u, v)]
}

func (hld *LCAHLD) build(cur, pre, dep int) int {
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
	hld.heavySon[cur] = heavySon
	hld.Parent[cur] = pre
	return subSize
}

func (hld *LCAHLD) markTop(cur, top int) {
	hld.top[cur] = top
	hld.dfn[cur] = hld.dfnId
	hld.dfnId++
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
