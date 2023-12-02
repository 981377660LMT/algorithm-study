// https://nyaannyaan.github.io/library/data-structure/union-find-with-potential.hpp
// UnionFindWithDist/UnionFindWithPotential
// 带权并查集(维护到每个组根节点距离的并查集)

// - 注意距离是`有向`的
//   例如维护和距离的并查集时,a->b 的距离是正数,b->a 的距离是负数
// - 如果组内两点距离存在矛盾(沿着不同边走距离不同),那么在组内会出现正环

// API:
//  Union(x,y,dist) : p(x) = p(y) + dist.如果组内两点距离存在矛盾(沿着不同边走距离不同),返回false.
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
	// WeightedUnionFindTrees()
	TreeOfButton()
}

// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=DSL_1_B&lang=jp
func WeightedUnionFindTrees() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	uf := NewUnionFindMapWithDist()
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

// F-ボタンの木
func TreeOfButton() {
	const INF int = 1e18

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	edges := make([][2]int, n-1)
	for i := range edges {
		fmt.Fscan(in, &edges[i][0], &edges[i][1])
		edges[i][0]--
		edges[i][1]--
	}
	starts := make([]int, n)
	for i := range starts {
		fmt.Fscan(in, &starts[i])
	}
	targets := make([]int, n)
	for i := range targets {
		fmt.Fscan(in, &targets[i])
	}

	tree := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	uf := NewUnionFindArrayWithDist(n)
	var dfs func(int, int) int
	dfs = func(cur, pre int) int {
		need := targets[cur] - starts[cur]
		for _, next := range tree[cur] {
			if next == pre {
				continue
			}
			sub := dfs(next, cur)
			uf.Union(cur, next, sub)
			need += sub
		}
		return need
	}

	dfs(0, -1)

	sum_, min_ := 0, INF
	for i := 0; i < n; i++ {
		diff := uf.Dist(i, 0)
		sum_ += diff
		min_ = min(min_, diff)
	}
	fmt.Fprintln(out, sum_-min_*n)
}

