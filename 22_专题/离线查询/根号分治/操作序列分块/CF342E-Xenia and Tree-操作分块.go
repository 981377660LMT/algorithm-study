// https://www.luogu.com.cn/problem/CF342E
// https://www.luogu.com.cn/problem/solution/P5443

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// https://www.luogu.com.cn/problem/CF342E
// 给定一棵 n 个节点的树，初始时 1 号节点为红色，其余为蓝色。
// 要求支持如下操作：
// 1.将一个节点变为红色。
// 2.询问节点 u 到最近红色节点的距离。
// !n,q<=1e5
//
// 我们设定一个阈值S，对操作序列每S个操作分一块。
// 对于每个块内的询问，我们对红点分类讨论：
// !1. 对于这个询问块之前的红点，我们预处理出来这些红点对于询问的点的距离；适合红点多的情况。
// !2. 对于在这个询问块内的红点，我们暴力枚举然后求距离。适合红点少的情况。
func XeniaAndTree(n int, edges [][2]int, operations [][2]int) (res []int) {
	adjList := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adjList[u] = append(adjList[u], v)
		adjList[v] = append(adjList[v], u)
	}

	preRed := []int{0} // 初始时根节点为红色
	distToPreRed := make([]int, n)
	updateDistToPreRed := func(starts []int) { // bfs, updateWithPreBlockMutations
		queue := append([]int{}, starts...)
		for i := range distToPreRed {
			distToPreRed[i] = -1
		}
		for _, s := range queue {
			distToPreRed[s] = 0
		}
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			for _, v := range adjList[u] {
				if distToPreRed[v] == -1 {
					distToPreRed[v] = distToPreRed[u] + 1
					queue = append(queue, v)
				}
			}
		}
	}

	lca := NewLCA(adjList, []int{0})

	// 操作序列分块.
	q := len(operations)
	block := UseBlock(q, 3*int(math.Sqrt(float64(q))+1)) // 减少分块个数
	blockStart, blockEnd, blockCount := block.blockStart, block.blockEnd, block.blockCount
	for bid := 0; bid < blockCount; bid++ {
		updateDistToPreRed(preRed)
		curRed := []int{}

		for i := blockStart[bid]; i < blockEnd[bid]; i++ {
			kind, u := operations[i][0], operations[i][1]
			if kind == 1 {
				curRed = append(curRed, u)
			} else {
				dist := distToPreRed[u]
				for _, v := range curRed {
					dist = min(dist, lca.Dist(u, v))
				}
				res = append(res, dist)
			}
		}

		for _, u := range curRed {
			preRed = append(preRed, u)
		}
	}

	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	edges := make([][2]int, n-1)
	for i := 0; i < n-1; i++ {
		fmt.Fscan(in, &edges[i][0], &edges[i][1])
		edges[i][0]--
		edges[i][1]--
	}
	ops := make([][2]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &ops[i][0], &ops[i][1])
		ops[i][1]--
	}

	res := XeniaAndTree(n, edges, ops)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type LCAFast struct {
	Depth, Parent      []int32
	Tree               [][]int
	dfn, top, heavySon []int32
	idToNode           []int32
	dfnId              int32
}

func NewLCA(tree [][]int, roots []int) *LCAFast {
	n := len(tree)
	dfn := make([]int32, n)      // vertex => dfn
	top := make([]int32, n)      // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	depth := make([]int32, n)    // 深度
	parent := make([]int32, n)   // 父结点
	heavySon := make([]int32, n) // 重儿子
	idToNode := make([]int32, n)
	res := &LCAFast{
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

func (hld *LCAFast) LCA(u, v int) int {
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

func (hld *LCAFast) Dist(u, v int) int {
	return int(hld.Depth[u] + hld.Depth[v] - 2*hld.Depth[hld.LCA(u, v)])
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

// blockSize = int(math.Sqrt(float64(len(nums)))+1)
func UseBlock(n int, blockSize int) struct {
	belong     []int // 下标所属的块.
	blockStart []int // 每个块的起始下标(包含).
	blockEnd   []int // 每个块的结束下标(不包含).
	blockCount int   // 块的数量.
} {
	blockCount := 1 + (n / blockSize)
	blockStart := make([]int, blockCount)
	blockEnd := make([]int, blockCount)
	belong := make([]int, n)
	for i := 0; i < blockCount; i++ {
		blockStart[i] = i * blockSize
		tmp := (i + 1) * blockSize
		if tmp > n {
			tmp = n
		}
		blockEnd[i] = tmp
	}
	for i := 0; i < n; i++ {
		belong[i] = i / blockSize
	}

	return struct {
		belong     []int
		blockStart []int
		blockEnd   []int
		blockCount int
	}{belong, blockStart, blockEnd, blockCount}
}
