// 可撤销线段树.
// RollbackableSegmentTree32/SegmentTree32Rollbackable
// !新增了 GetTime 和 Rollback 方法.

package main

import (
	"fmt"
)

func main() {
	seg := NewSegmentTree32Rollbackable(10, func(index int32) E { return 10 })
	fmt.Println(seg.GetAll())
	seg.Update(0, 5)
	fmt.Println(seg.GetAll())
	fmt.Println(seg.Query(0, 5))
	time := seg.GetTime()
	seg.Update(0, 3)
	fmt.Println(seg.GetAll())
	seg.Rollback(time)
	fmt.Println(seg.GetAll())
}

const INF32 int32 = 1 << 30

// PointSetRangeMin

type E = int32

func (*SegmentTree32Rollbackable) e() E        { return INF32 }
func (*SegmentTree32Rollbackable) op(a, b E) E { return min32(a, b) }
func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

type SegmentTree32Rollbackable struct {
	n, size int32
	seg     *rollbackArray32
}

func NewSegmentTree32Rollbackable(n int32, f func(int32) E) *SegmentTree32Rollbackable {
	res := &SegmentTree32Rollbackable{}
	size := int32(1)
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := int32(0); i < n; i++ {
		seg[i+size] = f(i)
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = newRollbackArray32From(seg)
	return res
}
func NewSegmentTree32RollbackableFrom(leaves []E) *SegmentTree32Rollbackable {
	res := &SegmentTree32Rollbackable{}
	n := int32(len(leaves))
	size := int32(1)
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := int32(0); i < n; i++ {
		seg[i+size] = leaves[i]
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = newRollbackArray32From(seg)
	return res
}

func (st *SegmentTree32Rollbackable) GetTime() int32 {
	return st.seg.GetTime()
}
func (st *SegmentTree32Rollbackable) Rollback(time int32) {
	st.seg.Rollback(time)
}

func (st *SegmentTree32Rollbackable) Get(index int32) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg.Get(index + st.size)
}
func (st *SegmentTree32Rollbackable) Set(index int32, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg.Set(index, value)
	for index >>= 1; index > 0; index >>= 1 {
		st.seg.Set(index, st.op(st.seg.Get(index<<1), st.seg.Get(index<<1|1)))
	}
}
func (st *SegmentTree32Rollbackable) Update(index int32, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg.Set(index, st.op(st.seg.Get(index), value))
	for index >>= 1; index > 0; index >>= 1 {
		st.seg.Set(index, st.op(st.seg.Get(index<<1), st.seg.Get(index<<1|1)))
	}
}

// [start, end)
func (st *SegmentTree32Rollbackable) Query(start, end int32) E {
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
			leftRes = st.op(leftRes, st.seg.Get(start))
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = st.op(st.seg.Get(end), rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return st.op(leftRes, rightRes)
}
func (st *SegmentTree32Rollbackable) QueryAll() E { return st.seg.Get(1) }
func (st *SegmentTree32Rollbackable) GetAll() []E {
	res := make([]E, st.n)
	copy(res, st.seg.GetAllMut()[st.size:st.size+st.n])
	return res
}

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTree32Rollbackable) MaxRight(left int32, predicate func(E) bool) int32 {
	if left == st.n {
		return st.n
	}
	left += st.size
	res := st.e()
	for {
		for left&1 == 0 {
			left >>= 1
		}
		if !predicate(st.op(res, st.seg.Get(left))) {
			for left < st.size {
				left <<= 1
				if tmp := st.op(res, st.seg.Get(left)); predicate(tmp) {
					res = tmp
					left++
				}
			}
			return left - st.size
		}
		res = st.op(res, st.seg.Get(left))
		left++
		if (left & -left) == left {
			break
		}
	}
	return st.n
}

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTree32Rollbackable) MinLeft(right int32, predicate func(E) bool) int32 {
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
		if !predicate(st.op(st.seg.Get(right), res)) {
			for right < st.size {
				right = right<<1 | 1
				if tmp := st.op(st.seg.Get(right), res); predicate(tmp) {
					res = tmp
					right--
				}
			}
			return right + 1 - st.size
		}
		res = st.op(st.seg.Get(right), res)
		if right&-right == right {
			break
		}
	}
	return 0
}

const mask int = 1<<32 - 1

type rollbackArray32 struct {
	n       int32
	data    []int32
	history []int // (index, value)
}

func newRollbackArray32(n int32, f func(index int32) int32) *rollbackArray32 {
	data := make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = f(i)
	}
	return &rollbackArray32{
		n:    n,
		data: data,
	}
}

func newRollbackArray32From(data []int32) *rollbackArray32 {
	return &rollbackArray32{n: int32(len(data)), data: data}
}

func (r *rollbackArray32) GetTime() int32 {
	return int32(len(r.history))
}

func (r *rollbackArray32) Rollback(time int32) {
	for int32(len(r.history)) > time {
		pair := r.history[len(r.history)-1]
		r.history = r.history[:len(r.history)-1]
		r.data[pair>>32] = int32(pair & mask)
	}
}

func (r *rollbackArray32) Get(index int32) int32 {
	return r.data[index]
}

func (r *rollbackArray32) Set(index int32, value int32) {
	r.history = append(r.history, int(index)<<32|int(r.data[index]))
	r.data[index] = value
}

func (r *rollbackArray32) GetAll() []int32 {
	return append(r.data[:0:0], r.data...)
}

func (r *rollbackArray32) GetAllMut() []int32 {
	return r.data
}

func (r *rollbackArray32) Len() int32 {
	return r.n
}