// https://leetcode.cn/problems/is-graph-bipartite/
// 785. 判断二分图
// 二分图的充要条件是: 不存在奇环.
func isBipartite(graph [][]int) bool {
	n := len(graph)
	uf := NewUnionFindArrayWithDist(n)
	for cur := range graph {
		for _, next := range graph[cur] {
			root1, root2 := uf.Find(cur), uf.Find(next)
			if root1 == root2 {
				cycleLen := uf.Dist(cur, next) + 1
				if cycleLen&1 == 1 {
					return false
				}
			} else {
				uf.Union(cur, next, 1)
			}
		}
	}
	return true
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type T = int

func e() T        { return 0 }
func op(x, y T) T { return x + y }
func inv(x T) T   { return -x }

// 数组实现的带权并查集(维护到每个组根节点距离的并查集).
// 用于维护环的权值，树上的距离等.
type UnionFindArrayWithDist struct {
	Part      int
	data      []int
	potential []T
}

func NewUnionFindArrayWithDist(n int) *UnionFindArrayWithDist {
	uf := &UnionFindArrayWithDist{
		Part:      n,
		data:      make([]int, n),
		potential: make([]T, n),
	}
	for i := range uf.data {
		uf.data[i] = -1
		uf.potential[i] = e()
	}
	return uf
}

// p[x] = p[y] + dist.
//
//	如果组内两点距离存在矛盾(沿着不同边走距离不同),返回false.
func (uf *UnionFindArrayWithDist) Union(x, y int, dist T) bool {
	dist = op(dist, op(uf.DistToRoot(y), inv(uf.DistToRoot(x))))
	x = uf.Find(x)
	y = uf.Find(y)
	if x == y {
		return dist == e()
	}
	if uf.data[x] < uf.data[y] {
		x, y = y, x
		dist = inv(dist)
	}
	uf.data[y] += uf.data[x]
	uf.data[x] = y
	uf.potential[x] = dist
	uf.Part--
	return true
}

// p[x] = p[y] + dist.
//
//	如果组内两点距离存在矛盾(沿着不同边走距离不同),返回false.
func (uf *UnionFindArrayWithDist) UnionWithCallback(x, y int, dist T, cb func(big, small int)) bool {
	dist = op(dist, op(uf.DistToRoot(y), inv(uf.DistToRoot(x))))
	x = uf.Find(x)
	y = uf.Find(y)
	if x == y {
		return dist == e()
	}
	if uf.data[x] < uf.data[y] {
		x, y = y, x
		dist = inv(dist)
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

func (uf *UnionFindArrayWithDist) Find(x int) int {
	if uf.data[x] < 0 {
		return x
	}
	root := uf.Find(uf.data[x])
	uf.potential[x] = op(uf.potential[x], uf.potential[uf.data[x]])
	uf.data[x] = root
	return root
}

// f[x]-f[find(x)].
//
//	点x到所在组根节点的距离.
func (uf *UnionFindArrayWithDist) DistToRoot(x int) T {
	uf.Find(x)
	return uf.potential[x]
}

// f[x] - f[y].
func (uf *UnionFindArrayWithDist) Dist(x, y int) T {
	return op(uf.DistToRoot(x), inv(uf.DistToRoot(y)))
}

func (uf *UnionFindArrayWithDist) GetSize(x int) int {
	return -uf.data[uf.Find(x)]
}

func (uf *UnionFindArrayWithDist) GetGroups() map[int][]int {
	res := make(map[int][]int)
	for i := range uf.data {
		root := uf.Find(i)
		res[root] = append(res[root], i)
	}
	return res
}

func (uf *UnionFindArrayWithDist) IsConnected(x, y int) bool {
	return uf.Find(x) == uf.Find(y)
}

//
//
//
//

// Map实现的带权并查集(维护到每个组根节点距离的并查集).
type UnionFindMapWithDist struct {
	Part      int
	data      map[int]int
	potential map[int]T
}

func NewUnionFindMapWithDist() *UnionFindMapWithDist {
	return &UnionFindMapWithDist{
		data:      make(map[int]int),
		potential: make(map[int]T),
	}
}

// p[x] = p[y] + dist.
//
//	如果组内两点距离存在矛盾(沿着不同边走距离不同),返回false.
func (uf *UnionFindMapWithDist) Union(x, y int, dist T) bool {
	dist = op(dist, op(uf.DistToRoot(y), inv(uf.DistToRoot(x))))
	x = uf.Find(x)
	y = uf.Find(y)
	if x == y {
		return dist == e()
	}
	if uf.data[x] < uf.data[y] {
		x, y = y, x
		dist = inv(dist)
	}
	uf.data[y] += uf.data[x]
	uf.data[x] = y
	uf.potential[x] = dist
	uf.Part--
	return true
}

// p[x] = p[y] + dist.
//
//	如果组内两点距离存在矛盾(沿着不同边走距离不同),返回false.
func (uf *UnionFindMapWithDist) UnionWithCallback(x, y int, dist T, cb func(big, small int)) bool {
	dist = op(dist, op(uf.DistToRoot(y), inv(uf.DistToRoot(x))))
	x = uf.Find(x)
	y = uf.Find(y)
	if x == y {
		return dist == e()
	}
	if uf.data[x] < uf.data[y] {
		x, y = y, x
		dist = inv(dist)
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

func (uf *UnionFindMapWithDist) Find(x int) int {
	if _, ok := uf.data[x]; !ok {
		uf.Add(x)
		return x
	}
	if uf.data[x] < 0 {
		return x
	}
	root := uf.Find(uf.data[x])
	uf.potential[x] = op(uf.potential[x], uf.potential[uf.data[x]])
	uf.data[x] = root
	return root
}

// f[x]-f[find(x)].
//
//	点x到所在组根节点的距离.
func (uf *UnionFindMapWithDist) DistToRoot(x int) T {
	uf.Find(x)
	return uf.potential[x]
}

// f[x] - f[y].
func (uf *UnionFindMapWithDist) Dist(x, y int) T {
	return op(uf.DistToRoot(x), inv(uf.DistToRoot(y)))
}

func (uf *UnionFindMapWithDist) GetSize(x int) int {
	return -uf.data[uf.Find(x)]
}

func (uf *UnionFindMapWithDist) GetGroups() map[int][]int {
	res := make(map[int][]int)
	for k := range uf.data {
		res[uf.Find(k)] = append(res[uf.Find(k)], k)
	}
	return res
}

func (uf *UnionFindMapWithDist) IsConnected(x, y int) bool {
	return uf.Find(x) == uf.Find(y)
}

func (uf *UnionFindMapWithDist) Add(x int) *UnionFindMapWithDist {
	if _, ok := uf.data[x]; !ok {
		uf.data[x] = -1
		uf.potential[x] = e()
		uf.Part++
	}
	return uf
}

func (uf *UnionFindMapWithDist) Has(x int) bool {
	_, ok := uf.data[x]
	return ok
}
