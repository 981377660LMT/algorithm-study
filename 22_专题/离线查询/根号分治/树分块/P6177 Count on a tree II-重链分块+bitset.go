// https://www.luogu.com.cn/problem/P6177
// 给定一棵树，每个点带点权，每次询问给出 u,v，求 u到 v的简单路径上的不同权值个数。强制在线。
// n<=4e4,q<=1e5
//
// 查询链颜色数，比较好的一种方法是使用 bitset，对值域建 bitset，答案就是 bitset 中 1 的数量。
// 那么现在的问题就是怎么把一条路径上的 bitset 并起来。
// !考虑轻重链剖分，询问时将路径上的若干条重链的 bitset 并起来即可。
// 由于重链上的点的 dfn 序是连续的，序列分块即可。
// 每次询问，跳重链分块计算这条重链的贡献即可。
// O(n^2/w) + qlogn(sqrt(n)+n/w)

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	values := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &values[i])
	}
	tree := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	D := NewDictionary()
	for i, v := range values {
		values[i] = D.Id(v)
	}

	lca := NewLCA(tree, []int{0})
	// top, depth, dfn, parent := lca.top, lca.Depth, lca.dfn, lca.Parent

	blockSize := 1000
	block := UseBlock(n, blockSize)
	belong, blockStart, blockEnd, blockCount := block.belong, block.blockStart, block.blockEnd, block.blockCount
	dp := make([][]Bitset, blockCount) // dp[i][j] 表示第 i 块到第 j 块的 bitset 的并 (0<=i<=j<blockCount)
	for i := 0; i < blockCount; i++ {
		tmp := make([]Bitset, blockCount)
		for j := 0; j < blockCount; j++ {
			tmp[j] = NewBitset(D.Size())
		}
		dp[i] = tmp
	}
	for i, v := range values {
		bid := belong[i]
		dp[bid][bid].Set(v)
	}
	for i := 0; i < blockCount; i++ {
		for j := i + 1; j < blockCount; j++ {
			dp[i][j].SetOr(dp[i][j-1], dp[j][j])
		}
	}
	resSet := NewBitset(D.Size())

	// [left,right]
	updateBlock := func(start, end int) {
		bid1, bid2 := belong[start], belong[end-1]
		if bid1 == bid2 {
			for i := start; i < end; i++ {
				resSet.Set(values[i])
			}
			return
		}
		for i := start; i < blockEnd[bid1]; i++ {
			resSet.Set(values[i])
		}
		for i := blockStart[bid2]; i < end; i++ {
			resSet.Set(values[i])
		}
		if bid1+1 <= bid2-1 {
			resSet.IOr(dp[bid1+1][bid2-1])
		}
	}

	queryTree := func(u, v int) int {
		resSet.ResetAll()
		// for top[u] != top[v] {
		// 	if depth[top[u]] < depth[top[v]] {
		// 		u, v = v, u
		// 	}
		// 	updateBlock(int(dfn[top[u]]), int(dfn[u]))
		// 	u = int(parent[top[u]])
		// }
		// if depth[u] > depth[v] {
		// 	u, v = v, u
		// }
		// updateBlock(int(dfn[u]), int(dfn[v]))
		lca.EnumeratePathDecomposition(u, v, true, func(start, end int) {
			updateBlock(start, end)
		})
		return resSet.OnesCount()
	}

	lastRes := 0
	for i := 0; i < q; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u = u ^ lastRes
		u, v = u-1, v-1
		lastRes = queryTree(u, v)
		fmt.Fprintln(out, lastRes)
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

func (hld *LCAFast) Dist(u, v int) int {
	return int(hld.Depth[u] + hld.Depth[v] - 2*hld.Depth[hld.LCA(u, v)])
}

func (hld *LCAFast) EnumeratePathDecomposition(u, v int, vertex bool, f func(start, end int)) {
	u32, v32 := int32(u), int32(v)
	for {
		if hld.top[u32] == hld.top[v32] {
			break
		}
		if hld.dfn[u32] < hld.dfn[v32] {
			a, b := hld.dfn[hld.top[v32]], hld.dfn[v32]
			if a > b {
				a, b = b, a
			}
			f(int(a), int(b+1))
			v32 = hld.Parent[hld.top[v32]]
		} else {
			a, b := hld.dfn[u32], hld.dfn[hld.top[u32]]
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

	if hld.dfn[u32] < hld.dfn[v32] {
		a, b := hld.dfn[u32]+edgeInt, hld.dfn[v32]
		if a > b {
			a, b = b, a
		}
		f(int(a), int(b+1))
	} else if hld.dfn[v32]+edgeInt <= hld.dfn[u32] {
		a, b := hld.dfn[u32], hld.dfn[v32]+edgeInt
		if a > b {
			a, b = b, a
		}
		f(int(a), int(b+1))
	}
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

type Bitset []uint

func NewBitset(n int) Bitset { return make(Bitset, n>>6+1) } // (n+64-1)>>6, 注意 n=0 的情况，n>>6+1的写法更好

func (b Bitset) Has(p int) bool { return b[p>>6]&(1<<(p&63)) != 0 }
func (b Bitset) Flip(p int)     { b[p>>6] ^= 1 << (p & 63) }
func (b Bitset) Set(p int)      { b[p>>6] |= 1 << (p & 63) }  // 置 1
func (b Bitset) Reset(p int)    { b[p>>6] &^= 1 << (p & 63) } // 置 0
func (b Bitset) ResetAll() {
	for i := range b {
		b[i] = 0
	}
}
func (b Bitset) IOr(other Bitset) {
	for i := range b {
		b[i] |= other[i]
	}
}
func (b Bitset) SetOr(other1, other2 Bitset) {
	for i := range b {
		b[i] = other1[i] | other2[i]
	}
}
func (b Bitset) OnesCount() (c int) {
	for _, v := range b {
		c += bits.OnesCount(v)
	}
	return
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

type V = int
type Dictionary struct {
	_idToValue []V
	_valueToId map[V]int
}

// A dictionary that maps values to unique ids.
func NewDictionary() *Dictionary {
	return &Dictionary{
		_valueToId: map[V]int{},
	}
}
func (d *Dictionary) Id(value V) int {
	res, ok := d._valueToId[value]
	if ok {
		return res
	}
	id := len(d._idToValue)
	d._idToValue = append(d._idToValue, value)
	d._valueToId[value] = id
	return id
}
func (d *Dictionary) Value(id int) V {
	return d._idToValue[id]
}
func (d *Dictionary) Has(value V) bool {
	_, ok := d._valueToId[value]
	return ok
}
func (d *Dictionary) Size() int {
	return len(d._idToValue)
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
