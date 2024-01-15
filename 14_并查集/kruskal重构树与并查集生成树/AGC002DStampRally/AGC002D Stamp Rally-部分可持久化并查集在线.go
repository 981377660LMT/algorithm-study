// AGC002D Stamp Rally-部分可持久化并查集在线
// https://drken1215.hatenablog.com/entry/2018/09/15/193722

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([][2]int, m)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		edges[i] = [2]int{u, v}
	}

	var q int
	fmt.Fscan(in, &q)
	queries := make([][3]int, q)
	for i := 0; i < q; i++ {
		var x, y, z int
		fmt.Fscan(in, &x, &y, &z)
		x--
		y--
		queries[i] = [3]int{x, y, z}
	}

	res := StampRally(n, edges, queries)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// https://atcoder.jp/contests/agc002/tasks/agc002_d
// 一张连通图，q 次询问从两个点 x 和 y 出发，
// 希望经过的点数量等于 z（每个点可以重复经过，但是重复经过只计算一次）
// 求经过的边最大编号最小是多少。
func StampRally(n int, edges [][2]int, queries [][3]int) []int {
	m, q := len(edges), len(queries)
	uf := NewPartiallyPersistentUnionFind(n)
	for i := 0; i < len(edges); i++ {
		u, v := edges[i][0], edges[i][1]
		uf.Union(u, v)
	}

	res := make([]int, q)
	for i := 0; i < q; i++ {
		x, y, z := queries[i][0], queries[i][1], queries[i][2]

		check := func(mid int) bool {
			if uf.IsConnected(mid, x, y) {
				size := uf.GetSize(mid, x)
				return size >= z
			} else {
				size1, size2 := uf.GetSize(mid, x), uf.GetSize(mid, y)
				return size1+size2 >= z
			}
		}

		left, right := 1, m // 二分版本(边的编号1-m)
		for left <= right {
			mid := (left + right) / 2
			if check(mid) {
				right = mid - 1
			} else {
				left = mid + 1
			}
		}
		res[i] = left
	}

	return res
}

// 部分可持久化并查集(初始版本为0).
type PartiallyPersistentUnionFind struct {
	Version int
	data    []int32
	last    []int32
	history [][][2]int32
}

func NewPartiallyPersistentUnionFind(n int) *PartiallyPersistentUnionFind {
	data := make([]int32, n)
	last := make([]int32, n)
	history := make([][][2]int32, n)
	for i := int32(0); i < int32(n); i++ {
		data[i] = -1
		last[i] = 1e9
		history[i] = [][2]int32{{0, -1}}
	}
	return &PartiallyPersistentUnionFind{data: data, last: last, history: history}
}

// 合并x和y所在的集合,返回当前版本号.
func (uf *PartiallyPersistentUnionFind) Union(x, y int) int {
	uf.Version++
	x, y = uf.Find(uf.Version, x), uf.Find(uf.Version, y)
	if x == y {
		return uf.Version
	}
	if uf.data[x] > uf.data[y] {
		x, y = y, x
	}
	uf.data[x] += uf.data[y]
	uf.history[x] = append(uf.history[x], [2]int32{int32(uf.Version), uf.data[x]})
	uf.data[y] = int32(x)
	uf.last[y] = int32(uf.Version)
	return uf.Version
}

func (uf *PartiallyPersistentUnionFind) Find(time, x int) int {
	if time < int(uf.last[x]) {
		return x
	}
	return uf.Find(time, int(uf.data[x]))
}

func (uf *PartiallyPersistentUnionFind) IsConnected(time, x int, y int) bool {
	return uf.Find(time, x) == uf.Find(time, y)
}

func (uf *PartiallyPersistentUnionFind) GetSize(time, x int) int {
	x = uf.Find(time, x)
	tmp := uf.history[x]
	time32 := int32(time)
	pos := sort.Search(len(tmp), func(i int) bool {
		return tmp[i][0] > time32
	})
	return -int(tmp[pos-1][1])
}

func (uf *PartiallyPersistentUnionFind) GetGroups(time int) map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < len(uf.data); i++ {
		root := uf.Find(time, i)
		groups[root] = append(groups[root], i)
	}
	return groups
}
