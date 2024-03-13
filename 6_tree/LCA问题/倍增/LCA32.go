package main

import "fmt"

func main() {
	edges := [][]int32{{0, 1}, {1, 2}, {1, 3}, {0, 4}, {3, 5}, {3, 6}, {5, 7}}
	tree := make([][]int32, 8)
	for _, e := range edges {
		u, v := e[0], e[1]
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}
	lca := NewLCA32(tree, []int32{0})
	fmt.Println(lca.LCAMultiPoint([]int32{5, 6, 7, 4})) // 0
}

type LCA32 struct {
	Depth, Parent           []int32
	Tree                    [][]int32
	lid, rid, top, heavySon []int32
	idToNode                []int32
	dfnId                   int32
}

func NewLCA32(tree [][]int32, roots []int32) *LCA32 {
	n := len(tree)
	lid := make([]int32, n) // vertex => dfn
	rid := make([]int32, n)
	top := make([]int32, n)      // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	depth := make([]int32, n)    // 深度
	parent := make([]int32, n)   // 父结点
	heavySon := make([]int32, n) // 重儿子
	idToNode := make([]int32, n)
	res := &LCA32{
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
		res._build(root, -1, 0)
		res._markTop(root, root)
	}
	return res
}

// 树上多个点的 LCA，就是 DFS 序最小的和 DFS 序最大的这两个点的 LCA.
func (hld *LCA32) LCAMultiPoint(nodes []int32) int32 {
	if len(nodes) == 1 {
		return nodes[0]
	}
	if len(nodes) == 2 {
		return hld.LCA(nodes[0], nodes[1])
	}
	minDfn, maxDfn := int32(1<<31-1), int32(-1)
	for _, root := range nodes {
		if hld.lid[root] < minDfn {
			minDfn = hld.lid[root]
		}
		if hld.lid[root] > maxDfn {
			maxDfn = hld.lid[root]
		}
	}
	u, v := hld.idToNode[minDfn], hld.idToNode[maxDfn]
	return hld.LCA(u, v)
}

func (hld *LCA32) LCA(u, v int32) int32 {
	for {
		if hld.lid[u] > hld.lid[v] {
			u, v = v, u
		}
		if hld.top[u] == hld.top[v] {
			return u
		}
		v = hld.Parent[hld.top[v]]
	}
}

func (hld *LCA32) Dist(u, v int32) int32 {
	return hld.Depth[u] + hld.Depth[v] - 2*hld.Depth[hld.LCA(u, v)]
}

func (hld *LCA32) EnumeratePathDecomposition(u, v int32, vertex bool, f func(start, end int32)) {
	for {
		if hld.top[u] == hld.top[v] {
			break
		}
		if hld.lid[u] < hld.lid[v] {
			a, b := hld.lid[hld.top[v]], hld.lid[v]
			if a > b {
				a, b = b, a
			}
			f(a, b+1)
			v = hld.Parent[hld.top[v]]
		} else {
			a, b := hld.lid[u], hld.lid[hld.top[u]]
			if a > b {
				a, b = b, a
			}
			f(a, b+1)
			u = hld.Parent[hld.top[u]]
		}
	}

	edgeInt := int32(1)
	if vertex {
		edgeInt = 0
	}

	if hld.lid[u] < hld.lid[v] {
		a, b := hld.lid[u]+edgeInt, hld.lid[v]
		if a > b {
			a, b = b, a
		}
		f(a, b+1)
	} else if hld.lid[v]+edgeInt <= hld.lid[u] {
		a, b := hld.lid[u], hld.lid[v]+edgeInt
		if a > b {
			a, b = b, a
		}
		f(a, b+1)
	}
}

// k: 0-based
//
//	如果不存在第k个祖先，返回-1
func (hld *LCA32) KthAncestor(root, k int32) int32 {
	if k > hld.Depth[root] {
		return -1
	}
	for {
		u := hld.top[root]
		if hld.lid[root]-k >= hld.lid[u] {
			return hld.idToNode[hld.lid[root]-k]
		}
		k -= hld.lid[root] - hld.lid[u] + 1
		root = hld.Parent[u]
	}
}

// 从 from 节点跳向 to 节点,跳过 step 个节点(0-indexed)
//
//	返回跳到的节点,如果不存在这样的节点,返回-1
func (hld *LCA32) Jump(from, to, step int32) int32 {
	if step == 1 {
		if from == to {
			return -1
		}
		if hld.IsInSubtree(to, from) {
			return hld.KthAncestor(to, (hld.Depth[to] - hld.Depth[from] - 1))
		}
		return hld.Parent[from]
	}
	c := hld.LCA(from, to)
	dac := hld.Depth[from] - hld.Depth[c]
	dbc := hld.Depth[to] - hld.Depth[c]
	if step > dac+dbc {
		return -1
	}
	if step <= dac {
		return hld.KthAncestor(from, step)
	}
	return hld.KthAncestor(to, dac+dbc-step)
}

// child 是否在 root 的子树中 (child和root不能相等)
func (hld *LCA32) IsInSubtree(child, root int32) bool {
	return hld.lid[root] <= hld.lid[child] && hld.lid[child] < hld.rid[root]
}

func (hld *LCA32) _build(cur, pre, dep int32) int32 {
	subSize, heavySize, heavySon := int32(1), int32(0), int32(-1)
	for _, next := range hld.Tree[cur] {
		if next != pre {
			nextSize := hld._build(next, cur, dep+1)
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

func (hld *LCA32) _markTop(cur, top int32) {
	hld.top[cur] = top
	hld.lid[cur] = hld.dfnId
	hld.idToNode[hld.dfnId] = cur
	hld.dfnId++
	if hld.heavySon[cur] != -1 {
		hld._markTop(hld.heavySon[cur], top)
		for _, next := range hld.Tree[cur] {
			if next != hld.heavySon[cur] && next != hld.Parent[cur] {
				hld._markTop(next, next)
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
