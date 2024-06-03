package main

import (
	"fmt"
	"math/bits"
)

// 0
// | \
// 1  2

func main() {
	edegs := [][]int32{{0, 2}, {0, 1}}
	tree := make([][]int32, 3)
	for _, e := range edegs {
		tree[e[0]] = append(tree[e[0]], e[1])
		tree[e[1]] = append(tree[e[1]], e[0])
	}
	fmt.Println(CountConnectedIntervals(tree)) // 5
}

const INF32 int32 = 1e9 + 10

// https://qoj.ac/contest/1277/problem/6674
// Connected Intervals
// 给定一棵树，求二元组(l,r)的个数使得[l,r]内所有的点在一个联通分量内.
// !连通块的个数等于点的个数减去边的个数
// 可以扫描线遍历r，用一个线段树维护(区间最小值 ,区间最小值出现次数)，
// index处的min为1就代表[index,r]这一段在一个联通分量内
func CountConnectedIntervals(tree [][]int32) int {
	n := int32(len(tree))
	seg := NewLazySegTree32(n, func(i int32) E { return E{min: 0, minCount: 1} })
	res := 0
	for r := int32(0); r < n; r++ {
		seg.Update(0, r+1, 1)          // 点
		for _, next := range tree[r] { // 边
			if next < r {
				seg.Update(0, next+1, -1)
			}
		}
		tmp := seg.Query(0, r+1)
		if tmp.min == 1 {
			res += int(tmp.minCount)
		}
	}
	return res
}

// RangeAddRangeMinCount

type E = struct {
	min, minCount int32
}

type Id = int32

func (*LazySegTree32) e() E   { return E{min: INF32, minCount: 0} }
func (*LazySegTree32) id() Id { return 0 }
func (*LazySegTree32) op(left, right E) E {
	if left.min > right.min {
		return right
	}
	if left.min < right.min {
		return left
	}
	return E{min: left.min, minCount: left.minCount + right.minCount}
}
func (*LazySegTree32) mapping(f Id, g E) E {
	g.min += f
	return g
}
func (*LazySegTree32) composition(f, g Id) Id {
	return f + g
}
func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
func max32(a, b int32) int32 {
	if a < b {
		return b
	}
	return a
}

// !template
type LazySegTree32 struct {
	n    int32
	size int32
	log  int32
	data []E
	lazy []Id
}

func NewLazySegTree32(n int32, f func(int32) E) *LazySegTree32 {
	tree := &LazySegTree32{}
	tree.n = n
	tree.log = int32(bits.Len32(uint32(n - 1)))
	tree.size = 1 << tree.log
	tree.data = make([]E, tree.size<<1)
	tree.lazy = make([]Id, tree.size)
	for i := range tree.data {
		tree.data[i] = tree.e()
	}
	for i := range tree.lazy {
		tree.lazy[i] = tree.id()
	}
	for i := int32(0); i < n; i++ {
		tree.data[tree.size+i] = f(i)
	}
	for i := tree.size - 1; i >= 1; i-- {
		tree.pushUp(i)
	}
	return tree
}

// 查询切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *LazySegTree32) Query(left, right int32) E {
	if left < 0 {
		left = 0
	}
	if right > tree.n {
		right = tree.n
	}
	if left >= right {
		return tree.e()
	}
	left += tree.size
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		if ((left >> i) << i) != left {
			tree.pushDown(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushDown((right - 1) >> i)
		}
	}
	sml, smr := tree.e(), tree.e()
	for left < right {
		if left&1 != 0 {
			sml = tree.op(sml, tree.data[left])
			left++
		}
		if right&1 != 0 {
			right--
			smr = tree.op(tree.data[right], smr)
		}
		left >>= 1
		right >>= 1
	}
	return tree.op(sml, smr)
}

// 更新切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *LazySegTree32) Update(left, right int32, f Id) {
	if left < 0 {
		left = 0
	}
	if right > tree.n {
		right = tree.n
	}
	if left >= right {
		return
	}
	left += tree.size
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		if ((left >> i) << i) != left {
			tree.pushDown(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushDown((right - 1) >> i)
		}
	}
	l2, r2 := left, right
	for left < right {
		if left&1 != 0 {
			tree.propagate(left, f)
			left++
		}
		if right&1 != 0 {
			right--
			tree.propagate(right, f)
		}
		left >>= 1
		right >>= 1
	}
	left = l2
	right = r2
	for i := int32(1); i <= tree.log; i++ {
		if ((left >> i) << i) != left {
			tree.pushUp(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushUp((right - 1) >> i)
		}
	}
}

func (tree *LazySegTree32) pushUp(root int32) {
	tree.data[root] = tree.op(tree.data[root<<1], tree.data[root<<1|1])
}
func (tree *LazySegTree32) pushDown(root int32) {
	if tree.lazy[root] != tree.id() {
		tree.propagate(root<<1, tree.lazy[root])
		tree.propagate(root<<1|1, tree.lazy[root])
		tree.lazy[root] = tree.id()
	}
}
func (tree *LazySegTree32) propagate(root int32, f Id) {
	tree.data[root] = tree.mapping(f, tree.data[root])
	// !叶子结点不需要更新lazy
	if root < tree.size {
		tree.lazy[root] = tree.composition(f, tree.lazy[root])
	}
}
