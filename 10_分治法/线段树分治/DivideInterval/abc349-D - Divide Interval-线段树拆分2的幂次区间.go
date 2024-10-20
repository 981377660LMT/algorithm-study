package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	demo()
	// abc349_d()
}

func demo() {
	count := 0
	D := NewDivideInterval(10)
	D.EnumerateSegment(3, 7, func(segmentId int) {
		count++
	}, false)
	fmt.Println(D.DivideCount(3, 7), count)
}

// D - Divide Interval (abc349 D)
// https://atcoder.jp/contests/abc349/tasks/abc349_d
// 给定[l,l+1,...,r−1,r)序列，拆分成最少的序列个数，使得每个序列形如[2^i*j,2^i*(j+1))。
// !即：拆分成若干个长度为2的幂的区间.
func abc349_d() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var L, R int
	fmt.Fscan(in, &L, &R)

	D := NewDivideInterval(1 << 62)
	var segs [][2]int
	D.EnumerateSegment(L, R, func(segmentId int) {
		start, end := D.IdToSegment(segmentId)
		segs = append(segs, [2]int{start, end})
	}, true)

	fmt.Fprintln(out, len(segs))
	for _, seg := range segs {
		fmt.Fprintln(out, seg[0], seg[1])
	}
}

type DivideInterval struct {
	Offset int // 线段树中一共offset+n个节点,offset+i对应原来的第i个节点.
	n      int
	log    int
}

// 线段树分割区间.
// 将长度为n的序列搬到长度为offset+n的线段树上, 以实现快速的区间操作.
func NewDivideInterval(n int) *DivideInterval {
	offset := 1
	log := 1
	for offset < n {
		offset <<= 1
		log++
	}
	return &DivideInterval{Offset: offset, n: n, log: log}
}

// 获取原下标为i的元素在树中的(叶子)编号.
func (d *DivideInterval) Id(rawIndex int) int {
	return rawIndex + d.Offset
}

// O(logn) 顺序遍历`[start,end)`区间对应的线段树节点.
// sorted表示是否按照节点编号从小到大的顺序遍历.
func (d *DivideInterval) EnumerateSegment(start, end int, f func(segmentId int), sorted bool) {
	if start < 0 {
		start = 0
	}
	if end > d.n {
		end = d.n
	}
	if start >= end {
		return
	}

	if sorted {
		for _, i := range d.getSegmentIds(start, end) {
			f(i)
		}
	} else {
		for start, end = start+d.Offset, end+d.Offset; start < end; start, end = start>>1, end>>1 {
			if start&1 == 1 {
				f(start)
				start++
			}
			if end&1 == 1 {
				end--
				f(end)
			}
		}
	}
}

func (d *DivideInterval) EnumeratePoint(index int, f func(segmentId int)) {
	if index < 0 || index >= d.n {
		return
	}
	index += d.Offset
	for index > 0 {
		f(index)
		index >>= 1
	}
}

// 线段树结点对应的区间.
func (d *DivideInterval) IdToSegment(id int) (start, end int) {
	if d.IsLeaf(id) {
		id -= d.Offset
		return id, id + 1
	}
	len := bits.Len64(uint64(id))
	start = id<<(d.log-len) - d.Offset
	end = start + (1 << (d.log - len))
	return
}

// O(n) 从根向叶子方向push.
func (d *DivideInterval) PushDown(f func(parent, child int)) {
	for p := 1; p < d.Offset; p++ {
		f(p, p<<1)
		f(p, p<<1|1)
	}
}

// O(n) 从叶子向根方向update.
func (d *DivideInterval) PushUp(f func(parent, child1, child2 int)) {
	for p := d.Offset - 1; p > 0; p-- {
		f(p, p<<1, p<<1|1)
	}
}

// 线段树的节点个数.
func (d *DivideInterval) Size() int {
	return d.Offset + d.n
}

func (d *DivideInterval) IsLeaf(segmentId int) bool {
	return segmentId >= d.Offset
}

func (d *DivideInterval) Depth(u int) int {
	if u == 0 {
		return 0
	}
	return bits.LeadingZeros64(uint64(u)) - 1
}

// 线段树(完全二叉树)中两个节点的最近公共祖先(两个二进制数字的最长公共前缀).
func (d *DivideInterval) Lca(u, v int) int {
	if u == v {
		return u
	}
	if u > v {
		u, v = v, u
	}
	depth1 := d.Depth(u)
	depth2 := d.Depth(v)
	diff := u ^ (v >> (depth2 - depth1))
	if diff == 0 {
		return u
	}
	len := bits.Len64(uint64(diff))
	return u >> len
}

// 区间拆分成的结点个数.
func (d *DivideInterval) DivideCount(start, end int) int {
	return bits.Len64(uint64(end - start))
}

func (d *DivideInterval) getSegmentIds(start, end int) []int {
	if !(0 <= start && start <= end && end <= d.n) {
		return nil
	}
	var leftRes, rightRes []int
	for start, end = start+d.Offset, end+d.Offset; start < end; start, end = start>>1, end>>1 {
		if start&1 == 1 {
			leftRes = append(leftRes, start)
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = append(rightRes, end)
		}
	}
	for i := len(rightRes) - 1; i >= 0; i-- {
		leftRes = append(leftRes, rightRes[i])
	}
	return leftRes
}
