package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const INF int = 1e18

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, Q int32
	fmt.Fscan(in, &N, &Q)
	H := make([]int, N)
	for i := int32(0); i < N; i++ {
		fmt.Fscan(in, &H[i])
	}
	queries := make([][2]int32, Q)
	for i := int32(0); i < Q; i++ {
		var l, r int32
		fmt.Fscan(in, &l, &r)
		l--
		queries[i] = [2]int32{l, r}
	}

	res := Buildings(H, queries)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// F - Buildings 2 (单调栈树)
// https://atcoder.jp/contests/abc379/tasks/abc379_f
//
// 给定n个建筑的高度hi，回答q个询问。每个询问给定l,r，问从l,l+1,...,r−1,r建筑往右看，都能看到的建筑的数量。
// 如果建筑i能看到建筑j，则不存在i<k<j，满足hk>hj。
//
// lca 在线解法:
// 1. 引入一个虚拟节点dummy，dummyHeight=INF，位于所有建筑的右侧。
// 2. 使用单调栈确定右侧第一个高度严格大于当前建筑的建筑，构建树。
// 3. 答案为结点lca(l+1,r)的深度。
func Buildings(heights []int, queries [][2]int32) []int32 {
	n := int32(len(heights))
	heights = append(heights, INF)
	stack := []int32{n}
	tree := make([][]int32, n+1)
	for i := int32(n - 1); i >= 0; i-- {
		for len(stack) > 0 && heights[stack[len(stack)-1]] <= heights[i] { // right strict
			stack = stack[:len(stack)-1]
		}
		p := stack[len(stack)-1]
		tree[p] = append(tree[p], i)
		tree[i] = append(tree[i], p)
		stack = append(stack, i)
	}

	depth := GetDepth(tree, n)
	lca := NewFastLca(tree, n)
	query := func(l, r int32) int32 { // 0<=l<=r<=n
		lca_ := lca.Lca(l+1, r)
		return depth[lca_]
	}

	res := make([]int32, len(queries))
	for i, v := range queries {
		res[i] = query(v[0], v[1])
	}
	return res
}

func GetDepth(tree [][]int32, root int32) []int32 {
	depth := make([]int32, len(tree))
	var dfs func(u, p int32)
	dfs = func(u, p int32) {
		for _, v := range tree[u] {
			if v == p {
				continue
			}
			depth[v] = depth[u] + 1
			dfs(v, u)
		}
	}
	dfs(root, -1)
	return depth
}

type FastLca struct {
	time     int32
	preOrder []int32
	i        []int32
	head     []int32
	a        []int32
	parent   []int32
}

func NewFastLca(tree [][]int32, root int32) *FastLca {
	res := &FastLca{
		preOrder: make([]int32, len(tree)),
		i:        make([]int32, len(tree)),
		head:     make([]int32, len(tree)),
		a:        make([]int32, len(tree)),
		parent:   make([]int32, len(tree)),
	}
	res._init(tree, root)
	return res
}

func NewFastLcaWithIsRoot(tree [][]int32, isRoot func(i int32) bool) *FastLca {
	res := &FastLca{
		preOrder: make([]int32, len(tree)),
		i:        make([]int32, len(tree)),
		head:     make([]int32, len(tree)),
		a:        make([]int32, len(tree)),
		parent:   make([]int32, len(tree)),
	}
	res._initWithIsRoot(tree, isRoot)
	return res
}

// floorLog: bits.Len32(uint32(n)) - 1
func (l *FastLca) Lca(x, y int32) int32 {
	var hb int32
	if a, b := l.i[x], l.i[y]; a == b {
		hb = a & -a
	} else {
		hb = 1 << (bits.Len32(uint32(a^b)) - 1)
	}
	tmp := l.a[x] & l.a[y] & -hb
	hz := tmp & -tmp
	ex := l._enterIntoStrip(x, hz)
	ey := l._enterIntoStrip(y, hz)
	if l.preOrder[ex] < l.preOrder[ey] {
		return ex
	} else {
		return ey
	}
}

func (l *FastLca) _init(tree [][]int32, root int32) {
	l.time = 0
	l._dfs1(tree, root, -1)
	l._dfs2(tree, root, -1, 0)
}

func (l *FastLca) _initWithIsRoot(tree [][]int32, isRoot func(i int32) bool) {
	l.time = 0
	for i := int32(0); i < int32(len(tree)); i++ {
		if isRoot(i) {
			l._dfs1(tree, i, -1)
			l._dfs2(tree, i, -1, 0)
		}
	}
}

func (l *FastLca) _dfs1(tree [][]int32, u, p int32) {
	l.parent[u] = p
	l.i[u] = l.time
	l.preOrder[u] = l.time
	l.time++
	for _, v := range tree[u] {
		if v == p {
			continue
		}
		l._dfs1(tree, v, u)
		if a, b := l.i[u], l.i[v]; a&-a < b&-b {
			l.i[u] = b
		}
	}
	l.head[l.i[u]] = u
}

func (l *FastLca) _dfs2(tree [][]int32, u, p, up int32) {
	l.a[u] = up | l.i[u]&-l.i[u]
	for _, v := range tree[u] {
		if v == p {
			continue
		}
		l._dfs2(tree, v, u, l.a[u])
	}
}

func (l *FastLca) _enterIntoStrip(x, hz int32) int32 {
	if a := l.i[x]; a&-a == hz {
		return x
	}
	tmp := l.a[x] & (hz - 1)
	hw := int32(1 << (bits.Len32(uint32(tmp)) - 1))
	return l.parent[l.head[l.i[x]&-hw|hw]]
}
