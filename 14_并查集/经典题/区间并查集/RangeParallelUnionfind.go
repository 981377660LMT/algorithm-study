// https://github.com/yosupo06/library-checker-problems/issues/934
// https://atcoder.jp/contests/yahoo-procon2018-final/submissions/8391439
// http://kmyk.github.io/competitive-programming-library/old/range-union-find-tree.inc.cpp
// https://yosupo.hatenablog.com/entry/2019/11/12/001535

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://atcoder.jp/contests/yahoo-procon2018-final/tasks/yahoo_procon2018_final_d
	// !前缀和后缀的LCP为lens[i]的字符串
	// !LCP => 并查集
	// 给定长为n的数组lens, 问是否存在一个长度为s的字符串,满足:
	// !s[0:i+1] 和 s[n-(i+1):n] 的最长公共前缀为 lens[i] (0<=i<n)
	// n<=3e5 0<=lens[i]<=i+1

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	lens := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &lens[i])
	}

	ufrp := NewUnionFindRangeParallel(n)
	for i := 0; i < n; i++ {
		ufrp.UnionParallelly(0, n-(i+1), lens[i]) // 各个位置的字符相同
	}

	uf := ufrp.Build()
	for i := 0; i < n; i++ {
		if lens[i] == i+1 {
			continue
		}
		if uf.IsConnected(lens[i], n-(i+1)+lens[i]) { // !s[len[i]]!=s[n-(i+1)+len[i]] (因为前后缀LCP只有len[i])
			fmt.Fprintln(out, "No")
			return
		}
	}
	fmt.Fprintln(out, "Yes")
}

// 并行合并的并查集.
type UnionFindRangeParallel struct {
	n    int
	ques [][][2]int
}

func NewUnionFindRangeParallel(n int) *UnionFindRangeParallel {
	return &UnionFindRangeParallel{n: n, ques: make([][][2]int, n+1)}
}

// !并行合并[(a,b),(a+1,b+1),...,(a+len-1,b+len-1)].
func (uf *UnionFindRangeParallel) UnionParallelly(a, b, len int) {
	if len == 0 {
		return
	}
	min_ := min(len, uf.n)
	uf.ques[min_] = append(uf.ques[min_], [2]int{a, b})
}

func (uf *UnionFindRangeParallel) Build() *_unionFindRange {
	res := _newUnionFindRange(uf.n)
	queue, nextQueue := [][2]int{}, [][2]int{}
	for di := uf.n; di >= 1; di-- {
		queue = append(queue, uf.ques[di]...)
		nextQueue = nextQueue[:0]
		for _, p := range queue {
			if res.IsConnected(p[0], p[1]) {
				continue
			}
			res.Union(p[0], p[1], func(big, small int) {})
			nextQueue = append(nextQueue, [2]int{p[0] + 1, p[1] + 1})
		}
		queue, nextQueue = nextQueue, queue
	}
	return res
}

type _unionFindRange struct {
	Part   int
	n      int
	parent []int
	rank   []int
}

func _newUnionFindRange(n int) *_unionFindRange {
	uf := &_unionFindRange{
		Part:   n,
		n:      n,
		parent: make([]int, n),
		rank:   make([]int, n),
	}
	for i := 0; i < n; i++ {
		uf.parent[i] = i
		uf.rank[i] = 1
	}
	return uf
}

func (uf *_unionFindRange) Find(x int) int {
	for x != uf.parent[x] {
		uf.parent[x] = uf.parent[uf.parent[x]]
		x = uf.parent[x]
	}
	return x
}

// Union 后, 大的编号的组会指向小的编号的组.
func (uf *_unionFindRange) Union(x, y int, f func(big, small int)) bool {
	if x < y {
		x, y = y, x
	}
	rootX := uf.Find(x)
	rootY := uf.Find(y)
	if rootX == rootY {
		return false
	}
	uf.parent[rootX] = rootY
	uf.rank[rootY] += uf.rank[rootX]
	uf.Part--
	f(rootY, rootX)
	return true
}

// UnionRange 合并区间 [left, right] 的所有元素, 返回合并次数.
func (uf *_unionFindRange) UnionRange(left, right int, f func(big, small int)) int {
	if left >= right {
		return 0
	}
	leftRoot := uf.Find(left)
	rightRoot := uf.Find(right)
	unionCount := 0
	for rightRoot != leftRoot {
		unionCount++
		uf.Union(rightRoot, rightRoot-1, f)
		rightRoot = uf.Find(rightRoot - 1)
	}
	return unionCount
}

func (uf *_unionFindRange) IsConnected(x, y int) bool {
	return uf.Find(x) == uf.Find(y)
}

func (uf *_unionFindRange) GetSize(x int) int {
	return uf.rank[uf.Find(x)]
}

func (uf *_unionFindRange) GetGroups() map[int][]int {
	group := make(map[int][]int)
	for i := 0; i < uf.n; i++ {
		group[uf.Find(i)] = append(group[uf.Find(i)], i)
	}
	return group
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
