// 线性Lca/FastLca
// https://codeforces.com/blog/lrvideckis
// https://codeforces.com/blog/entry/125371
// lhttps://blog.csdn.net/kksleric/article/details/7836649
// https://zhuanlan.zhihu.com/p/79423299?utm\_psn=1859183512999059456
// Schieber-Vishkin O(n),O(1) 求lca，

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	yosupo()
}

// https://judge.yosupo.jp/problem/lca
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	tree := make([][]int32, n)
	for i := 1; i < n; i++ {
		var parent int32
		fmt.Fscan(in, &parent)
		tree[parent] = append(tree[parent], int32(i))
	}
	bl := NewLcaOnTreeBySchieberVishkin(tree, 0)
	for i := 0; i < q; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		fmt.Fprintln(out, bl.Lca(u, v))
	}
}

// O(n)时空间预处理，O(1)查询LCA。
type LcaOnTreeBySchieberVishkin struct {
	time     int32
	preOrder []int32
	i        []int32
	head     []int32
	a        []int32
	parent   []int32
}

func NewLcaOnTreeBySchieberVishkin(tree [][]int32, root int32) *LcaOnTreeBySchieberVishkin {
	res := &LcaOnTreeBySchieberVishkin{
		preOrder: make([]int32, len(tree)),
		i:        make([]int32, len(tree)),
		head:     make([]int32, len(tree)),
		a:        make([]int32, len(tree)),
		parent:   make([]int32, len(tree)),
	}
	res._init(tree, root)
	return res
}

func NewLcaOnTreeBySchieberVishkinWithIsRoot(tree [][]int32, isRoot func(i int32) bool) *LcaOnTreeBySchieberVishkin {
	res := &LcaOnTreeBySchieberVishkin{
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
func (l *LcaOnTreeBySchieberVishkin) Lca(x, y int32) int32 {
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

func (l *LcaOnTreeBySchieberVishkin) _init(tree [][]int32, root int32) {
	l.time = 0
	l._dfs1(tree, root, -1)
	l._dfs2(tree, root, -1, 0)
}

func (l *LcaOnTreeBySchieberVishkin) _initWithIsRoot(tree [][]int32, isRoot func(i int32) bool) {
	l.time = 0
	for i := int32(0); i < int32(len(tree)); i++ {
		if isRoot(i) {
			l._dfs1(tree, i, -1)
			l._dfs2(tree, i, -1, 0)
		}
	}
}

func (l *LcaOnTreeBySchieberVishkin) _dfs1(tree [][]int32, u, p int32) {
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

func (l *LcaOnTreeBySchieberVishkin) _dfs2(tree [][]int32, u, p, up int32) {
	l.a[u] = up | l.i[u]&-l.i[u]
	for _, v := range tree[u] {
		if v == p {
			continue
		}
		l._dfs2(tree, v, u, l.a[u])
	}
}

func (l *LcaOnTreeBySchieberVishkin) _enterIntoStrip(x, hz int32) int32 {
	if a := l.i[x]; a&-a == hz {
		return x
	}
	tmp := l.a[x] & (hz - 1)
	hw := int32(1 << (bits.Len32(uint32(tmp)) - 1))
	return l.parent[l.head[l.i[x]&-hw|hw]]
}
