// 可撤销线段树.
// RollbackableSegmentTree/SegmentTreeRollbackable
// !新增了 GetTime 和 Rollback 方法.

package main

import (
	"fmt"
)

func main() {
	seg := NewSegmentTreeRollbackable(10, func(index int32) E { return 10 })
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

func (*SegmentTreeRollbackable) e() E        { return INF32 }
func (*SegmentTreeRollbackable) op(a, b E) E { return min32(a, b) }
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

type SegmentTreeRollbackable struct {
	n, size int32
	seg     *rollbackArray[E]
}

func NewSegmentTreeRollbackable(n int32, f func(int32) E) *SegmentTreeRollbackable {
	res := &SegmentTreeRollbackable{}
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
	res.seg = newRollbackArrayFrom(seg)
	return res
}
func NewSegmentTreeRollbackableFrom(leaves []E) *SegmentTreeRollbackable {
	res := &SegmentTreeRollbackable{}
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
	res.seg = newRollbackArrayFrom(seg)
	return res
}

func (st *SegmentTreeRollbackable) GetTime() int32 {
	return st.seg.GetTime()
}
func (st *SegmentTreeRollbackable) Rollback(time int32) {
	st.seg.Rollback(time)
}

func (st *SegmentTreeRollbackable) Get(index int32) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg.Get(index + st.size)
}
func (st *SegmentTreeRollbackable) Set(index int32, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg.Set(index, value)
	for index >>= 1; index > 0; index >>= 1 {
		st.seg.Set(index, st.op(st.seg.Get(index<<1), st.seg.Get(index<<1|1)))
	}
}
func (st *SegmentTreeRollbackable) Update(index int32, value E) {
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
func (st *SegmentTreeRollbackable) Query(start, end int32) E {
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
func (st *SegmentTreeRollbackable) QueryAll() E { return st.seg.Get(1) }
func (st *SegmentTreeRollbackable) GetAll() []E {
	res := make([]E, st.n)
	copy(res, st.seg.GetAllMut()[st.size:st.size+st.n])
	return res
}

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTreeRollbackable) MaxRight(left int32, predicate func(E) bool) int32 {
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
func (st *SegmentTreeRollbackable) MinLeft(right int32, predicate func(E) bool) int32 {
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

type historyItem[V any] struct {
	index int32
	value V
}

type rollbackArray[V any] struct {
	n       int32
	data    []V
	history []historyItem[V]
}

func newRollbackArray[V any](n int32, f func(index int32) V) *rollbackArray[V] {
	data := make([]V, n)
	for i := int32(0); i < n; i++ {
		data[i] = f(i)
	}
	return &rollbackArray[V]{
		n:    n,
		data: data,
	}
}

func newRollbackArrayFrom[V any](data []V) *rollbackArray[V] {
	return &rollbackArray[V]{n: int32(len(data)), data: data}
}

func (r *rollbackArray[V]) GetTime() int32 {
	return int32(len(r.history))
}

func (r *rollbackArray[V]) Rollback(time int32) {
	for int32(len(r.history)) > time {
		pair := r.history[len(r.history)-1]
		r.history = r.history[:len(r.history)-1]
		r.data[pair.index] = pair.value
	}
}

func (r *rollbackArray[V]) Undo() bool {
	if len(r.history) == 0 {
		return false
	}
	pair := r.history[len(r.history)-1]
	r.history = r.history[:len(r.history)-1]
	r.data[pair.index] = pair.value
	return true
}

func (r *rollbackArray[V]) Get(index int32) V {
	return r.data[index]
}

func (r *rollbackArray[V]) Set(index int32, value V) {
	r.history = append(r.history, historyItem[V]{index: index, value: r.data[index]})
	r.data[index] = value
}

func (r *rollbackArray[V]) GetAll() []V {
	return append(r.data[:0:0], r.data...)
}

func (r *rollbackArray[V]) GetAllMut() []V {
	return r.data
}

func (r *rollbackArray[V]) Len() int32 {
	return r.n
}
