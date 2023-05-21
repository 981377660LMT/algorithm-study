// https://nyaannyaan.github.io/library/data-structure/union-find-with-potential.hpp
// UnionFindWithDist/UnionFindWithPotential
// 带权并查集(维护到每个组根节点距离的并查集)

// - 注意距离是`有向`的
//   例如维护和距离的并查集时,a->b 的距离是正数,b->a 的距离是负数
// - 如果组内两点距离存在矛盾(沿着不同边走距离不同),那么在组内会出现正环

// API:
//  Union(x,y,dist) : p(x) = p(y) + dist.
//  Find(x) : 返回x所在组的根节点.
//  Dist(x,y) : 返回x到y的距离.
//  DistToRoot(x) : 返回x到所在组根节点的距离.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	WeightedUnionFindTrees()
}

// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=DSL_1_B&lang=jp
func WeightedUnionFindTrees() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	uf := NewUnionFindArrayWithDist(n)
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 0 {
			var x, y, w int
			fmt.Fscan(in, &x, &y, &w)
			uf.Union(y, x, w)
		} else {
			var x, y int
			fmt.Fscan(in, &x, &y)
			if !uf.IsConnected(x, y) {
				fmt.Fprintln(out, "?")
			} else {
				fmt.Fprintln(out, uf.Dist(y, x))
			}
		}
	}
}

type T = int

func (uf *UnionFindWithDist) e() T        { return 0 }
func (uf *UnionFindWithDist) op(x, y T) T { return x + y }
func (uf *UnionFindWithDist) inv(x T) T   { return -x }

type UnionFindWithDist struct {
	Part      int
	data      []int
	potential []T
}

func NewUnionFindArrayWithDist(n int) *UnionFindWithDist {
	uf := &UnionFindWithDist{
		Part:      n,
		data:      make([]int, n),
		potential: make([]T, n),
	}
	for i := range uf.data {
		uf.data[i] = -1
		uf.potential[i] = uf.e()
	}
	return uf
}

// p[x] = p[y] + dist.
//  如果组内两点距离存在矛盾(沿着不同边走距离不同),返回false.
func (uf *UnionFindWithDist) Union(x, y int, dist T) bool {
	dist = uf.op(dist, uf.op(uf.DistToRoot(y), uf.inv(uf.DistToRoot(x))))
	x = uf.Find(x)
	y = uf.Find(y)
	if x == y {
		return dist == uf.e()
	}
	if uf.data[x] < uf.data[y] {
		x, y = y, x
		dist = uf.inv(dist)
	}
	uf.data[y] += uf.data[x]
	uf.data[x] = y
	uf.potential[x] = dist
	uf.Part--
	return true
}

// p[x] = p[y] + dist.
//  如果组内两点距离存在矛盾(沿着不同边走距离不同),返回false.
func (uf *UnionFindWithDist) UnionWithCallback(x, y int, dist T, cb func(big, small int)) bool {
	dist = uf.op(dist, uf.op(uf.DistToRoot(y), uf.inv(uf.DistToRoot(x))))
	x = uf.Find(x)
	y = uf.Find(y)
	if x == y {
		return dist == uf.e()
	}
	if uf.data[x] < uf.data[y] {
		x, y = y, x
		dist = uf.inv(dist)
	}
	uf.data[y] += uf.data[x]
	uf.data[x] = y
	uf.potential[x] = dist
	uf.Part--
	if cb != nil {
		cb(y, x)
	}
	return true
}

func (uf *UnionFindWithDist) Find(x int) int {
	if uf.data[x] < 0 {
		return x
	}
	root := uf.Find(uf.data[x])
	uf.potential[x] = uf.op(uf.potential[x], uf.potential[uf.data[x]])
	uf.data[x] = root
	return root
}

// f[x]-f[find(x)].
//  点x到所在组根节点的距离.
func (uf *UnionFindWithDist) DistToRoot(x int) T {
	uf.Find(x)
	return uf.potential[x]
}

// f[x] - f[y].
func (uf *UnionFindWithDist) Dist(x, y int) T {
	return uf.op(uf.DistToRoot(x), uf.inv(uf.DistToRoot(y)))
}

func (uf *UnionFindWithDist) GetSize(x int) int {
	return -uf.data[uf.Find(x)]
}

func (uf *UnionFindWithDist) IsConnected(x, y int) bool {
	return uf.Find(x) == uf.Find(y)
}
