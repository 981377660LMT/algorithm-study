// 区间赋值线段树.
// Api:
//   1. Build(n int32, f func(i int32) E)
//   2. Query(start, end int32) E
//   3. QueryAll() E
//   4. Assign(start, end int32, value E)

package main

import (
	"fmt"
	"math/bits"
	"strconv"
	"strings"
)

// 区间赋值线段树.
type SegmentTreeRangeAssign[E any] struct {
	n    int32
	seg  *segmentTree32[E]
	cut  *fastSet32
	data []E
	e    func() E
	op   func(e1, e2 E) E
	pow  func(e E, x int32) E
}

// template
func NewSegmentTreeRangeAssign[E any](e func() E, op func(e1, e2 E) E, pow func(e E, x int32) E) *SegmentTreeRangeAssign[E] {
	res := &SegmentTreeRangeAssign[E]{e: e, op: op, pow: pow}
	return res
}

func (st *SegmentTreeRangeAssign[E]) Build(n int32, f func(i int32) E) {
	st.n = n
	st.seg = newSegmentTree32(st.e, st.op)
	st.seg.Build(n, f)
	st.cut = newFastSet32From(n, func(i int32) bool { return true })
	st.data = st.seg.GetAll()
}

func (st *SegmentTreeRangeAssign[E]) Query(start, end int32) E {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return st.e()
	}
	a, b, c := st.cut.Prev(start), st.cut.Next(start), st.cut.Prev(end)
	if a == c {
		return st.pow(st.data[a], end-start)
	}
	x := st.pow(st.data[a], b-start)
	y := st.seg.Query(b, c)
	z := st.pow(st.data[c], end-c)
	return st.op(st.op(x, y), z)
}

func (st *SegmentTreeRangeAssign[E]) QueryAll() E {
	return st.seg.QueryAll()
}

func (st *SegmentTreeRangeAssign[E]) Assign(start, end int32, value E) {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return
	}
	a, b := st.cut.Prev(start), st.cut.Next(end)
	if a < start {
		st.seg.Set(a, st.pow(st.data[a], start-a))
	}
	if end < b {
		y := st.data[st.cut.Prev(end)]
		st.data[end] = y
		st.cut.Insert(end)
		st.seg.Set(end, st.pow(y, b-end))
	}
	st.cut.Enumerate(start+1, end, func(i int32) { st.seg.Set(i, st.e()); st.cut.Erase(i) })
	st.data[start] = value
	st.cut.Insert(start)
	st.seg.Set(start, st.pow(value, end-start))
}

func (st *SegmentTreeRangeAssign[E]) GetAll() []E {
	res := make([]E, st.n)
	p := int32(0)
	for p < st.n {
		q := st.cut.Next(p + 1)
		val := st.data[p]
		for i := p; i < q; i++ {
			res[i] = val
		}
		p = q
	}
	return res
}

const INF32 int32 = 1 << 30

type segmentTree32[E any] struct {
	n, size int32
	seg     []E
	e       func() E
	op      func(a, b E) E
}

func newSegmentTree32[E any](e func() E, op func(a, b E) E) *segmentTree32[E] {
	res := &segmentTree32[E]{e: e, op: op}
	return res
}

func (st *segmentTree32[E]) Build(n int32, f func(i int32) E) {
	size := int32(1)
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = st.e()
	}
	for i := int32(0); i < n; i++ {
		seg[i+size] = f(i)
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = st.op(seg[i<<1], seg[i<<1|1])
	}
	st.n = n
	st.size = size
	st.seg = seg
}
func (st *segmentTree32[E]) Get(index int32) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *segmentTree32[E]) Set(index int32, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *segmentTree32[E]) Update(index int32, value E) {
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
func (st *segmentTree32[E]) Query(start, end int32) E {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return st.e()
	}
	leftRes, rightRes := st.e(), st.e()
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
func (st *segmentTree32[E]) QueryAll() E { return st.seg[1] }
func (st *segmentTree32[E]) GetAll() []E {
	res := make([]E, st.n)
	copy(res, st.seg[st.size:st.size+st.n])
	return res
}

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (st *segmentTree32[E]) MaxRight(left int32, predicate func(E) bool) int32 {
	if left == st.n {
		return st.n
	}
	left += st.size
	res := st.e()
	for {
		for left&1 == 0 {
			left >>= 1
		}
		if !predicate(st.op(res, st.seg[left])) {
			for left < st.size {
				left <<= 1
				if tmp := st.op(res, st.seg[left]); predicate(tmp) {
					res = tmp
					left++
				}
			}
			return left - st.size
		}
		res = st.op(res, st.seg[left])
		left++
		if (left & -left) == left {
			break
		}
	}
	return st.n
}

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 predicate
func (st *segmentTree32[E]) MinLeft(right int32, predicate func(E) bool) int32 {
	if right == 0 {
		return 0
	}
	right += st.size
	res := st.e()
	for {
		right--
		for right > 1 && right&1 == 1 {
			right >>= 1
		}
		if !predicate(st.op(st.seg[right], res)) {
			for right < st.size {
				right = right<<1 | 1
				if tmp := st.op(st.seg[right], res); predicate(tmp) {
					res = tmp
					right--
				}
			}
			return right + 1 - st.size
		}
		res = st.op(st.seg[right], res)
		if right&-right == right {
			break
		}
	}
	return 0
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

func newFastSet32From(n int32, f func(i int32) bool) *fastSet32 {
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

// 返回最大的 right 使得 [left,right) 内的值满足 check.
// !注意check内的right不包含，使用时需要right-1.
// right<=upper.
func MaxRight(left int, check func(right int) bool, upper int) int {
	ok, ng := left, upper+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
// left>=lower.
func MinLeft(right int, check func(left int) bool, lower int) int {
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

func abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}
