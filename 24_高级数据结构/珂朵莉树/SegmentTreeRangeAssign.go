package main

import (
	"fmt"
	"math/bits"
	"strconv"
	"strings"
)

func main() {
	seg := NewSegmentTreeRangeAssign(10)
	seg.Assign(0, 7, 1)
	fmt.Println(seg.Query(0, 10))
}

// 区间赋值的线段树.
type SegmentTreeRangeAssign struct {
	n    int
	seg  *SegmentTree
	cut  *FastSet
	data []E
}

const INF int = 1e18

type E = int

func e() E             { return 0 }
func op(a, b E) E      { return a + b }
func pow(a E, k int) E { return a * k }

func NewSegmentTreeRangeAssign(n int) *SegmentTreeRangeAssign {
	return NewSegmentTreeRangeAssignFrom(n, func(i int) E { return e() })
}

func NewSegmentTreeRangeAssignFrom(n int, f func(i int) E) *SegmentTreeRangeAssign {
	res := &SegmentTreeRangeAssign{
		n:   n,
		seg: NewSegmentTree(n, func(i int) E { return f(i) }),
		cut: NewFastSetFrom(n, func(i int) bool { return true }),
	}
	res.data = res.seg.GetAll()
	return res
}

func (st *SegmentTreeRangeAssign) Query(start, end int) E {
	a, b, c := st.cut.Prev(start), st.cut.Next(start), st.cut.Prev(end)
	if a == c {
		return pow(st.data[a], end-start)
	}
	x := pow(st.data[a], b-start)
	y := st.seg.Query(b, c)
	z := pow(st.data[c], end-c)
	return op(op(x, y), z)
}

func (st *SegmentTreeRangeAssign) Assign(start, end int, value E) {
	a, b := st.cut.Prev(start), st.cut.Next(end)
	if a < start {
		st.seg.Set(a, pow(st.data[a], start-a))
	}
	if end < b {
		y := st.data[st.cut.Prev(end)]
		st.data[end] = y
		st.cut.Insert(end)
		st.seg.Set(end, pow(y, b-end))
	}
	st.cut.Enumerate(start+1, end, func(i int) { st.seg.Set(i, e()); st.cut.Erase(i) })
	st.data[start] = value
	st.cut.Insert(start)
	st.seg.Set(start, pow(value, end-start))
}

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

type SegmentTree struct {
	n, size int
	seg     []E
}

func NewSegmentTree(n int, f func(int) E) *SegmentTree {
	res := &SegmentTree{}
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = e()
	}
	for i := 0; i < n; i++ {
		seg[i+size] = f(i)
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}
func NewSegmentTreeFrom(leaves []E) *SegmentTree {
	res := &SegmentTree{}
	n := len(leaves)
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = e()
	}
	for i := 0; i < n; i++ {
		seg[i+size] = leaves[i]
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}
func (st *SegmentTree) Get(index int) E {
	if index < 0 || index >= st.n {
		return e()
	}
	return st.seg[index+st.size]
}
func (st *SegmentTree) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTree) Update(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = op(st.seg[index], value)
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

// [start, end)
func (st *SegmentTree) Query(start, end int) E {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return e()
	}
	leftRes, rightRes := e(), e()
	start += st.size
	end += st.size
	for start < end {
		if start&1 == 1 {
			leftRes = op(leftRes, st.seg[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = op(st.seg[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return op(leftRes, rightRes)
}

func (st *SegmentTree) QueryAll() E { return st.seg[1] }

func (st *SegmentTree) GetAll() []E {
	res := make([]E, st.n)
	copy(res, st.seg[st.size:st.size+st.n])
	return res
}

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTree) MaxRight(left int, predicate func(E) bool) int {
	if left == st.n {
		return st.n
	}
	left += st.size
	res := e()
	for {
		for left&1 == 0 {
			left >>= 1
		}
		if !predicate(op(res, st.seg[left])) {
			for left < st.size {
				left <<= 1
				if tmp := op(res, st.seg[left]); predicate(tmp) {
					res = tmp
					left++
				}
			}
			return left - st.size
		}
		res = op(res, st.seg[left])
		left++
		if (left & -left) == left {
			break
		}
	}
	return st.n
}

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTree) MinLeft(right int, predicate func(E) bool) int {
	if right == 0 {
		return 0
	}
	right += st.size
	res := e()
	for {
		right--
		for right > 1 && right&1 == 1 {
			right >>= 1
		}
		if !predicate(op(st.seg[right], res)) {
			for right < st.size {
				right = right<<1 | 1
				if tmp := op(st.seg[right], res); predicate(tmp) {
					res = tmp
					right--
				}
			}
			return right + 1 - st.size
		}
		res = op(st.seg[right], res)
		if right&-right == right {
			break
		}
	}
	return 0
}

type FastSet struct {
	n, lg int
	seg   [][]int
}

func NewFastSet(n int) *FastSet {
	res := &FastSet{n: n}
	seg := [][]int{}
	n_ := n
	for {
		seg = append(seg, make([]int, (n_+63)>>6))
		n_ = (n_ + 63) >> 6
		if n_ <= 1 {
			break
		}
	}
	res.seg = seg
	res.lg = len(seg)
	return res
}

func NewFastSetFrom(n int, f func(i int) bool) *FastSet {
	res := NewFastSet(n)
	for i := 0; i < n; i++ {
		if f(i) {
			res.seg[0][i>>6] |= 1 << (i & 63)
		}
	}
	for h := 0; h < res.lg-1; h++ {
		for i := 0; i < len(res.seg[h]); i++ {
			if res.seg[h][i] != 0 {
				res.seg[h+1][i>>6] |= 1 << (i & 63)
			}
		}
	}
	return res
}

func (fs *FastSet) Has(i int) bool {
	return (fs.seg[0][i>>6]>>(i&63))&1 != 0
}

func (fs *FastSet) Insert(i int) bool {
	if fs.Has(i) {
		return false
	}
	for h := 0; h < fs.lg; h++ {
		fs.seg[h][i>>6] |= 1 << (i & 63)
		i >>= 6
	}
	return true
}

func (fs *FastSet) Erase(i int) bool {
	if !fs.Has(i) {
		return false
	}
	for h := 0; h < fs.lg; h++ {
		cache := fs.seg[h]
		cache[i>>6] &= ^(1 << (i & 63))
		if cache[i>>6] != 0 {
			break
		}
		i >>= 6
	}
	return true
}

// 返回大于等于i的最小元素.如果不存在,返回n.
func (fs *FastSet) Next(i int) int {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}

	for h := 0; h < fs.lg; h++ {
		cache := fs.seg[h]
		if i>>6 == len(cache) {
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
func (fs *FastSet) Prev(i int) int {
	if i < 0 {
		return -1
	}
	if i >= fs.n {
		i = fs.n - 1
	}

	for h := 0; h < fs.lg; h++ {
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
func (fs *FastSet) Enumerate(start, end int, f func(i int)) {
	x := start - 1
	for {
		x = fs.Next(x + 1)
		if x >= end {
			break
		}
		f(x)
	}
}

func (fs *FastSet) String() string {
	res := []string{}
	for i := 0; i < fs.n; i++ {
		if fs.Has(i) {
			res = append(res, strconv.Itoa(i))
		}
	}
	return fmt.Sprintf("FastSet{%v}", strings.Join(res, ", "))
}

func (fs *FastSet) Size() int {
	return fs.n
}

func (*FastSet) bsr(x int) int {
	return 63 - bits.LeadingZeros(uint(x))
}

func (*FastSet) bsf(x int) int {
	return bits.TrailingZeros(uint(x))
}
