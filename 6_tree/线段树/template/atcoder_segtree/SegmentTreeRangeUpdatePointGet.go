// https://maspypy.github.io/library/ds/segtree/dual_segtree.hpp

package main

import (
	"bytes"
	"fmt"
)

type Fancy struct {
	seg *SegmentTreeDual
	len int
}

func Constructor() Fancy {
	return Fancy{seg: NewSegmentTreeDual(1e5 + 10)}
}

func (this *Fancy) Append(val int) {
	this.seg.Update(this.len, this.len+1, Id{add: val, mul: 1})
	this.len++
}

func (this *Fancy) AddAll(inc int) {
	this.seg.Update(0, this.len, Id{add: inc, mul: 1})
}

func (this *Fancy) MultAll(m int) {
	this.seg.Update(0, this.len, Id{mul: m})
}

func (this *Fancy) GetIndex(idx int) int {
	if idx >= this.len {
		return -1
	}
	return this.seg.Get(idx).add
}

const INF int = 1e18
const MOD int = 1e9 + 7

// RangeAffinePointGet

type Id = struct{ mul, add int }

func (*SegmentTreeDual) id() Id { return Id{mul: 1} }
func (*SegmentTreeDual) composition(f, g Id) Id {
	return Id{mul: f.mul * g.mul % MOD, add: (f.mul*g.add + f.add) % MOD}
}

// RangeAssignPointGet

type SegmentTreeDual struct {
	n            int
	size, height int
	lazy         []Id
}

func NewSegmentTreeDual(n int) *SegmentTreeDual {
	res := &SegmentTreeDual{}
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
	res.n = n
	res.size = size
	res.height = height
	res.lazy = lazy
	return res
}

func (seg *SegmentTreeDual) Get(index int) Id {
	index += seg.size
	for i := seg.height; i > 0; i-- {
		seg.propagate(index >> i)
	}
	return seg.lazy[index]
}

func (seg *SegmentTreeDual) Update(left, right int, value Id) {
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
	for i := seg.height; i > 0; i-- {
		if (left>>i)<<i != left {
			seg.propagate(left >> i)
		}
		if (right>>i)<<i != right {
			seg.propagate((right - 1) >> i)
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

func (seg *SegmentTreeDual) propagate(k int) {
	if seg.lazy[k] != seg.id() {
		seg.lazy[k<<1] = seg.composition(seg.lazy[k], seg.lazy[k<<1])
		seg.lazy[k<<1|1] = seg.composition(seg.lazy[k], seg.lazy[k<<1|1])
		seg.lazy[k] = seg.id()
	}
}

func (st *SegmentTreeDual) String() string {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i := 0; i < st.n; i++ {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(fmt.Sprint(st.Get(i)))
	}
	buf.WriteByte(']')
	return buf.String()
}
