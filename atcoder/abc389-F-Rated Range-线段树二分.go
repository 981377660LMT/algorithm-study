// F - Rated Range
// https://atcoder.jp/contests/abc389/tasks/abc389_f
// !总共n场比赛，第i场比赛如果评分在[Li,Ri]，则评分加1。现有Q次查询，每次查询给出始分数X，求n场比赛后的总分。
// n,q<=2e5, Li,Ri,X<=5e5.
//
// !注意到 f(x) 是单调不减的，因此每次区间加都是一段连续的区间.

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MAX int32 = 1 << 19

	var n int32
	fmt.Fscan(in, &n)

	seg := NewSegmentTreeDual32(MAX, func(i int32) Id { return i })
	for i := int32(0); i < n; i++ {
		var l, r int32
		fmt.Fscan(in, &l, &r)

		a := MinLeft32(MAX, func(x int32) bool { return seg.Get(x) >= l }, 0)
		b := MinLeft32(MAX, func(x int32) bool { return seg.Get(x) > r }, a)
		seg.Update(a, b, 1)
	}

	res := seg.GetAll()

	var q int32
	fmt.Fscan(in, &q)
	for i := int32(0); i < q; i++ {
		var x int32
		fmt.Fscan(in, &x)
		fmt.Println(res[x])
	}
}

// RangeAddPointGet

type Id = int32

const COMMUTATIVE = true

func (*SegmentTreeDual32) id() Id                 { return 0 }
func (*SegmentTreeDual32) composition(f, g Id) Id { return f + g }

type SegmentTreeDual32 struct {
	n         int32
	size, log int32
	lazy      []Id
	unit      Id
}

func NewSegmentTreeDual32(n int32, f func(i int32) Id) *SegmentTreeDual32 {
	res := &SegmentTreeDual32{}
	log := int32(1)
	for 1<<log < n {
		log++
	}
	size := int32(1 << log)
	lazy := make([]Id, 2*size)
	unit := res.id()
	for i := int32(0); i < size; i++ {
		lazy[i] = unit
	}
	for i := int32(0); i < n; i++ {
		lazy[size+i] = f(i)
	}
	res.n = n
	res.size = size
	res.log = log
	res.lazy = lazy
	res.unit = unit
	return res
}
func (seg *SegmentTreeDual32) Get(index int32) Id {
	index += seg.size
	for i := seg.log; i > 0; i-- {
		seg.propagate(index >> i)
	}
	return seg.lazy[index]
}
func (seg *SegmentTreeDual32) Set(index int32, value Id) {
	index += seg.size
	for i := seg.log; i > 0; i-- {
		seg.propagate(index >> i)
	}
	seg.lazy[index] = value
}
func (seg *SegmentTreeDual32) GetAll() []Id {
	for i := int32(0); i < seg.size; i++ {
		seg.propagate(i)
	}
	return seg.lazy[seg.size : seg.size+seg.n]
}
func (seg *SegmentTreeDual32) Update(left, right int32, value Id) {
	if left < 0 {
		left = 0
	}
	if right > seg.n {
		right = seg.n
	}
	if left >= right {
		return
	}
	left += seg.size
	right += seg.size
	if !COMMUTATIVE {
		for i := seg.log; i > 0; i-- {
			if (left>>i)<<i != left {
				seg.propagate(left >> i)
			}
			if (right>>i)<<i != right {
				seg.propagate((right - 1) >> i)
			}
		}
	}
	for left < right {
		if left&1 > 0 {
			seg.lazy[left] = seg.composition(value, seg.lazy[left])
			left++
		}
		if right&1 > 0 {
			right--
			seg.lazy[right] = seg.composition(value, seg.lazy[right])
		}
		left >>= 1
		right >>= 1
	}
}
func (seg *SegmentTreeDual32) propagate(k int32) {
	if seg.lazy[k] != seg.unit {
		seg.lazy[k<<1] = seg.composition(seg.lazy[k], seg.lazy[k<<1])
		seg.lazy[k<<1|1] = seg.composition(seg.lazy[k], seg.lazy[k<<1|1])
		seg.lazy[k] = seg.unit
	}
}
func (st *SegmentTreeDual32) String() string {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i := int32(0); i < st.n; i++ {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(fmt.Sprint(st.Get(i)))
	}
	buf.WriteByte(']')
	return buf.String()
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
// left>=lower.
func MinLeft32(right int32, check func(left int32) bool, lower int32) int32 {
	ok, ng := right, lower-1
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}
