// 100475. 连接两棵树后最大目标节点数目 I
// https://leetcode.cn/problems/maximize-the-number-of-target-nodes-after-connecting-trees-i/description/
//
// !如果节点 u 和节点 v 之间路径的边数小于等于 k ，那么我们称节点 u 是节点 v 的 目标节点 。
// 注意 ，一个节点一定是它自己的 目标节点 。
//
// !请你返回一个长度为 n 的整数数组 answer ，answer[i] 表示将第一棵树中的一个节点与第二棵树中的一个节点连接一条边后，
// !第一棵树中节点 i 的 目标节点 数目的 最大值 。
//
// 注意 ，每个查询相互独立。意味着进行下一次查询之前，你需要先把刚添加的边给删掉。
//
// 2 <= n,m <= 1e5
// 0<=k<=1e5
//
// !等高线汇集算法求到每个顶点，距离<=v的点的个数
// !相当于 VertexAddRangeContourSum 问题.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	abc()
}

// https://atcoder.jp/contests/yahoo-procon2018-final/tasks/yahoo_procon2018_final_c
// !距离树中点node距离恰好为k的点的个数.
func abc() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, Q int32
	fmt.Fscan(in, &N, &Q)

	tree := make([][]int32, N)
	for i := int32(0); i < N-1; i++ {
		var a, b int32
		fmt.Fscan(in, &a, &b)
		a, b = a-1, b-1
		tree[a] = append(tree[a], b)
		tree[b] = append(tree[b], a)
	}

	C := NewContourQueryRange(N, tree)
	presum := make([]int32, C.Size()+1)
	for i := int32(0); i < N; i++ {
		C.EnumeratePoint(i, func(v int32) {
			presum[v+1]++
		})
	}
	for i := int32(1); i < int32(len(presum)); i++ {
		presum[i] += presum[i-1]
	}

	// !距离树中点node距离<=threshold的点的个数.
	query := func(node int32, threshold int32) int32 {
		if threshold < 0 {
			return 0
		}
		res := int32(1)
		C.EnumerateRange(node, 0, threshold+1, func(lo, hi int32) {
			res += presum[hi] - presum[lo]
		})
		return res
	}

	for i := int32(0); i < Q; i++ {
		var node, threshold int32
		fmt.Fscan(in, &node, &threshold)
		node--
		fmt.Println(query(node, threshold) - query(node, threshold-1))
	}
}

// 100475. 连接两棵树后最大目标节点数目 I
func maxTargetNodes(edges1 [][]int, edges2 [][]int, k int) []int {
	n, m := int32(len(edges1)+1), int32(len(edges2)+1)
	tree1, tree2 := make([][]int32, n), make([][]int32, m)
	for _, e := range edges1 {
		u, v := int32(e[0]), int32(e[1])
		tree1[u] = append(tree1[u], v)
		tree1[v] = append(tree1[v], u)
	}
	for _, e := range edges2 {
		u, v := int32(e[0]), int32(e[1])
		tree2[u] = append(tree2[u], v)
		tree2[v] = append(tree2[v], u)
	}

	C1, C2 := NewContourQueryRange(n, tree1), NewContourQueryRange(m, tree2)
	presum1, presum2 := make([]int32, C1.Size()+1), make([]int32, C2.Size()+1)
	for i := int32(0); i < n; i++ {
		C1.EnumeratePoint(i, func(v int32) {
			presum1[v+1]++
		})
	}
	for i := int32(0); i < m; i++ {
		C2.EnumeratePoint(i, func(v int32) {
			presum2[v+1]++
		})
	}
	for i := int32(1); i < int32(len(presum1)); i++ {
		presum1[i] += presum1[i-1]
	}
	for i := int32(1); i < int32(len(presum2)); i++ {
		presum2[i] += presum2[i-1]
	}

	// 距离树中点node距离<=threshold的点的个数.
	rangeCount := func(contour *ContourQueryRange, presum []int32, node int32, threshold int32) int32 {
		res := int32(1)
		contour.EnumerateRange(node, 0, threshold+1, func(lo, hi int32) {
			res += presum[hi] - presum[lo]
		})
		return res
	}

	maxCount2 := int32(0)
	if k > 0 {
		for i := int32(0); i < m; i++ {
			count := rangeCount(C2, presum2, i, int32(k-1))
			maxCount2 = max32(maxCount2, count)
		}
	}

	res := make([]int, n)
	for i := int32(0); i < n; i++ {
		count := rangeCount(C1, presum1, i, int32(k))
		res[i] = int(count + maxCount2)
	}
	return res
}

