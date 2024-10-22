// Api:
//
// NewRangeAddMakeIncreasing(n int32, f func(i int32) int) *RangeAddRangeMakeIncreasing
// Get(i int32) int
// GetAll() []int
// Set(i int32, x int)
// RangeAdd(start, end int32, x int)
// RangeAssign(start, end int32, x int)
// !MakeIncreasing(start, end int32) -> 从start到end的区间变成非严格单调递增.从前缀开始修正.

package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"strconv"
	"strings"
)

func main() {
	demo()
}

func demo() {
	{
		A := []int{1, 2, 3, 4, 5}
		seg := NewRangeAddMakeIncreasing(int32(len(A)), func(i int32) int { return A[i] })

		fmt.Println(seg.GetAll())

		seg.RangeAdd(1, 3, 2)
		fmt.Println(seg.GetAll())
		seg.Set(1, 10)
		seg.MakeIncreasing(1, 3)
		fmt.Println(seg.GetAll())
		seg.MakeIncreasing(0, 5)
		fmt.Println(seg.GetAll())
	}

	{
		A := []int{3, 1, 4}
		seg := NewRangeAddMakeIncreasing(int32(len(A)), func(i int32) int { return A[i] })
		seg.MakeIncreasing(0, 3)
		fmt.Println(seg.GetAll())
	}
}

// SegmentTreeRangeAddRangeIncreasing
type RangeAddRangeMakeIncreasing struct {
	n      int32
	s1, s2 *fastSet32
	seg    *SegmentTreeDual32
}

func NewRangeAddMakeIncreasing(n int32, f func(i int32) int) *RangeAddRangeMakeIncreasing {
	seg := NewSegmentTreeDual32(n, func(i int32) Id { return f(i) })
	s1 := NewFastSet32From(n, func(i int32) bool { return true })
	s2 := newFastSet32(n)
	for i := int32(1); i < n; i++ {
		if f(i-1) > f(i) {
			s2.Insert(i)
		}
	}
	return &RangeAddRangeMakeIncreasing{n: n, s1: s1, s2: s2, seg: seg}
}

func (r *RangeAddRangeMakeIncreasing) Get(i int32) int {
	return r.seg.Get(r.s1.Prev(i))
}

func (r *RangeAddRangeMakeIncreasing) GetAll() []int {
	res := r.seg.GetAll()
	p := int32(0)
	for i := int32(0); i < r.n; i++ {
		if r.s1.Has(i) {
			p = i
		}
		res[i] = res[p]
	}
	return res
}

func (r *RangeAddRangeMakeIncreasing) Set(i int32, x int) {
	r.split(i)
	r.split(i + 1)
	r.seg.Set(i, x)
	r.s2.Insert(i)
	r.s2.Insert(i + 1)
}

func (r *RangeAddRangeMakeIncreasing) RangeAdd(start, end int32, x int) {
	r.split(start)
	r.split(end)
	if x < 0 {
		r.s2.Insert(start)
	}
	if x > 0 {
		r.s2.Insert(end)
	}
	r.seg.Update(start, end, x)
}

func (r *RangeAddRangeMakeIncreasing) RangeAssign(start, end int32, x int) {
	r.split(start)
	r.split(end)
	r.s1.Insert(start)
	r.s1.Insert(end)
	r.s1.Enumerate(start, end, func(i int32) { r.s1.Erase(i) })
	r.s1.Insert(start)
	r.seg.Set(start, x)
}

func (r *RangeAddRangeMakeIncreasing) MakeIncreasing(start, end int32) {
	r.split(start)
	r.split(end)
	r.s2.Enumerate(start+1, end, func(i int32) {
		mx := r.Get(i - 1)
		for i < end {
			r.s2.Erase(i)
			now := r.Get(i)
			if mx < now {
				break
			}
			r.s1.Erase(i)
			i = r.s1.Next(i)
		}
	})
}

func (r *RangeAddRangeMakeIncreasing) split(p int32) {
	if p == 0 || p == r.n || r.s1.Has(p) {
		return
	}
	r.seg.Set(p, r.Get(p))
	r.s1.Insert(p)
}

