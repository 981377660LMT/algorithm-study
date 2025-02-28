package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	p2894()
}

// P2894 [USACO08FEB] Hotel G
// https://www.luogu.com.cn/problem/P2894
// 第一行输入 n,m，n 代表有 n 个房间 (1≤n≤50,000)，编号为 1∼n，开始都为空房，m 表示以下有 m 行操作 (1≤m<50,000)，以下每行先输入一个数 i ，表示一种操作：
// 若 i 为 1，表示查询房间，再输入一个数 x，表示在 1,2,...,n 房间中找到长度为 x 的连续空房，输出连续 x 个房间中左端的房间号，尽量让这个房间号最小，若找不到长度为 x 的连续空房，输出 0。若找得到，在这 x 个空房间中住上人。
// 若 i 为 2，表示退房，再输入两个数 x,y 代表房间号 x∼x+y−1 退房，即让房间为空。
func p2894() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int32
	fmt.Fscan(in, &n, &m)

	seg := NewSegTreeLongest0(n)
	for i := int32(0); i < m; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 1 {
			var x int32
			fmt.Fscan(in, &x)
			pos := seg.FindFirst(0, x)
			if pos == -1 {
				fmt.Fprintln(out, 0)
			} else {
				fmt.Fprintln(out, pos+1)
				seg.Update(pos, pos+x, 1)
			}
		} else {
			var x, y int32
			fmt.Fscan(in, &x, &y)
			seg.Update(x-1, x+y-1, 0)
		}
	}
}

type node []struct {
	l, r             int32
	max0, pre0, suf0 int32
	lazy             int32
}

func (t node) propagate(i int32, lazy int32) {
	o := &t[i]
	size := int32(0)
	if lazy == 0 {
		size = o.r - o.l + 1
	}
	o.max0 = size
	o.pre0 = size
	o.suf0 = size
	o.lazy = lazy
}

func (t node) pushDown(i int32) {
	if t[i].lazy != -1 {
		v := t[i].lazy
		t.propagate(i<<1, v)
		t.propagate(i<<1|1, v)
		t[i].lazy = -1
	}
}

func (t node) pushUp(i int32) {
	o := &t[i]
	lo, ro := t[i<<1], t[i<<1|1]
	o.pre0 = lo.pre0
	if lo.pre0 == lo.r-lo.l+1 {
		o.pre0 += ro.pre0
	}
	o.suf0 = ro.suf0
	if ro.suf0 == ro.r-ro.l+1 {
		o.suf0 += lo.suf0
	}
	o.max0 = max(lo.max0, ro.max0, lo.suf0+ro.pre0)
}

// !维护区间最大连续0的懒标记线段树.
func NewSegTreeLongest0(n int32) node {
	t := make(node, 2<<bits.Len32(uint32(n-1)))
	t.build(1, 0, n-1)
	return t
}

func (t node) Update(start, end int32, v int32) {
	t.update(1, start, end-1, v)
}

func (t node) FindFirst(start, size int32) int32 {
	return t.findFirst(start+1, size)
}

func (t node) findFirst(i, size int32) int32 {
	o := &t[i]
	if o.max0 < size {
		return -1
	}
	if o.l == o.r {
		return o.l
	}

	t.pushDown(i)

	idx := t.findFirst(i<<1, size)
	if idx == -1 {
		if t[i<<1].suf0+t[i<<1|1].pre0 >= size {
			m := (o.l + o.r) >> 1
			return m - t[i<<1].suf0 + 1
		}
		idx = t.findFirst(i<<1|1, size)
	}
	return idx
}

func (t node) update(i, L, R int32, v int32) {
	if L <= t[i].l && t[i].r <= R {
		t.propagate(i, v)
		return
	}
	t.pushDown(i)
	m := (t[i].l + t[i].r) >> 1
	if L <= m {
		t.update(i<<1, L, R, v)
	}
	if m < R {
		t.update(i<<1|1, L, R, v)
	}
	t.pushUp(i)
}

func (t node) build(i, l, r int32) {
	o := &t[i]
	o.l, o.r = l, r
	t.propagate(i, 0)
	o.lazy = -1
	if l == r {
		return
	}
	m := (l + r) >> 1
	t.build(i<<1, l, m)
	t.build(i<<1|1, m+1, r)
}

type interval struct {
	start, end int32
}

type Allocator struct {
	seg     node
	buckets map[int][]interval
}

func Constructor(n int) Allocator {
	return Allocator{
		seg:     NewSegTreeLongest0(int32(n)),
		buckets: make(map[int][]interval),
	}
}

func (t Allocator) Allocate(size int, mID int) int {
	pos := t.seg.FindFirst(0, int32(size))
	if pos == -1 {
		return -1
	}
	start, end := pos, pos+int32(size)
	t.buckets[mID] = append(t.buckets[mID], interval{start, end})
	t.seg.Update(start, end, 1)
	return int(pos)
}

func (t Allocator) FreeMemory(mID int) (res int) {
	for _, se := range t.buckets[mID] {
		t.seg.Update(se.start, se.end, 0)
		res += int(se.end - se.start)
	}
	delete(t.buckets, mID)
	return res
}
