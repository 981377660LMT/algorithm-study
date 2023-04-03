// 又叫做 64-ary tree / Van Emde Boas Tree
// !时间复杂度:O(log64n)
// https://zhuanlan.zhihu.com/p/107238627
// 使用场景:
// 1. 在存储IP地址的时候， 需要快速查找某个IP地址（2 ^32大小)是否在访问的列表中，
//    或者需要找到比这个IP地址大一点或者小一点的IP作为重新分配的IP。
// 2. 一条路上开了很多商店，用int来表示商店的位置（假设位置为1-256之间的数），
//    不断插入，删除商店，同时需要找到离某个商店最近的商店在哪里。

// !Insert/Erase/Prev/Next/Has/Enumerate

package main

import (
	"fmt"
	"math/bits"
	"strconv"
	"strings"
)

const INF int = 1e18

type Finder struct {
	n, lg int
	seg   [][]int
}

func NewFinder(n int) *Finder {
	res := &Finder{n: n}
	seg := [][]int{}
	n_ := n
	for {
		seg = append(seg, make([]int, (n_+63)/64))
		n_ = (n_ + 63) / 64
		if n_ <= 1 {
			break
		}
	}
	res.seg = seg
	res.lg = len(seg)
	return res
}

func (fs *Finder) Has(i int) bool {
	return (fs.seg[0][i/64]>>(i%64))&1 != 0
}

func (fs *Finder) Insert(i int) {
	for h := 0; h < fs.lg; h++ {
		fs.seg[h][i/64] |= 1 << (i % 64)
		i /= 64
	}
}

func (fs *Finder) Erase(i int) {
	for h := 0; h < fs.lg; h++ {
		fs.seg[h][i/64] &= ^(1 << (i % 64))
		if fs.seg[h][i/64] != 0 {
			break
		}
		i /= 64
	}
}

// 返回大于等于i的最小元素.如果不存在,返回n.
func (fs *Finder) Next(i int) int {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}

	for h := 0; h < fs.lg; h++ {
		if i/64 == len(fs.seg[h]) {
			break
		}
		d := fs.seg[h][i/64] >> (i % 64)
		if d == 0 {
			i = i/64 + 1
			continue
		}
		// find
		i += fs.bsf(d)
		for g := h - 1; g >= 0; g-- {
			i *= 64
			i += fs.bsf(fs.seg[g][i/64])
		}

		return i
	}

	return fs.n
}

// 返回小于等于i的最大元素.如果不存在,返回-1.
func (fs *Finder) Prev(i int) int {
	if i < 0 {
		return -1
	}
	if i >= fs.n {
		i = fs.n - 1
	}

	for h := 0; h < fs.lg; h++ {
		if i == -1 {
			break
		}
		d := fs.seg[h][i/64] << (63 - i%64)
		if d == 0 {
			i = i/64 - 1
			continue
		}
		// find
		i += fs.bsr(d) - 63
		for g := h - 1; g >= 0; g-- {
			i *= 64
			i += fs.bsr(fs.seg[g][i/64])
		}

		return i
	}

	return -1
}

// 遍历[start,end)区间内的元素.
func (fs *Finder) Enumerate(start, end int, f func(i int)) {
	x := start - 1
	for {
		x = fs.Next(x + 1)
		if x >= end {
			break
		}
		f(x)
	}
}

func (fs *Finder) String() string {
	res := []string{}
	for i := 0; i < fs.n; i++ {
		if fs.Has(i) {
			res = append(res, strconv.Itoa(i))
		}
	}
	return fmt.Sprintf("Finder{%v}", strings.Join(res, ", "))
}

func (*Finder) bsr(x int) int {
	return 63 - bits.LeadingZeros(uint(x))
}

func (*Finder) bsf(x int) int {
	return bits.TrailingZeros(uint(x))
}

// 2612. 最少翻转操作数
// https://leetcode.cn/problems/minimum-reverse-operations/
func minReverseOperations(n int, p int, banned []int, k int) []int {
	finder := [2]*Finder{NewFinder(n), NewFinder(n)}
	for i := 0; i < n; i++ {
		finder[i&1].Insert(i)
	}
	for _, i := range banned {
		finder[i&1].Erase(i)
	}

	getRange := func(i int) (int, int) {
		return max(i-k+1, k-i-1), min(i+k-1, 2*n-k-i-1)
	}

	setUsed := func(u int) {
		finder[u&1].Erase(u)
	}

	findUnused := func(u int) int {
		left, right := getRange(u)
		pre := finder[(u+k+1)&1].Prev(right)
		if pre != -1 && left <= pre && pre <= right {
			return pre
		}
		next := finder[(u+k+1)&1].Next(left)
		if next != n && left <= next && next <= right {
			return next
		}
		return -1
	}

	dist := OnlineBfs(n, p, setUsed, findUnused)
	res := make([]int, n)
	for i, d := range dist {
		if d == INF {
			res[i] = -1
		} else {
			res[i] = d
		}
	}
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// 在线bfs.
//   不预先给出图，而是通过两个函数 setUsed 和 findUnused 来在线寻找边.
//   setUsed(u)：将 u 标记为已访问。
//   findUnused(u)：找到和 u 邻接的一个未访问过的点。如果不存在, 返回 `-1`。

func OnlineBfs(
	n int, start int,
	setUsed func(u int), findUnused func(cur int) (next int),
) (dist []int) {
	dist = make([]int, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0
	queue := []int{start}
	setUsed(start)

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for {
			next := findUnused(cur)
			if next == -1 {
				break
			}
			dist[next] = dist[cur] + 1 // weight
			queue = append(queue, next)
			setUsed(next)
		}
	}

	return
}
