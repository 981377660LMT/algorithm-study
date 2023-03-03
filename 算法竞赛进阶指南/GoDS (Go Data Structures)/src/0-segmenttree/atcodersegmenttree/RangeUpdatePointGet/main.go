// !区间修改单点查询的线段树 Dual-Segment-Tree(双対セグメント木)
// 双対セグメント木は遅延伝搬セグメント木の作用素モノイドのみを取り出したセグメント木を指す.
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=DSL_2_D&lang=ja
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	seg := NewSegmentTree(n)
	for i := 0; i < q; i++ {
		var com int
		fmt.Fscan(in, &com)
		if com == 0 {
			var l, r, x int
			fmt.Fscan(in, &l, &r, &x)
			seg.Update(l, r+1, x)
		} else if com == 1 {
			var k int
			fmt.Fscan(in, &k)
			fmt.Fprintln(out, seg.Get(k))
		}
	}

}

// RangeAssignPointGet
type Id = int

func (*SegmentTree) id() Id                 { return 1<<31 - 1 }
func (*SegmentTree) composition(f, g Id) Id { return g }

type SegmentTree struct {
	size, height int
	lazy         []Id
}

func NewSegmentTree(n int) *SegmentTree {
	res := &SegmentTree{}
	size := 1
	height := 0
	for size < n {
		size <<= 1
		height++
	}
	lazy := make([]Id, 2*size)
	for i := 0; i < 2*size; i++ {
		lazy[i] = res.id()
	}
	res.size = size
	res.height = height
	res.lazy = lazy
	return res
}
func (seg *SegmentTree) Get(index int) Id {
	index += seg.size
	seg.thrust(index)
	return seg.lazy[index]
}
func (seg *SegmentTree) Update(left, right int, value Id) {
	left += seg.size
	right += seg.size - 1
	seg.thrust(left)
	seg.thrust(right)
	l, r := left, right+1
	for l < r {
		if l&1 == 1 {
			seg.lazy[l] = seg.composition(seg.lazy[l], value)
			l++
		}
		if r&1 == 1 {
			r--
			seg.lazy[r] = seg.composition(seg.lazy[r], value)
		}
		l >>= 1
		r >>= 1
	}
}
func (seg *SegmentTree) thrust(k int) {
	for i := seg.height; i > 0; i-- {
		seg.propagate(k >> uint(i))
	}
}
func (seg *SegmentTree) propagate(k int) {
	if seg.lazy[k] != seg.id() {
		seg.lazy[k<<1] = seg.composition(seg.lazy[k<<1], seg.lazy[k])
		seg.lazy[k<<1|1] = seg.composition(seg.lazy[k<<1|1], seg.lazy[k])
		seg.lazy[k] = seg.id()
	}
}
