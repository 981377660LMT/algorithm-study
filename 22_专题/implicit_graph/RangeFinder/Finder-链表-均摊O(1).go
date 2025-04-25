package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	fs := NewFinder(5)
	fmt.Println(fs) // Finder{0, 1, 2, 3, 4}

	fs.Erase(2)
	fmt.Println(fs) // Finder{0, 1, 3, 4}

	fmt.Println("Has(2)?", fs.Has(2))    // false
	fmt.Println("Next(2) =", fs.Next(2)) // 3
	fmt.Println("Prev(2) =", fs.Prev(2)) // 1

	debug := func() {
		for i := 0; i < fs.n; i++ {
			fmt.Printf("Prev(%d) = %d\n", i, fs.Prev(i))
		}
		for i := 0; i < fs.n; i++ {
			fmt.Printf("Next(%d) = %d\n", i, fs.Next(i))
		}
	}
	_ = debug

	debug()
}

type Finder struct {
	n          int
	exist      []bool
	prev, next []int
}

// 建立一个包含0到n-1的集合.
func NewFinder(n int) *Finder {
	res := &Finder{
		n:     n,
		exist: make([]bool, n),
		prev:  make([]int, n),
		next:  make([]int, n),
	}
	for i := 0; i < n; i++ {
		res.exist[i] = true
		res.prev[i] = i - 1
		res.next[i] = i + 1
	}
	return res
}

// 0<=i<n.
func (fs *Finder) Has(i int) bool {
	return i >= 0 && i < fs.n && fs.exist[i]
}

// 0<=i<n.
func (fs *Finder) Erase(i int) bool {
	if !fs.Has(i) {
		return false
	}
	fs.exist[i] = false
	return true
}

// 返回小于等于i的最大元素.如果不存在,返回-1.
// 0<=i<n.
func (fs *Finder) Prev(i int) int {
	if i < 0 {
		return -1
	}
	if i >= fs.n {
		i = fs.n - 1
	}
	if fs.exist[i] {
		return i
	}

	realPrev := fs.prev[i]
	for realPrev >= 0 && !fs.exist[realPrev] {
		realPrev = fs.prev[realPrev]
	}
	cur := i
	for cur >= 0 && cur != realPrev {
		tmp := fs.prev[cur]
		fs.prev[cur] = realPrev
		cur = tmp
	}
	return realPrev
}

// 返回大于等于i的最小元素.如果不存在,返回n.
// 0<=i<n.
func (fs *Finder) Next(i int) int {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}
	if fs.exist[i] {
		return i
	}

	realNext := fs.next[i]
	for realNext < fs.n && !fs.exist[realNext] {
		realNext = fs.next[realNext]
	}
	cur := i
	for cur < fs.n && cur != realNext {
		tmp := fs.next[cur]
		fs.next[cur] = realNext
		cur = tmp
	}
	return realNext
}

// 遍历[start,end)区间内的元素.
func (fs *Finder) Enumerate(start, end int, f func(i int)) {
	if start < 0 {
		start = 0
	}
	if end > fs.n {
		end = fs.n
	}
	if start >= end {
		return
	}

	for x := fs.Next(start); x < end; x = fs.Next(x + 1) {
		f(x)
	}
}

func (fs *Finder) String() string {
	var res []string
	fs.Enumerate(0, fs.n, func(i int) {
		res = append(res, strconv.Itoa(i))
	})
	return fmt.Sprintf("Finder{%v}", strings.Join(res, ", "))
}

const INF int = 1e18

// 2612. 最少翻转操作数
// https://leetcode.cn/problems/minimum-reverse-operations/
func minReverseOperations(n int, p int, banned []int, k int) []int {
	finder := [2]*Finder{NewFinder(n), NewFinder(n)}

	for i := 0; i < n; i++ {
		finder[(i&1)^1].Erase(i)
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
		next := finder[(u+k+1)&1].Next(left)
		if left <= next && next <= right {
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