// !注意不包含距离0.
type ContourQueryRange struct {
	_n          int32
	_v          []int32
	_comp       []int32
	_dep        []int32
	_infoIdx    []int32
	_infoIndptr []int32
	_compRange  []int32
}

func NewContourQueryRange(n int32, tree [][]int32) *ContourQueryRange {
	p := int32(0)
	compRange := []int32{0}
	V := []int32{}
	comp := []int32{}
	dep := []int32{}
	f := func(par []int32, vs []int32, color []int8) {
		n := int32(len(par))
		dist := make([]int32, n)

		for v := int32(1); v < n; v++ {
			dist[v] = dist[par[v]] + 1
		}

		for c1 := int8(0); c1 < 2; c1++ {
			var A, B []int32
			for v := int32(0); v < n; v++ {
				if color[v] == c1 {
					A = append(A, v)
				}
				if color[v] > c1 {
					B = append(B, v)
				}
			}
			if len(A) == 0 || len(B) == 0 {
				return
			}
			mxA := int32(0)
			mxB := int32(0)
			for _, v := range A {
				V = append(V, vs[v])
				comp = append(comp, p)
				dep = append(dep, dist[v])
				mxA = max32(mxA, dist[v])
			}
			compRange = append(compRange, compRange[len(compRange)-1]+mxA+1)
			p++
			for _, v := range B {
				V = append(V, vs[v])
				comp = append(comp, p)
				dep = append(dep, dist[v])
				mxB = max32(mxB, dist[v])
			}
			compRange = append(compRange, compRange[len(compRange)-1]+mxB+1)
			p++
		}
	}

	CentroidDecomposition2(n, tree, f)
	infoIndptr := make([]int32, n+1)
	for _, v := range V {
		infoIndptr[v+1]++
	}
	for v := int32(0); v < n; v++ {
		infoIndptr[v+1] += infoIndptr[v]
	}
	counter := append([]int32{}, infoIndptr...)
	infoIdx := make([]int32, infoIndptr[len(infoIndptr)-1])
	for i := int32(0); i < int32(len(V)); i++ {
		infoIdx[counter[V[i]]] = i
		counter[V[i]]++
	}
	return &ContourQueryRange{
		_n:          n,
		_v:          V,
		_comp:       comp,
		_dep:        dep,
		_infoIdx:    infoIdx,
		_infoIndptr: infoIndptr,
		_compRange:  compRange,
	}
}

func (cqr *ContourQueryRange) Size() int32 {
	return cqr._compRange[len(cqr._compRange)-1]
}

func (cqr *ContourQueryRange) EnumerateRange(node int32, start int32, end int32, f func(int32, int32)) {
	for k := cqr._infoIndptr[node]; k < cqr._infoIndptr[node+1]; k++ {
		idx := cqr._infoIdx[k]
		p := cqr._comp[idx] ^ 1
		lo := start - cqr._dep[idx]
		hi := end - cqr._dep[idx]
		L := cqr._compRange[p]
		R := cqr._compRange[p+1]
		n := R - L
		lo = max32(lo, 0)
		hi = min32(hi, n)
		if lo < hi {
			f(L+lo, L+hi)
		}
	}
}

func (cqr *ContourQueryRange) EnumeratePoint(v int32, f func(int32)) {
	for k := cqr._infoIndptr[v]; k < cqr._infoIndptr[v+1]; k++ {
		idx := cqr._infoIdx[k]
		p := cqr._comp[idx]
		f(cqr._compRange[p] + cqr._dep[idx])
	}
}

