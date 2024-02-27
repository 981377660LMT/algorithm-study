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
	fmt.Println(lca.LCAMultiPoint([]int{5, 6, 7, 4})) // 0
}

type LCAFast struct {
	Depth, Parent           []int32
	Tree                    [][]int
	lid, rid, top, heavySon []int32
	idToNode                []int32
	dfnId                   int32
}

func NewLCA(tree [][]int, roots []int) *LCAFast {
	n := len(tree)
	lid := make([]int32, n) // vertex => dfn
	rid := make([]int32, n)
	top := make([]int32, n)      // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	depth := make([]int32, n)    // 深度
	parent := make([]int32, n)   // 父结点
	heavySon := make([]int32, n) // 重儿子
	idToNode := make([]int32, n)
	res := &LCAFast{
		Tree:     tree,
		lid:      lid,
		rid:      rid,
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

// 树上多个点的 LCA，就是 DFS 序最小的和 DFS 序最大的这两个点的 LCA.
func (hld *LCAFast) LCAMultiPoint(nodes []int) int {
	if len(nodes) == 1 {
		return nodes[0]
	}
	if len(nodes) == 2 {
		return hld.LCA(nodes[0], nodes[1])
	}
	minDfn, maxDfn := int32(1<<31-1), int32(-1)
	for _, root := range nodes {
		root32 := int32(root)
		if hld.lid[root32] < minDfn {
			minDfn = hld.lid[root32]
		}
		if hld.lid[root32] > maxDfn {
			maxDfn = hld.lid[root32]
		}
	}
	u, v := hld.idToNode[minDfn], hld.idToNode[maxDfn]
	return hld.LCA(int(u), int(v))
}

func (hld *LCAFast) LCA(u, v int) int {
	u32, v32 := int32(u), int32(v)
	for {
		if hld.lid[u32] > hld.lid[v32] {
			u32, v32 = v32, u32
		}
		if hld.top[u32] == hld.top[v32] {
			return int(u32)
		}
		v32 = hld.Parent[hld.top[v32]]
	}
}

func (hld *LCAFast) Dist(u, v int) int {
	return int(hld.Depth[u] + hld.Depth[v] - 2*hld.Depth[hld.LCA(u, v)])
}

func (hld *LCAFast) EnumeratePathDecomposition(u, v int, vertex bool, f func(start, end int)) {
	u32, v32 := int32(u), int32(v)
	for {
		if hld.top[u32] == hld.top[v32] {
			break
		}
		if hld.lid[u32] < hld.lid[v32] {
			a, b := hld.lid[hld.top[v32]], hld.lid[v32]
			if a > b {
				a, b = b, a
			}
			f(int(a), int(b+1))
			v32 = hld.Parent[hld.top[v32]]
		} else {
			a, b := hld.lid[u32], hld.lid[hld.top[u32]]
			if a > b {
				a, b = b, a
			}
			f(int(a), int(b+1))
			u32 = hld.Parent[hld.top[u32]]
		}
	}

	edgeInt := int32(1)
	if vertex {
		edgeInt = 0
	}

	if hld.lid[u32] < hld.lid[v32] {
		a, b := hld.lid[u32]+edgeInt, hld.lid[v32]
		if a > b {
			a, b = b, a
		}
		f(int(a), int(b+1))
	} else if hld.lid[v32]+edgeInt <= hld.lid[u32] {
		a, b := hld.lid[u32], hld.lid[v32]+edgeInt
		if a > b {
			a, b = b, a
		}
		f(int(a), int(b+1))
	}
}

// k: 0-based
//
//	如果不存在第k个祖先，返回-1
func (hld *LCAFast) KthAncestor(root, k int) int {
	root32 := int32(root)
	k32 := int32(k)
	if k32 > hld.Depth[root32] {
		return -1
	}
	for {
		u := hld.top[root32]
		if hld.lid[root32]-k32 >= hld.lid[u] {
			return int(hld.idToNode[hld.lid[root32]-k32])
		}
		k32 -= hld.lid[root32] - hld.lid[u] + 1
		root32 = hld.Parent[u]
	}
}

// 从 from 节点跳向 to 节点,跳过 step 个节点(0-indexed)
//
//	返回跳到的节点,如果不存在这样的节点,返回-1
func (hld *LCAFast) Jump(from, to, step int) int {
	if step == 1 {
		if from == to {
			return -1
		}
		if hld.IsInSubtree(to, from) {
			return hld.KthAncestor(to, int(hld.Depth[to]-hld.Depth[from]-1))
		}
		return int(hld.Parent[from])
	}
	c := hld.LCA(from, to)
	dac := hld.Depth[from] - hld.Depth[c]
	dbc := hld.Depth[to] - hld.Depth[c]
	if step > int(dac+dbc) {
		return -1
	}
	if step <= int(dac) {
		return hld.KthAncestor(from, step)
	}
	return hld.KthAncestor(to, int(dac+dbc-int32(step)))
}

// child 是否在 root 的子树中 (child和root不能相等)
func (hld *LCAFast) IsInSubtree(child, root int) bool {
	return hld.lid[root] <= hld.lid[child] && hld.lid[child] < hld.rid[root]
}

func (hld *LCAFast) _build(cur, pre, dep int32) int {
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

func (hld *LCAFast) _markTop(cur, top int32) {
	hld.top[cur] = top
	hld.lid[cur] = hld.dfnId
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
	hld.rid[cur] = hld.dfnId
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
