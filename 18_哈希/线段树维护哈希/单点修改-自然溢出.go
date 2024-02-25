package main

import "fmt"

func main() {
	s := "abcdbc"
	seg := NewSegmentTree(len(s), func(i int) E { return FromElement(uint(s[i])) })
	fmt.Println(seg.Query(0, 3))
	fmt.Println(seg.Query(3, 6))
	seg.Set(3, FromElement('a'))
	fmt.Println(seg.Query(0, 3))
	fmt.Println(seg.Query(3, 6))
}

const N int = 1e5 + 10

var BASEPOW0 [N]uint
var BASEPOW1 [N]uint

func init() {
	BASEPOW0[0] = 1
	BASEPOW1[0] = 1
	for i := 1; i < N; i++ {
		BASEPOW0[i] = BASEPOW0[i-1] * BASE0
		BASEPOW1[i] = BASEPOW1[i-1] * BASE1
	}
}

// PointSetRangeHash
// 131/13331/1713302033171(回文素数)
const BASE0 = 131
const BASE1 = 13331

type E = struct {
	len          int
	hash1, hash2 uint
}

func FromElement(c uint) E {
	return E{len: 1, hash1: c, hash2: c}
}

func (*SegmentTreeHash) e() E { return E{} }
func (*SegmentTreeHash) op(a, b E) E {
	return E{
		len:   a.len + b.len,
		hash1: a.hash1*BASEPOW0[b.len] + b.hash1,
		hash2: a.hash2*BASEPOW1[b.len] + b.hash2,
	}
}

type SegmentTreeHash struct {
	n, size int
	seg     []E
	unit    E
}

func NewSegmentTree(n int, f func(int) E) *SegmentTreeHash {
	res := &SegmentTreeHash{}

	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := 0; i < n; i++ {
		seg[i+size] = f(i)
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	res.unit = res.e()
	return res
}
func (st *SegmentTreeHash) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.unit
	}
	return st.seg[index+st.size]
}
func (st *SegmentTreeHash) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTreeHash) Update(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = st.op(st.seg[index], value)
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

// [start, end)
func (st *SegmentTreeHash) Query(start, end int) E {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return st.unit
	}
	leftRes, rightRes := st.unit, st.unit
	start += st.size
	end += st.size
	for start < end {
		if start&1 == 1 {
			leftRes = st.op(leftRes, st.seg[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = st.op(st.seg[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return st.op(leftRes, rightRes)
}

func (st *SegmentTreeHash) QueryAll() E { return st.seg[1] }

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