func CentroidDecomposition2(
	n int32, tree [][]int32,
	f func(parent, vertex []int32, color []int8),
) {
	if n == 1 {
		return
	}
	queue := make([]int32, n)
	parent := make([]int32, n)
	for i := range parent {
		parent[i] = -1
	}
	l := int32(0)
	r := int32(0)
	queue[r] = int32(0)
	r++
	for l < r {
		v := queue[l]
		l++
		for _, next := range tree[v] {
			if next != parent[v] {
				queue[r] = next
				parent[next] = v
				r++
			}
		}
	}
	if r != n {
		panic("r should be equal to n")
	}
	newIdx := make([]int32, n)
	for i := int32(0); i < n; i++ {
		newIdx[queue[i]] = i
	}
	tmp := make([]int32, n)
	for i := int32(0); i < n; i++ {
		tmp[i] = -1
	}
	for i := int32(1); i < n; i++ {
		j := parent[i]
		tmp[newIdx[i]] = newIdx[j]
	}
	parent = tmp

	real := make([]bool, n)
	for i := range real {
		real[i] = true
	}
	centroidDecomposition2Dfs(parent, queue, real, f)
}

func centroidDecomposition2Dfs(
	parent []int32, vs []int32, real []bool,
	f func(parent, vertex []int32, color []int8),
) {
	n := int32(len(parent))
	if n <= 1 {
		panic("N should be greater than or equal to 2")
	}
	if n == 2 {
		if real[0] && real[1] {
			color := []int8{0, 1}
			f(parent, vs, color)
		}
		return
	}
	c := int32(-1)
	sz := make([]int32, n)
	for i := range sz {
		sz[i] = 1
	}
	for i := n - 1; i >= 0; i-- {
		if sz[i] >= (n+1)>>1 {
			c = i
			break
		}
		sz[parent[i]] += sz[i]
	}
	color := make([]int8, n)
	ord := make([]int32, n)
	for i := range color {
		color[i] = -1
		ord[i] = -1
	}
	take := int32(0)
	ord[c] = 0
	p := int32(1)
	for v := int32(1); v < n; v++ {
		if parent[v] == c && take+sz[v] <= (n-1)>>1 {
			color[v] = 0
			ord[v] = p
			p++
			take += sz[v]
		}
	}
	for i := int32(1); i < n; i++ {
		if color[parent[i]] == 0 {
			color[i] = 0
			ord[i] = p
			p++
		}
	}
	n0 := p - 1
	for a := parent[c]; a != -1; a = parent[a] {
		color[a] = 1
		ord[a] = p
		p++
	}
	for i := int32(0); i < n; i++ {
		if i != c && color[i] == -1 {
			color[i] = 1
			ord[i] = p
			p++
		}
	}
	if p != n {
		panic("p should be equal to N")
	}
	n1 := n - 1 - n0
	par0 := make([]int32, n0+1)
	for i := range par0 {
		par0[i] = -1
	}
	par1 := make([]int32, n1+1)
	for i := range par1 {
		par1[i] = -1
	}
	par2 := make([]int32, n)
	for i := range par2 {
		par2[i] = -1
	}
	V0 := make([]int32, n0+1)
	V1 := make([]int32, n1+1)
	V2 := make([]int32, n)
	rea0 := make([]bool, n0+1)
	rea1 := make([]bool, n1+1)
	rea2 := make([]bool, n)
	for v := int32(0); v < n; v++ {
		i := ord[v]
		V2[i] = vs[v]
		rea2[i] = real[v]
		if color[v] != 1 {
			V0[i] = vs[v]
			rea0[i] = real[v]
		}
		if color[v] != 0 {
			V1[max32(i-n0, 0)] = vs[v]
			rea1[max32(i-n0, 0)] = real[v]
		}
	}
	for v := int32(1); v < n; v++ {
		a := ord[v]
		b := ord[parent[v]]
		if a > b {
			a, b = b, a
		}
		par2[b] = a
		if color[v] != 1 && color[parent[v]] != 1 {
			par0[b] = a
		}
		if color[v] != 0 && color[parent[v]] != 0 {
			par1[max32(b-n0, 0)] = max32(a-n0, 0)
		}
	}
	color = make([]int8, n)
	for i := int32(0); i < n; i++ {
		color[i] = -1
	}
	for i := int32(1); i < n; i++ {
		if rea2[i] {
			if i <= n0 {
				color[i] = 0
			} else {
				color[i] = 1
			}
		}
	}
	f(par2, V2, color)
	centroidDecomposition2Dfs(par0, V0, rea0, f)
	centroidDecomposition2Dfs(par1, V1, rea1, f)
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

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
