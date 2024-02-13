// P3603 雪辉-树分块+bitset
// https://www.luogu.com.cn/problem/P3603
// https://oi-wiki.org/ds/tree-decompose/
// 给你一棵 n 个节点且带点权的树，q 个询问，
// !每个询问给你多条链，请你输出这几条链的点的集合并的颜色数和 mex(从0开始)。
// n,q<=1e5, 点权<=3e4
// 方法1：树上撒点分块法(预处理出互为祖孙关系的关键点之间的数的信息)，不采用，难以求解
// !方法2：跳重链分块(路径上的若干条重链的 bitset 并起来)
// 重链上的点的 dfn 序是连续的，序列分块即可

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q, f int
	fmt.Fscan(in, &n, &q, &f)
	maxValue := 0
	values := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &values[i])
		maxValue = max(maxValue, values[i])
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

	lca := NewLCA(tree, []int{0})
	dfnValues := make([]int, n) // !注意要将value按照dfn序存储
	for i := 0; i < n; i++ {
		dfn := lca.dfn[i]
		dfnValues[dfn] = values[i]
	}

	blockSize := 900
	block := UseBlock(n, blockSize)
	belong, blockStart, blockEnd, blockCount := block.belong, block.blockStart, block.blockEnd, block.blockCount
	dp := make([][]Bitset, blockCount) // dp[i][j] 表示第 i 块到第 j 块的 bitset 的并 (0<=i<=j<blockCount)
	for i := 0; i < blockCount; i++ {
		tmp := make([]Bitset, blockCount)
		for j := i; j < blockCount; j++ { // j可以从i开始
			tmp[j] = NewBitset(maxValue + 1)
		}
		dp[i] = tmp
	}
	for i, v := range dfnValues {
		bid := belong[i]
		dp[bid][bid].Set(v)
	}
	for i := 0; i < blockCount; i++ {
		for j := i + 1; j < blockCount; j++ {
			dp[i][j].SetOr(dp[i][j-1], dp[j][j])
		}
	}
	resSet := NewBitset(maxValue + 1)

	getMex := func(b Bitset) int {
		for i := 0; i <= maxValue; i++ {
			if !b.Has(i) {
				return i
			}
		}
		return maxValue + 1
	}

	updateBlock := func(start, end int) {
		bid1, bid2 := belong[start], belong[end-1]
		if bid1 == bid2 {
			for i := start; i < end; i++ {
				resSet.Set(dfnValues[i])
			}
			return
		}
		if bid1+1 <= bid2-1 {
			resSet.IOr(dp[bid1+1][bid2-1])
		}
		for i := start; i < blockEnd[bid1]; i++ {
			resSet.Set(dfnValues[i])
		}
		for i := blockStart[bid2]; i < end; i++ {
			resSet.Set(dfnValues[i])
		}
	}

	queryTree := func(u, v int) int {
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

	queryChains := func(chains [][2]int) (count int, mex int) {
		resSet.ResetAll()
		for _, chain := range chains {
			u, v := chain[0], chain[1]
			queryTree(u, v)
		}
		count = resSet.OnesCount()
		mex = getMex(resSet)
		return
	}

	lastKind, lastMex := 0, 0
	normalize := func(x int) int {
		if f == 0 {
			return x
		}
		return (lastKind + lastMex) ^ x
	}
	for i := 0; i < q; i++ {
		var m int
		fmt.Fscan(in, &m)
		chains := make([][2]int, 0, m)
		for j := 0; j < m; j++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u, v = normalize(u), normalize(v)
			u, v = u-1, v-1
			chains = append(chains, [2]int{u, v})
		}
		lastKind, lastMex = queryChains(chains)
		fmt.Fprintln(out, lastKind, lastMex)
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

// 返回下标 >= p 的第一个 0 的下标，若不存在则返回-1.
func (b Bitset) Next0(p int) int {
	if i := p >> 6; i < len(b) {
		v := b[i]
		if p&63 != 0 {
			v |= ^(^uint(0) << (p & 63))
		}
		if ^v != 0 {
			return i<<6 | bits.TrailingZeros(^v)
		}
		for i++; i < len(b); i++ {
			if ^b[i] != 0 {
				return i<<6 | bits.TrailingZeros(^b[i])
			}
		}
	}
	return -1
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
