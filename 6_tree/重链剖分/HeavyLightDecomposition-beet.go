// !比 `HeavyLightDecomposition-nyann.go` 快一些

// https://beet-aizu.github.io/library/tree/heavylightdecomposition.cpp
// HL分解将树上的路径分成logn条,分割之后只需要op操作logn条链即可
// 如果原问题可以在序列上O(X)时间解决,那么在树上就可以在O(Xlogn)时间解决
// !如果op运算不满足交换律,需要使用w=lca(u,v)过渡,合成forEach(w,u)和forEach(w,v)的结果

// Usage:
//  hld := NewHeavyLightDecomposition(n)
//  for i := 0; i < n-1; i++ {
//      hld.AddEdge(u, v)
//  }
//
//  hld.Build(root)
//  hld.QueryPath(u, v, vertex, f)
//  hld.QueryNonCommutativePath(u, v, vertex, f)
//  hld.QuerySubTree(u, vertex, f)
//  hld.Id(u)
//  hld.LCA(u, v)

// 树的欧拉序编号:
//             0 [0,7)
//					 /       \
//					/         \
//   		e1 /           \ e4
//        /             \
//       /	             \
//      1 [1,4)        	  2 [4,7)
//     / \               / \
// e2 /	  \ e3       e5 /   \ e6
//   /	   \           /     \
// 3 [2,3)  4 [3,4)  5 [5,6)  6 [6,7)

// 点的表示(0 <= vid <= n-1):
//   一个点的起点终点用欧拉序[down,up) (0 <= down < up <= n) 表示.
//   !点[down,up)的编号为 `down`.

// 边的表示(1 <= eid <= n-1):
//   !边的序号用较深的那个顶点的欧拉序的起点编号表示.
//   如上图, 0 -> 1 的边表示为 1, 1 -> 4 的边表示为 3
//   !点 [in,out) 到父亲的边的序号为 `in`.

package main

import "fmt"

func main() {
	hld := NewHeavyLightDecomposition(7)
	hld.AddEdge(0, 1)
	hld.AddEdge(0, 2)
	hld.AddEdge(1, 3)
	hld.AddEdge(1, 4)
	hld.AddEdge(2, 5)
	hld.AddEdge(2, 6)
	hld.Build(0)
	for i := 0; i < 7; i++ {
		fmt.Println(hld.Id(i))
	}

	// 查询边权
	hld.QueryPath(3, 6, false, func(start, end int) {
		fmt.Println("path1", start, end)
		// path1 6 7
		// path1 4 5
		// path1 1 3
	})
	hld.QueryPath(2, 5, false, func(start, end int) {
		fmt.Println("path2", start, end)
		// path2 5 6
	})
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
