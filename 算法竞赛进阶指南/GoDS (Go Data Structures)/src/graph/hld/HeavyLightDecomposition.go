// https://beet-aizu.github.io/library/tree/heavylightdecomposition.cpp
// HL分解将树上的路径分成logn条,分割之后只需要op操作logn条链即可
// 如果原问题可以在数组上O(X)时间解决,那么在树上就可以在O(Xlogn)时间解决
// forEach 处理顶点查询
// forEachEdge 处理边查询
// 如果op运算不满足交换律,需要使用w=lca(u,v)过渡,合成forEach(w,u)和forEach(w,v)的结果

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// 单点加/路径和查询
	// 0 vertex add => 顶点加
	// 1 root1 root2 => 路径和查询
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	hld := NewHeavyLightDecomposition(n)
	values := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &values[i])
	}
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		hld.AddEdge(u, v)
	}
	hld.Build(0)

	bit := NewBITArray(n)
	for i := 0; i < n; i++ {
		id := hld.Get(i)
		bit.Add(id+1, values[i])
	}

	for i := 0; i < q; i++ {
		var op, vertex, add, root1, root2 int
		fmt.Fscan(in, &op)
		if op == 0 {
			fmt.Fscan(in, &vertex, &add)
			bit.Add(hld.Get(vertex)+1, add)
		} else {
			fmt.Fscan(in, &root1, &root2)
			res := 0
			hld.ForEach(root1, root2, func(l, r int) {
				res += bit.QueryRange(l+1, r+1-1)
			})
			fmt.Fprintln(out, res)
		}
	}
}

type HeavyLightDecomposition struct {
	tree                                                  [][]int
	dfn, dfnToNode, top, subSize, depth, parent, heavySon []int
	dfnId                                                 int // !从0开始
}

// !注意：
//  1. dfn 是 0-indexed 的.
//  2. 构建 HLD 需要调用 `Build` 方法.
//  3. 回调函数参数的 dfn区间 是左闭右开的, 即`[left, right)`.
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
		tree:      tree,
		dfn:       dfn,
		dfnToNode: dfnToNode,
		top:       top,
		subSize:   subSize,
		depth:     depth,
		parent:    parent,
		heavySon:  heavySon,
	}
}

// 添加无向边 u-v.
func (hld *HeavyLightDecomposition) AddEdge(u, v int) {
	hld.tree[u] = append(hld.tree[u], v)
	hld.tree[v] = append(hld.tree[v], u)
}

func (hld *HeavyLightDecomposition) Build(root int) {
	hld.build(root, -1, 0)
	hld.markTop(root, root)
}

// 返回树节点 u 对应的 dfs 序号.
//  0 <= u < n, 0 <= id < n.
func (hld *HeavyLightDecomposition) Get(u int) int {
	return hld.dfn[u]
}

// 处理树节点u到v的路径上的所有顶点.
//  回调函数内的参数是左闭右开的 dfn 区间, 即[left, right).
//   0<=left<=right<=n
func (hld *HeavyLightDecomposition) ForEach(u, v int, cb func(left, right int)) {
	for {
		if hld.dfn[u] > hld.dfn[v] {
			u, v = v, u
		}
		cb(max(hld.dfn[hld.top[v]], hld.dfn[u]), hld.dfn[v]+1)
		if hld.top[u] != hld.top[v] {
			v = hld.parent[hld.top[v]]
		} else {
			break
		}
	}
}

// 处理树节点u到v的路径上的所有边.
//  回调函数内的参数是左闭右开的 dfn 区间, 即[left, right)
//   0<=left<=right<=n
func (hld *HeavyLightDecomposition) ForEachEdge(u, v int, cb func(left, right int)) {
	for {
		if hld.dfn[u] > hld.dfn[v] {
			u, v = v, u
		}
		if hld.top[u] != hld.top[v] {
			cb(hld.dfn[hld.top[v]], hld.dfn[v]+1)
			v = hld.parent[hld.top[v]]
		} else {
			if u != v {
				cb(hld.dfn[u]+1, hld.dfn[v]+1)
			}
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
		v = hld.parent[hld.top[v]]
	}
}

func (hld *HeavyLightDecomposition) Dist(u, v int) int {
	return hld.depth[u] + hld.depth[v] - 2*hld.depth[hld.LCA(u, v)]
}

func (hld *HeavyLightDecomposition) build(cur, pre, dep int) int {
	subSize, heavySize, heavySon := 1, 0, -1
	for _, next := range hld.tree[cur] {
		if next != pre {
			nextSize := hld.build(next, cur, dep+1)
			subSize += nextSize
			if nextSize > heavySize {
				heavySize, heavySon = nextSize, next
			}
		}
	}
	hld.depth[cur] = dep
	hld.subSize[cur] = subSize
	hld.heavySon[cur] = heavySon
	hld.parent[cur] = pre
	return subSize
}

func (hld *HeavyLightDecomposition) markTop(cur, top int) {
	hld.top[cur] = top
	hld.dfn[cur] = hld.dfnId
	hld.dfnId++
	hld.dfnToNode[hld.dfn[cur]] = cur
	if hld.heavySon[cur] != -1 {
		hld.markTop(hld.heavySon[cur], top)
		for _, next := range hld.tree[cur] {
			if next != hld.heavySon[cur] && next != hld.parent[cur] {
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

//
//
//

// !Point Add Range Sum, 1-based.
type BITArray struct {
	n    int
	tree []int
}

func NewBITArrayWithIntSlice(nums []int) *BITArray {
	bitArray := &BITArray{
		n:    len(nums),
		tree: make([]int, len(nums)+1),
	}
	bitArray.Build(nums)
	return bitArray
}

func NewBITArray(n int) *BITArray {
	return &BITArray{n: n, tree: make([]int, n+1)}
}

func (b *BITArray) Build(nums []int) {
	for i := 1; i < len(b.tree); i++ {
		b.tree[i] += nums[i-1]
		if j := i + (i & -i); j < len(b.tree) {
			b.tree[j] += b.tree[i]
		}
	}
}

// 位置 index 增加 delta
//  1<=i<=n
func (b *BITArray) Add(index int, delta int) {
	for ; index < len(b.tree); index += index & -index {
		b.tree[index] += delta
	}
}

// 求前缀和
//  1<=i<=n
func (b *BITArray) Query(index int) (res int) {
	if index > b.n {
		index = b.n
	}
	for ; index > 0; index -= index & -index {
		res += b.tree[index]
	}
	return
}

// 1<=left<=right<=n
func (b *BITArray) QueryRange(left, right int) int {
	return b.Query(right) - b.Query(left-1)
}

func (b *BITArray) Len() int {
	return b.n
}

func (b *BITArray) String() string {
	sb := strings.Builder{}
	sb.WriteString("BITArray{")
	for i := 1; i <= b.n; i++ {
		sb.WriteString(fmt.Sprintf("%d", b.QueryRange(i, i)))
		if i != b.n {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("}")
	return sb.String()
}
