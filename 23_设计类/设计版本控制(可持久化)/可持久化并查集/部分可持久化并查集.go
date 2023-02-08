package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const INF int = 1e18

type PartiallyPersistentUnionFind struct {
	CurVersion int
	Dead       []int
	parent     []int
	height     []int
	size       []int
	history    [][][2]int
}

func NewPartiallyPersistentUnionFind(n int) *PartiallyPersistentUnionFind {
	uf := &PartiallyPersistentUnionFind{
		Dead:    make([]int, n),
		parent:  make([]int, n),
		height:  make([]int, n),
		size:    make([]int, n),
		history: make([][][2]int, n),
	}
	for i := 0; i < n; i++ {
		uf.Dead[i] = INF
		uf.parent[i] = i
		uf.height[i] = 1
		uf.size[i] = 1
		uf.history[i] = [][2]int{{0, 1}}
	}
	return uf
}

// 合并x和y所在的集合,返回当前版本号
func (uf *PartiallyPersistentUnionFind) Union(x int, y int) int {
	px, py := uf.Find(uf.CurVersion, x), uf.Find(uf.CurVersion, y)
	if px == py {
		uf.CurVersion++
		return uf.CurVersion
	}
	if uf.height[py] < uf.height[px] {
		uf.parent[py] = px
		uf.Dead[py] = uf.CurVersion
		uf.size[px] += uf.size[py]
		uf.history[px] = append(uf.history[px], [2]int{uf.CurVersion, uf.size[px]})
	} else {
		uf.parent[px] = py
		uf.Dead[px] = uf.CurVersion
		uf.size[py] += uf.size[px]
		uf.history[py] = append(uf.history[py], [2]int{uf.CurVersion, uf.size[py]})
		cand := uf.height[px] + 1
		if cand > uf.height[py] {
			uf.height[py] = cand
		}
	}

	uf.CurVersion++
	return uf.CurVersion
}

func (uf *PartiallyPersistentUnionFind) Find(time, x int) int {
	for uf.Dead[x] < time {
		x = uf.parent[x]
	}
	return x
}

func (uf *PartiallyPersistentUnionFind) IsConnected(time, x int, y int) bool {
	return uf.Find(time, x) == uf.Find(time, y)
}

func (uf *PartiallyPersistentUnionFind) GetSize(time, x int) int {
	y := uf.Find(time, x)
	pos := sort.Search(len(uf.history[y]), func(i int) bool {
		return uf.history[y][i][0] > time
	})
	return uf.history[y][pos-1][1]
}

func (uf *PartiallyPersistentUnionFind) GetGroups(time int) map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < len(uf.parent); i++ {
		root := uf.Find(time, i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func main() {
	// # https://atcoder.jp/contests/code-thanks-festival-2017-open/tasks/code_thanks_festival_2017_h
	// # 给定n个集合,初始时第i个集合只有一个元素i (i=1,2,...,n)
	// # 之后进行m次合并操作,每次合并ai和bi所在的集合
	// # 如果ai和bi在同一个集合,则无事发生
	// # 给定q个询问,问ai和bi是在第几次操作后第一次连通的,如果不连通则输出-1

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
		if !uf.IsConnected(uf.CurVersion, a, b) {
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
