package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	UnionSets()
	// demo()
}

func demo() {
	uf := NewPartiallyPersistentUnionFind(10)
	v := uf.Union(0, 1)
	fmt.Println(uf.GetSize(v, 0), v)
	fmt.Println(uf.GetSize(0, 0), v)
}

func UnionSets() {
	// https://atcoder.jp/contests/code-thanks-festival-2017-open/tasks/code_thanks_festival_2017_h
	// 给定n个集合,初始时第i个集合只有一个元素i (i=1,2,...,n)
	// 之后进行m次合并操作,每次合并ai和bi所在的集合
	// 如果ai和bi在同一个集合,则无事发生
	// 给定q个询问,问ai和bi是在第几次操作后第一次连通的,如果不连通则输出-1

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	uf := NewPartiallyPersistentUnionFind(n)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a, b = a-1, b-1
		uf.Union(a, b)
	}

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a, b = a-1, b-1
		if !uf.IsConnected(uf.Version, a, b) {
			fmt.Fprintln(out, -1)
			continue
		}

		left, right := 0, m
		for left <= right {
			mid := (left + right) >> 1
			if uf.IsConnected(mid, a, b) {
				right = mid - 1
			} else {
				left = mid + 1
			}
		}
		fmt.Fprintln(out, left)
	}
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