type Id = int

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
	res := make([]Id, seg.n)
	copy(res, seg.lazy[seg.size:seg.size+seg.n])
	return res
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

type fastSet32 struct {
	n, lg int32
	seg   [][]uint64
	size  int32
}

func newFastSet32(n int32) *fastSet32 {
	res := &fastSet32{n: n}
	seg := [][]uint64{}
	n_ := n
	for {
		seg = append(seg, make([]uint64, (n_+63)>>6))
		n_ = (n_ + 63) >> 6
		if n_ <= 1 {
			break
		}
	}
	res.seg = seg
	res.lg = int32(len(seg))
	return res
}

func NewFastSet32From(n int32, f func(i int32) bool) *fastSet32 {
	res := newFastSet32(n)
	for i := int32(0); i < n; i++ {
		if f(i) {
			res.seg[0][i>>6] |= 1 << (i & 63)
			res.size++
		}
	}
	for h := int32(0); h < res.lg-1; h++ {
		for i := 0; i < len(res.seg[h]); i++ {
			if res.seg[h][i] != 0 {
				res.seg[h+1][i>>6] |= 1 << (i & 63)
			}
		}
	}
	return res
}

func (fs *fastSet32) Has(i int32) bool {
	return (fs.seg[0][i>>6]>>(i&63))&1 != 0
}

func (fs *fastSet32) Insert(i int32) bool {
	if fs.Has(i) {
		return false
	}
	for h := int32(0); h < fs.lg; h++ {
		fs.seg[h][i>>6] |= 1 << (i & 63)
		i >>= 6
	}
	fs.size++
	return true
}

func (fs *fastSet32) Erase(i int32) bool {
	if !fs.Has(i) {
		return false
	}
	for h := int32(0); h < fs.lg; h++ {
		cache := fs.seg[h]
		cache[i>>6] &= ^(1 << (i & 63))
		if cache[i>>6] != 0 {
			break
		}
		i >>= 6
	}
	fs.size--
	return true
}

// 返回大于等于i的最小元素.如果不存在,返回n.
func (fs *fastSet32) Next(i int32) int32 {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}

	for h := int32(0); h < fs.lg; h++ {
		cache := fs.seg[h]
		if i>>6 == int32(len(cache)) {
			break
		}
		d := cache[i>>6] >> (i & 63)
		if d == 0 {
			i = i>>6 + 1
			continue
		}
		// find
		i += fs.bsf(d)
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsf(fs.seg[g][i>>6])
		}

		return i
	}

	return fs.n
}

// 返回小于等于i的最大元素.如果不存在,返回-1.
func (fs *fastSet32) Prev(i int32) int32 {
	if i < 0 {
		return -1
	}
	if i >= fs.n {
		i = fs.n - 1
	}

	for h := int32(0); h < fs.lg; h++ {
		if i == -1 {
			break
		}
		d := fs.seg[h][i>>6] << (63 - i&63)
		if d == 0 {
			i = i>>6 - 1
			continue
		}
		// find
		i += fs.bsr(d) - 63
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsr(fs.seg[g][i>>6])
		}

		return i
	}

	return -1
}

// 遍历[start,end)区间内的元素.
func (fs *fastSet32) Enumerate(start, end int32, f func(i int32)) {
	for x := fs.Next(start); x < end; x = fs.Next(x + 1) {
		f(x)
	}
}

func (fs *fastSet32) String() string {
	res := []string{}
	for i := int32(0); i < fs.n; i++ {
		if fs.Has(i) {
			res = append(res, strconv.Itoa(int(i)))
		}
	}
	return fmt.Sprintf("FastSet{%v}", strings.Join(res, ", "))
}

func (fs *fastSet32) Size() int32 {
	return fs.size
}

func (*fastSet32) bsr(x uint64) int32 {
	return 63 - int32(bits.LeadingZeros64(x))
}

func (*fastSet32) bsf(x uint64) int32 {
	return int32(bits.TrailingZeros64(x))
}
