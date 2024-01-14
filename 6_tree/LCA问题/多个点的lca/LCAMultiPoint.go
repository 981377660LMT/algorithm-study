package main

import "fmt"

func main() {
	edges := [][]int{{0, 1}, {1, 2}, {1, 3}, {0, 4}, {3, 5}, {3, 6}, {5, 7}}
	tree := make([][]int, 8)
	for _, e := range edges {
		u, v := e[0], e[1]
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}
	lca := NewLCA(tree, []int{0})
	fmt.Println(lca.LCAMultiPoint([]int{5, 6, 7})) // 0
}

// 树上多个点的 LCA，就是 DFS 序最小的和 DFS 序最大的这两个点的 LCA.
type LCAMultiPoint struct {
	Depth, Parent      []int32
	Tree               [][]int
	dfn, top, heavySon []int32
	idToNode           []int32
	dfnId              int32
}

func NewLCA(tree [][]int, roots []int) *LCAMultiPoint {
	n := len(tree)
	dfn := make([]int32, n)      // vertex => dfn
	top := make([]int32, n)      // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	depth := make([]int32, n)    // 深度
	parent := make([]int32, n)   // 父结点
	heavySon := make([]int32, n) // 重儿子
	idToNode := make([]int32, n)
	res := &LCAMultiPoint{
		Tree:     tree,
		dfn:      dfn,
		top:      top,
		Depth:    depth,
		Parent:   parent,
		heavySon: heavySon,
		idToNode: idToNode,
	}
	for _, root := range roots {
		root32 := int32(root)
		res._build(root32, -1, 0)
		res._markTop(root32, root32)
	}
	return res
}

func (hld *LCAMultiPoint) LCAMultiPoint(nodes []int) int {
	if len(nodes) == 1 {
		return nodes[0]
	}
	if len(nodes) == 2 {
		return hld.LCA(nodes[0], nodes[1])
	}
	minDfn, maxDfn := int32(1<<31-1), int32(-1)
	for _, root := range nodes {
		root32 := int32(root)
		if hld.dfn[root32] < minDfn {
			minDfn = hld.dfn[root32]
		}
		if hld.dfn[root32] > maxDfn {
			maxDfn = hld.dfn[root32]
		}
	}
	u, v := hld.idToNode[minDfn], hld.idToNode[maxDfn]
	return hld.LCA(int(u), int(v))
}

func (hld *LCAMultiPoint) LCA(u, v int) int {
	u32, v32 := int32(u), int32(v)
	for {
		if hld.dfn[u32] > hld.dfn[v32] {
			u32, v32 = v32, u32
		}
		if hld.top[u32] == hld.top[v32] {
			return int(u32)
		}
		v32 = hld.Parent[hld.top[v32]]
	}
}

func (hld *LCAMultiPoint) Dist(u, v int) int {
	return int(hld.Depth[u] + hld.Depth[v] - 2*hld.Depth[hld.LCA(u, v)])
}

func (hld *LCAMultiPoint) _build(cur, pre, dep int32) int {
	subSize, heavySize, heavySon := 1, 0, int32(-1)
	for _, next := range hld.Tree[cur] {
		next32 := int32(next)
		if next32 != pre {
			nextSize := hld._build(next32, cur, dep+1)
			subSize += nextSize
			if nextSize > heavySize {
				heavySize, heavySon = nextSize, next32
			}
		}
	}
	hld.Depth[cur] = dep
	hld.heavySon[cur] = heavySon
	hld.Parent[cur] = pre
	return subSize
}

func (hld *LCAMultiPoint) _markTop(cur, top int32) {
	hld.top[cur] = top
	hld.dfn[cur] = hld.dfnId
	hld.idToNode[hld.dfnId] = cur
	hld.dfnId++
	if hld.heavySon[cur] != -1 {
		hld._markTop(hld.heavySon[cur], top)
		for _, next := range hld.Tree[cur] {
			next32 := int32(next)
			if next32 != hld.heavySon[cur] && next32 != hld.Parent[cur] {
				hld._markTop(next32, next32)
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
