// https://leetcode.cn/problems/find-x-value-of-array-ii/solutions/3656585/xian-duan-shu-by-tsreaper-nuku/
// !每次操作修改一个值，求以某个元素为开头的子数组，乘积 modk=x 的有几个.
// 线段树维护区间子数组乘积的模

package main

const MaxK int = 5

func resultArray(nums []int, k int, queries [][]int) []int {
	type E struct {
		mul int
		mod [MaxK]int
	}

	e := func() E {
		return E{mul: 1}
	}

	op := func(a, b E) E {
		res := E{}
		res.mul = a.mul * b.mul % k
		res.mod = a.mod    // 子数组完全位于左侧
		for i := range k { // 子数组末尾位于右侧
			res.mod[i*a.mul%k] += b.mod[i]
		}
		return res
	}

	of := func(x int) E {
		res := E{}
		res.mul = x
		res.mod[x%k] = 1
		return res
	}

	seg := NewSegmentTree(e, op)
	seg.Build(len(nums), func(i int) E { return of(nums[i]) })

	n := len(nums)
	res := make([]int, len(queries))
	for qi, qv := range queries {
		index, value, start, x := qv[0], qv[1], qv[2], qv[3]
		seg.Set(index, of(value))
		res[qi] = seg.Query(start, n).mod[x]
	}
	return res
}

// PointSetRangeMulMod
type SegmentTree[E any] struct {
	n, size int
	seg     []E
	e       func() E
	op      func(E, E) E
}

func NewSegmentTree[E any](e func() E, op func(E, E) E) *SegmentTree[E] {
	res := &SegmentTree[E]{
		e:  e,
		op: op,
	}
	return res
}

func (s *SegmentTree[E]) Build(n int, f func(int) E) {
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = s.e()
	}
	for i := 0; i < n; i++ {
		seg[i+size] = f(i)
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = s.op(seg[i<<1], seg[i<<1|1])
	}
	s.n = n
	s.size = size
	s.seg = seg
}

func (s *SegmentTree[E]) Get(index int) E {
	if index < 0 || index >= s.n {
		return s.e()
	}
	return s.seg[index+s.size]
}
func (s *SegmentTree[E]) Set(index int, value E) {
	if index < 0 || index >= s.n {
		return
	}
	index += s.size
	s.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		s.seg[index] = s.op(s.seg[index<<1], s.seg[index<<1|1])
	}
}
func (s *SegmentTree[E]) Update(index int, value E) {
	if index < 0 || index >= s.n {
		return
	}
	index += s.size
	s.seg[index] = s.op(s.seg[index], value)
	for index >>= 1; index > 0; index >>= 1 {
		s.seg[index] = s.op(s.seg[index<<1], s.seg[index<<1|1])
	}
}

// [start, end)
func (s *SegmentTree[E]) Query(start, end int) E {
	if start < 0 {
		start = 0
	}
	if end > s.n {
		end = s.n
	}
	if start >= end {
		return s.e()
	}
	leftRes, rightRes := s.e(), s.e()
	start += s.size
	end += s.size
	for start < end {
		if start&1 == 1 {
			leftRes = s.op(leftRes, s.seg[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = s.op(s.seg[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return s.op(leftRes, rightRes)
}
func (s *SegmentTree[E]) QueryAll() E { return s.seg[1] }
func (s *SegmentTree[E]) GetAll() []E {
	res := make([]E, s.n)
	copy(res, s.seg[s.size:s.size+s.n])
	return res
}
